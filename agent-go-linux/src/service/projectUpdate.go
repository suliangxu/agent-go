package service

import (
	"configure"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	newUrl string = ""
	count int = 0  // 用于接收内容中 横杠（-） 的个数
	//newProjectName string  // 新项目的名称
	newProjectLocalPath string  // 项目文件夹搬移之后的本机路径
)

func projectUpdate(url string)  {

	newDir := ""
	t := strings.Split(configure.ProjectPath, "/")

	for i:=1; i<len(t)-1; i++{
		newDir = newDir + "/" + t[i]
	}

	filename := newDir + "/" + "install.sh"
	Info("filename: " + filename)
	configure.NewProjectShellPath = filename
	Info("ProjectPath: " + configure.ProjectPath)
	Info("NewProjectShellPath: " + configure.NewProjectShellPath)

	download(url, filename)


	path := configure.NewProjectShellPath
	//configpath := newProjectLocalPath + "/configure/agent.conf"
	command := "nohup " + "sh " + path + " &"
	Info("command: " + command)
	cmd := exec.Command("nohup", "sh", path, "&")
	fmt.Println("run new main process")

	if err := cmd.Start(); err != nil {
		fmt.Println("run new process error: " + err.Error())
	}
	fmt.Println("new main process pid: " + strconv.Itoa(cmd.Process.Pid))

	//killDaemon()
	//SignalHandler()
	//kill(configure.MainProcessPid)

}


func SignalHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		Info("等待 os.Interrupt 信号......")
		<-c //阻塞等待
		ProcessKill() //工作
		os.Exit(0)
	}()
}

func ProcessKill()  {
	RemoveDir(configure.ProjectPath)
	fmt.Println("Remove success")
	go func() {
		path := configure.NewProjectShellPath
		//configpath := newProjectLocalPath + "/configure/agent.conf"
		command := "nohup " + "sh " + path + " &"
		Info("command: " + command)
		cmd := exec.Command("nohup", "sh", path, "&")
		fmt.Println("run new main process")

		if err := cmd.Start(); err != nil {
			fmt.Println("run new process error: " + err.Error())
		}
		fmt.Println("new main process pid: " + strconv.Itoa(cmd.Process.Pid))

	}()

	time.Sleep(time.Second * 5)
	//Movefile("logs/log.txt")
	//RemoveProject()
}


func RemoveDir(dirName string)  {
	if err := os.RemoveAll(dirName); err != nil {
		Error(err.Error())
	} else {
		Info("Remove success")
	}
}