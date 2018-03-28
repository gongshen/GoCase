package main

import (

	"os/signal"
	"os"
	"fmt"
	"syscall"
	"sync"
	"time"
	"os/exec"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

func main() {
	go func() {
		time.Sleep(5*time.Second)
		sendSignal()
	}()
	handleSignal()
}

func handleSignal() {
	sigRecv:=make(chan os.Signal,1)
	args:=[]os.Signal{syscall.SIGINT,syscall.SIGQUIT}
	signal.Notify(sigRecv,args...)
// 如果系统发送的信号是SIGINT,那么signal处理程序会把它封装好后发送给 sigRecv 和 sigRecv2
// 如果系统发送的信号是SIGBUS,那么signal处理程序会把它封装好后发送给 sigRecv 和 sigRecv2
	sigRecv2:=make(chan os.Signal,1)
	args2:=[]os.Signal{syscall.SIGINT,syscall.SIGBUS}
	signal.Notify(sigRecv2,args2...)

	// 恢复系统默认的操作,这时 sigRecv通道不接收值,那么for语句被一直阻塞
	// 使用close关闭通道
	fmt.Println("Wait for 2 seconds...")
	time.Sleep(2*time.Second)
	fmt.Println("Stop notification...")
	signal.Stop(sigRecv)
	close(sigRecv)
	fmt.Println("Done.[sigRecv]")
	var wg sync.WaitGroup
	wg.Add(2)
	// 因为这两个for循环都会被阻塞,所以才有并发执行
	go func() {
		for sig:=range sigRecv{
			fmt.Printf("Received a signal from sigRecv: %s\n",sig)
		}
		fmt.Printf("End.[sigRecv]")
		wg.Done()
	}()

	go func() {
		for sig:=range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
			fmt.Printf("End.[sigRecv2]")
			wg.Done()
	}()
	wg.Wait()
}

// 获得进程ID
func sendSignal()  {
	cmds:=[]*exec.Cmd{
		exec.Command("ps","aux"),
		exec.Command("grep","signal"),
		exec.Command("grep","grep"),
		exec.Command("grep","-v","go run"),
		exec.Command("awk","{print $2}"),
	}
	output,err:=runCmds(cmds)
	if err!=nil{
		fmt.Printf("Command Execution Error: %s\n",err)
		return
	}
	pids,err:=getPids(output)
	if err!=nil{
		fmt.Printf("PID Parsing Error: %s\n",err)
		return
	}
	fmt.Printf("Target PID: %s\n",pids)

	// 根据pid寻找相关的程序
	for _,pid:=range pids{
		proc,err:=os.FindProcess(pid)
		if err!=nil{
			fmt.Printf("Process find error: %s\n",err)
			return
		}
		sig:=syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process(pid=%d)...\n",sig,proc)
		err=proc.Signal(sig)
		if err!=nil{
			fmt.Printf("Signal Sending Error: %s\n",err)
			return
		}
	}
}

func runCmds(cmds []*exec.Cmd)([]string,error){
	if cmds==nil || len(cmds) == 0{
		return nil,errors.New("The cmd slice is invailed!")
	}
	var first=true
	var output []byte
	var err error
	for _,cmd:=range cmds{
		fmt.Printf("Run command: %v\n",getCmdPlaintext(cmd))
		// 循环读和写
		if !first{
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin=&stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout=&stdoutBuf
		if err=cmd.Start();err!=nil{
			return nil,getError(err,cmd)
		}
		if err=cmd.Wait();err!=nil{
			return nil,getError(err,cmd)
		}
		output:=stdoutBuf.Bytes()
		fmt.Printf("Output:\n%s\n",output)
		if first{
			first=false
		}
	}
	var lines []string
	var outputBuf bytes.Buffer
	for{
		line,err:=outputBuf.ReadBytes('\n')
		if err!=nil{
			if err==io.EOF{
				break
			}else{
				return nil,getError(err,nil)
			}
		}
		lines=append(lines,string(line))
	}
	return lines,nil
}

func getCmdPlaintext(cmd *exec.Cmd)string{
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)
	// 为什么这里使用 [1:]
	// 不然会出现  ps ps aux
	for _,arg:=range cmd.Args[1:]{
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}

// 错误处理函数
func getError(err error,cmd *exec.Cmd)error{
	var errMsg string
	if cmd!=nil{
		errMsg=fmt.Sprintf("%s [%s %v]",err,(*cmd).Path,(*cmd).Args)
	}else {
		errMsg=fmt.Sprintf("%s",err)
	}
	// 这里为什么

	return errors.New(errMsg)
}

func getPids(strs []string)([]int,error){
	var pids []int
	for _,str:=range strs{
		pid,err:=strconv.Atoi(strings.TrimSpace(str))
		if err!=nil{
			return nil,err
		}
		pids=append(pids,pid)
	}
	return pids,nil
}