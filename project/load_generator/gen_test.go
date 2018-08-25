package loadgen

import (
	"testing"
	"time"
)

var printDetail = true

func TestStart(t *testing.T) {
	server := NewTCPServer()
	defer server.Close()
	serverAddr := "127.0.0.1:8081"
	t.Logf("启动被测软件(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		t.Fatalf("被测软件启动失败：(地址=%s)", serverAddr)
		t.FailNow()
	}
	pset := ParamSet{
		Caller:     NewTCPComm(serverAddr),
		TimeoutNS:  50 * time.Millisecond,
		DurationNS: 10 * time.Second,
		LPS:        uint32(1000),
		ResultCh:   make(chan *CallResult, 50),
	}
	t.Logf("初始化载荷发生器（超时时间=%v，负载持续时间=%v，每秒载荷量=%d）", pset.TimeoutNS, pset.DurationNS, pset.LPS)
	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("载荷发生器初始化失败：%s\n", err)
		t.FailNow()
	}
	t.Log("启动载荷发生器。。。")
	gen.Start()
	countMap := make(map[RetCode]int)
	for r := range pset.ResultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n", r.ID, r.Code, r.Msg, r.Elapse)
		}
	}
	var total int
	t.Log("开始计数。。。")
	for k, v := range countMap {
		codePlain := GetRetCodePlain(k)
		t.Logf("Code Plain:%s(%d),Count:%d.\n", codePlain, k, v)
		total += v
	}
	t.Logf("Total:%d.\n", total)
	successCount := countMap[RET_CODE_SUCCESS]
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	t.Logf("每秒载荷量：%d;每秒事务量：%f\n", pset.LPS, tps)
}
