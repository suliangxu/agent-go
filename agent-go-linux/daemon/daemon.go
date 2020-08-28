package main

/*
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>

void judgeChildProcess(pid_t pc){
	pid_t pr;

	do{
		pr = waitpid(pc, NULL, WNOHANG);
		if(pr == 0){
			printf("No child exited\n");
		}
	} while (pr == 0);

	if (pr == pc){
		printf("successfully release child %d\n", pr);
	} else {
		printf("some error occured\n");
	}
}
 */
import "C"

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	t := strings.Split(path, "/")

	for i:=1; i<len(t)-1; i++{
		newDir = newDir + "/" + t[i]
	}

	return newDir
}


func restart_agent()  {

	projectPath := getProjectPath()

	Info("command: " + "nohup " + projectPath+"/bin/main " + projectPath + "/configure/agent.conf " + "&")
	cmd := exec.Command("nohup", projectPath+"/bin/main", projectPath+"/configure/agent.conf", "&")
	Info("run new main process")

	if err := cmd.Start(); err != nil {
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

	fmt.Println(daemonPid)
	fmt.Fprintln(w, daemonPid)

	Info("mianPid: " + strconv.Itoa(mainPid))

	c <- true
}


func main()  {
	c = make(chan bool, 1)
	daemonPid = os.Getpid()
	info := os.Args[0]
	path = ""
	t := strings.Split(info, "/")

	for i:=1; i<len(t)-1; i++{
		path = path + "/" + t[i]
	}

	logPath = path + "/daemonLog.txt"

	go func() {
		http.HandleFunc("/daemonInfo", getPid)
		http.ListenAndServe(":10010", nil)
	}()

	<- c
	for {
		<- time.After(time.Minute * 1)

		if err := syscall.Kill(mainPid, 0); err == nil {
			Info("进程号" + strconv.Itoa(mainPid) + "存在")

			pid_c := C.int(mainPid)
			C.judgeChildProcess(pid_c)

			continue
		}

		Info("进程号不存在")

		restart_agent()
		<-c
	}

}
