package loadgen

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

/*
================================1、Generator==========================================
*/
//载荷发生器状态2
const (
	STATUS_ORIGINAL uint32 = 0
	STATUS_STARTING uint32 = 1
	STATUS_STARTED  uint32 = 2
	STATUS_STOPPING uint32 = 3
	STATUS_STOPPED  uint32 = 4
)

//载荷发生器的接口
type Generator interface {
	Start() bool
	Stop() bool
	Status() uint32
	CallCount() int64
}

//载荷发生器
type myGenerator struct {
	caller      Caller
	timeoutNS   time.Duration      //处理超时时间
	durationNS  time.Duration      //载荷持续时间
	lps         uint32             //每秒载荷量
	concurrency uint32             //载荷并发量
	tickets     GoTickets          //gorouting票池
	ctx         context.Context    //上下文
	cancelFunc  context.CancelFunc //取消函数
	callCount   int64              //调用计数
	status      uint32             //状态
	resultCh    chan *CallResult   //调用结果通道
}

func NewGenerator(pset ParamSet) (Generator, error) {

	log.Println("新建一个载荷发生器。。。")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		durationNS: pset.DurationNS,
		lps:        pset.LPS,
		resultCh:   pset.ResultCh,
		status:     STATUS_ORIGINAL,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("初始化一个载荷发生器。。。")
	var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	buf.WriteString(fmt.Sprintf("完成。（concurency=%d)", gen.concurrency))
	log.Println(buf.String())
	return nil
}

func (gen *myGenerator) Start() bool {
	log.Println("启动载荷发生器...")
	if !atomic.CompareAndSwapUint32(&gen.status, STATUS_ORIGINAL, STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(&gen.status, STATUS_STOPPED, STATUS_STARTING) {
			return false
		}
	}
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		log.Printf("设置节流阀（%v）。。。", interval)
		throttle = time.Tick(interval)
	}
	gen.ctx, gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)
	gen.callCount = 0
	atomic.StoreUint32(&gen.status, STATUS_STARTED)
	go func() {
		log.Println("开始生成载荷。。。")
		gen.genLoad(throttle)
		log.Printf("停止。（调用计数：%d）", gen.callCount)
	}()
	return true
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}
		gen.asynCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

func (gen *myGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(&gen.status, STATUS_STARTED, STATUS_STOPPING) {
		return false
	}
	gen.cancelFunc()
	for {
		if atomic.LoadUint32(&gen.status) == STATUS_STOPPED {
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

func (gen *myGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}

func (gen *myGenerator) asynCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("异步调用恐慌！（错误：%s）", err)
				} else {
					errMsg = fmt.Sprintf("异步调用恐慌！（线索：%#v）", err)
				}
				log.Println(errMsg)
				result := &CallResult{
					ID:   -1,
					Code: RET_CODE_FATAL_CALL,
					Msg:  errMsg,
				}
				gen.sendResult(result)
			}
			gen.tickets.Return()
		}()
		rawReq := gen.caller.BuildReq()
		//调用状态：0-未调用；1-调用完成；2-调用超时
		var callStatus uint32
		timer := time.AfterFunc(gen.timeoutNS, func() {
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			result := &CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("超时！（期望时间：< %v）", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.sendResult(result)
		})
		rawResp := gen.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		timer.Stop()
		var result *CallResult
		if rawResp.Err != nil {
			result = &CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   RET_CODE_ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse,
			}
		} else {
			result = gen.caller.CheckReqAndResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()

}

func (gen *myGenerator) callOne(rawReq *RawReq) *RawResp {
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &RawResp{ID: -1, Err: errors.New("非法的请求！")}
	}
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	var rawResp RawResp
	if err != nil {
		errMsg := fmt.Sprintf("异步调用错误：%s.", err)
		rawResp = RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime,
		}
	} else {
		rawResp = RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime,
		}
	}
	return &rawResp
}

func (gen *myGenerator) sendResult(result *CallResult) bool {
	if atomic.LoadUint32(&gen.status) != STATUS_STARTED {
		gen.printIgnoredResult(result, "载荷发生器已经停止了！")
		return false
	}
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "结果通道值满！")
		return false
	}
}

func (gen *myGenerator) printIgnoredResult(result *CallResult, cause string) {
	resultMsg := fmt.Sprintf("ID:=%d,Code=%d,Msg=%s,Elapse=%v", result.ID, result.Code, result.Msg, result.Elapse)
	log.Printf("忽略的结果：%s.（原因：%s）\n", resultMsg, cause)
}

func (gen *myGenerator) prepareToStop(ctxErr error) {
	log.Printf("准备停止载荷发生器（原因：%s）...", ctxErr)
	atomic.CompareAndSwapUint32(&gen.status, STATUS_STARTED, STATUS_STOPPING)
	log.Println("关闭返回值通道！")
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, STATUS_STOPPED)
}

/*
================================3、Caller==========================================
*/
//原生请求
type RawReq struct {
	ID  int64
	Req []byte
}

//原生响应
type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration //响应的耗时
}

//调用结果代码类型
const (
	RET_CODE_SUCCESS              RetCode = 0    //成功
	RET_CODE_WARNING_CALL_TIMEOUT         = 1001 //调用超时警告
	RET_CODE_ERROR_CALL                   = 2001 // 调用错误。
	RET_CODE_ERROR_RESPONSE               = 2002 // 响应内容错误。
	RET_CODE_ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
	RET_CODE_FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

//根据结果代码返回相应的文字解释
func GetRetCodePlain(code RetCode) string {
	var codePlain string
	switch code {
	case RET_CODE_SUCCESS:
		codePlain = "调用成功！"
	case RET_CODE_FATAL_CALL:
		codePlain = "调用过程中发生了致命错误！"
	case RET_CODE_ERROR_CALEE:
		codePlain = "调用错误！"
	case RET_CODE_ERROR_RESPONSE:
		codePlain = "响应内容错误！"
	case RET_CODE_ERROR_CALL:
		codePlain = "被测软件的内部错误！"
	case RET_CODE_WARNING_CALL_TIMEOUT:
		codePlain = "调用超时警告！"
	default:
		codePlain = "未知返回结果！"
	}
	return codePlain
}

//载荷发生器的调用器
type Caller interface {
	BuildReq() RawReq
	Call(req []byte, timeoutNS time.Duration) (resp []byte, err error)
	CheckReqAndResp(req RawReq, resp RawResp) *CallResult
}

/*
================================4、GoTickets==========================================
*/
type GoTickets interface {
	Take()             //取一张票
	Return()           //归还一张票
	Active() bool      //是否激活
	Total() uint32     //总票数
	Remainder() uint32 //剩余的票数
}

func NewGoTickets(total uint32) (GoTickets, error) {
	mytickets := myGoTickets{}
	if !mytickets.init(total) {
		errMsg := fmt.Sprintf("gorouting票不能被初始化！（total=%d）\n", total)
		return nil, errors.New(errMsg)
	}
	return &mytickets, nil
}

type myGoTickets struct {
	total    uint32
	ticketCh chan struct{} //票的容器
	active   bool          //是否激活
}

func (mytickets *myGoTickets) init(total uint32) bool {
	if mytickets.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	mytickets.ticketCh = ch
	mytickets.total = total
	mytickets.active = true
	return true
}

func (mytickets *myGoTickets) Take() {
	<-mytickets.ticketCh
}

func (mytickets *myGoTickets) Return() {
	mytickets.ticketCh <- struct{}{}
}

func (mytickets *myGoTickets) Active() bool {
	return mytickets.active
}

func (mytickets *myGoTickets) Total() uint32 {
	return mytickets.total
}

func (mytickets *myGoTickets) Remainder() uint32 {
	return uint32(len(mytickets.ticketCh))
}

/*
================================5、CallResult==========================================
*/
type CallResult struct {
	ID     int64
	Req    RawReq
	Resp   RawResp
	Code   RetCode
	Msg    string
	Elapse time.Duration
}

type RetCode int

/*
================================6、ParamSet==========================================
*/

type ParamSet struct {
	Caller     Caller
	TimeoutNS  time.Duration
	LPS        uint32
	DurationNS time.Duration
	ResultCh   chan *CallResult
}

func (pset *ParamSet) Check() error {
	var errMsgs []string
	if pset.Caller == nil {
		errMsgs = append(errMsgs, "非法的Caller！")
	}
	if pset.DurationNS == 0 {
		errMsgs = append(errMsgs, "非法的durationNS！")
	}
	if pset.TimeoutNS == 0 {
		errMsgs = append(errMsgs, "非法的timeoutNS！")
	}
	if pset.LPS == 0 {
		errMsgs = append(errMsgs, "非法lps！")
	}
	if pset.ResultCh == nil {
		errMsgs = append(errMsgs, "非法的resultCh！")
	}
	var buf bytes.Buffer
	buf.WriteString("开始检查参数值。。。")
	if errMsgs != nil {
		errMsg := strings.Join(errMsgs, " ")
		buf.WriteString(fmt.Sprintf("未通过！（%s）", errMsg))
		log.Println(buf.String())
		return errors.New(errMsg)
	}
	buf.WriteString(fmt.Sprintf("通过！（timeoutNS=%s,lps=%d,durationNS=%s)", pset.TimeoutNS, pset.LPS, pset.DurationNS))
	log.Println(buf.String())
	return nil
}

/*
================================7、Caller调用器的实现==========================================
*/
const (
	DELIM = '\n'
)

//操作符切片
var operators = []string{"+", "-", "*", "/"}

type TCPComm struct {
	addr string
}

func NewTCPComm(addr string) Caller {
	return &TCPComm{addr: addr}
}

func (comm *TCPComm) BuildReq() RawReq {
	id := time.Now().UnixNano()
	sreq := ServerReq{
		ID: id,
		Operands: []int{
			int(rand.Int31n(1000) + 1),
			int(rand.Int31n(1000) + 1),
		},
		Operator: func() string {
			return operators[rand.Int31n(100)%4]
		}(),
	}
	bytes, err := json.Marshal(sreq)
	if err != nil {
		panic(err)
	}
	rawReq := RawReq{ID: id, Req: bytes}
	return rawReq
}

func (comm *TCPComm) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", comm.addr, timeoutNS)
	if err != nil {
		return nil, err
	}
	_, err = write(conn, req, DELIM)
	if err != nil {
		return nil, err
	}
	return read(conn, DELIM)
}

func (comm *TCPComm) CheckReqAndResp(rawReq RawReq, rawResp RawResp) *CallResult {
	var commResult CallResult
	commResult.ID = rawResp.ID
	commResult.Req = rawReq
	commResult.Resp = rawResp
	var sreq ServerReq
	err := json.Unmarshal(rawReq.Req, &sreq)
	if err != nil {
		commResult.Code = RET_CODE_FATAL_CALL
		commResult.Msg = fmt.Sprintf("请求的格式化错误：%s!\n", string(rawReq.Req))
		return &commResult
	}
	var sresp ServerResp
	err = json.Unmarshal(rawResp.Resp, &sresp)
	if err != nil {
		commResult.Code = RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("响应的格式化错误：%s!\n", string(rawResp.Resp))
		return &commResult
	}
	if sreq.ID != sresp.ID {
		commResult.Code = RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("原始请求与原始响应的id不相同。(%d != %d)\n", rawReq.ID, rawResp.ID)
		return &commResult
	}
	if sresp.Err != nil {
		commResult.Code = RET_CODE_ERROR_CALEE
		commResult.Msg = fmt.Sprintf("被测软件不正常：%s\n", sresp.Err)
		return &commResult
	}
	if sresp.Result != op(sreq.Operands, sreq.Operator) {
		commResult.Code = RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf(
			"结果值不相等：%s\n",
			genFormula(sreq.Operands, sreq.Operator, sresp.Result, false))
		return &commResult
	}
	commResult.Code = RET_CODE_SUCCESS
	commResult.Msg = fmt.Sprintf("成功。(%s)\n", sresp.Formula)
	return &commResult
}

func read(conn net.Conn, delim byte) ([]byte, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return nil, err
		}
		readByte := readBytes[0]
		if readByte == delim {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.Bytes(), nil
}

func write(conn net.Conn, req []byte, delim byte) (int, error) {
	writer := bufio.NewWriter(conn)
	n, err := writer.Write(req)
	if err == nil {
		writer.WriteByte(delim)
	}
	if err == nil {
		writer.Flush()
	}
	return n, err
}

/*
================================8、被测软件实现==========================================
*/
type ServerReq struct {
	ID       int64
	Operands []int
	Operator string
}

type ServerResp struct {
	ID      int64
	Formula string
	Result  int
	Err     error
}

func op(operands []int, operator string) int {
	var result int
	switch {
	case operator == "+":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result = result + v
			}
		}
	case operator == "-":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result = result - v
			}
		}
	case operator == "*":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result = result * v
			}
		}
	case operator == "/":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result = result / v
			}
		}
	}
	return result
}

func genFormula(operands []int, operator string, result int, equal bool) string {
	var buff bytes.Buffer
	n := len(operands)
	for i := 0; i < n; i++ {
		if i > 0 {
			buff.WriteString(" ")
			buff.WriteString(operator)
			buff.WriteString(" ")
		}
		buff.WriteString(strconv.Itoa(operands[i]))
	}
	if equal {
		buff.WriteString("=")
	} else {
		buff.WriteString("!=")
	}
	buff.WriteString(strconv.Itoa(result))
	return buff.String()
}

func reqHandler(conn net.Conn) {
	var errMsg string
	var sresp ServerResp
	req, err := read(conn, DELIM)
	if err != nil {
		errMsg = fmt.Sprintf("服务器：请求读取失败：%s", errMsg)
	} else {
		var sreq ServerReq
		err := json.Unmarshal(req, &sreq)
		if err != nil {
			errMsg = fmt.Sprintf("服务器：请求解析失败：%s", err)
		} else {
			sresp.ID = sreq.ID
			sresp.Result = op(sreq.Operands, sreq.Operator)
			sresp.Formula = genFormula(sreq.Operands, sreq.Operator, sresp.Result, true)
		}
	}
	if errMsg != "" {
		sresp.Err = errors.New(errMsg)
	}
	bytes, err := json.Marshal(sresp)
	if err != nil {
		log.Fatalf("服务器：响应整顿失败：%s", err)
	}
	_, err = write(conn, bytes, DELIM)
	if err != nil {
		log.Fatalf("服务器：响应写入失败：%s", err)
	}
}

type TCPServer struct {
	listener net.Listener
	active   uint32 //0-未激活；1-激活
}

func NewTCPServer() *TCPServer {
	return &TCPServer{}
}

func (server *TCPServer) init(addr string) error {
	if !atomic.CompareAndSwapUint32(&server.active, 0, 1) {
		return nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		atomic.StoreUint32(&server.active, 0)
		return err
	}
	server.listener = ln
	return nil
}

func (server *TCPServer) Listen(addr string) error {
	err := server.init(addr)
	if err != nil {
		return err
	}
	go func() {
		for {
			if atomic.LoadUint32(&server.active) != 1 {
				break
			}
			conn, err := server.listener.Accept()
			if err != nil {
				if atomic.LoadUint32(&server.active) == 1 {
					log.Fatalf("服务器：请求接收失败：%s\n", err)
				} else {
					log.Println("服务器：接收阻塞，因为连接已关闭")
				}
				continue
			}
			go reqHandler(conn)
		}
	}()
	return nil
}

func (server *TCPServer) Close() bool {
	if !atomic.CompareAndSwapUint32(&server.active, 1, 0) {
		return false
	}
	server.listener.Close()
	return true
}
