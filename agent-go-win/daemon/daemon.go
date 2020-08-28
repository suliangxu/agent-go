package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mainPid int
	daemonPid int
	path string
	logPath string
	c chan bool
)

func Info(info string)  {
	info = strings.Split(time.Now().String(),".")[0] +" " + info + "\n"

	f, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f.Write([]byte(info))
	f.Close()
}

func getProjectPath() string {
	newDir := ""
	t := strings.Split(path, "\\")
	newDir = t[0]
	//fmt.Println("path: " + path)

	for i:=1; i<len(t)-1; i++{
		newDir = newDir + "\\" + t[i]
	}

	//fmt.Println("newPath: " + newDir)

	return newDir
}


func restart_agent()  {

	projectPath := getProjectPath()

	Info("cmd /C start" + "/min " + projectPath+"\\bin\\main " + projectPath + "\\configure\\agent.conf ")
	cmd := exec.Command("cmd", "/C", "start", "/min", projectPath+"\\bin\\main", projectPath+"\\configure\\agent.conf")
	Info("run new main process")

	if err := cmd.Start(); err != nil {
		fmt.Println("run new error")
		Info("run new process error: " + err.Error())
	}
	Info("new main process pid: " + strconv.Itoa(cmd.Process.Pid))

}

func getPid(w http.ResponseWriter, r *http.Request)  {
	var FormMap =  sync.Map{}

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			FormMap.Store(k, strings.Join(v, ""))
		}
	}

	//fmt.Println(r.Form)
	info, _ := FormMap.Load("mainPid")
	//fmt.Println(info)
	mainPid, _ = strconv.Atoi(info.(string))

	//fmt.Println(daemonPid)
	fmt.Fprintln(w, daemonPid)

	Info("mianPid: " + strconv.Itoa(mainPid))

	c <- true
}

func Check() bool {
	//fmt.Println(mainPid)
	cmd := exec.Command("cmd", "/C", "query", "process", "|", "findstr", strconv.Itoa(mainPid))

	out, _ := cmd.Output()

	info := string(out)
	//fmt.Println(info)
	info = strings.Replace(info, "\n", "", -1)

	if info == ""{
		//fmt.Println("false")
		return false
	}

	//fmt.Println("true")
	return true
}

func main()  {
	c = make(chan bool)
	daemonPid = os.Getpid()
	info := os.Args[0]
	path = ""
	t := strings.Split(info, "\\")

	path = t[0]
	for i:=1; i<len(t)-1; i++{
		path = path + "\\" + t[i]
	}

	logPath = path + "\\daemonLog.txt"

	fmt.Println("监控程序已启动!")
	fmt.Println("注意：不要关闭此窗口！！！")

	//fmt.Println(logPath)

	go func() {
		http.HandleFunc("/daemonInfo", getPid)
		http.ListenAndServe(":10010", nil)
	}()


	<- c
	for {
		<- time.After(time.Minute * 1)

		if Check() {
			Info("进程号" + strconv.Itoa(mainPid) + "存在")
			continue
		}

		Info("进程号不存在")

		restart_agent()
		<-c
	}

}
