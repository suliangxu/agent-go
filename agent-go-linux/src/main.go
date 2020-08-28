package main

import (
	"configure"
	"net"
	"os"
	"runtime"
	"service"
	"strconv"
	"strings"
	"time"
	"util/commonUtil"
)

func ProcessCongfigPath()  {
	arg := os.Args[1]
	arg = strings.Replace(arg, "\\", "/", -1)
	configure.ConfigPath = arg

	//configure.ConfigPath = os.Args[1]
	//service.Info("获取到配置文件路径为" + configure.ConfigPath)

	// /Users/hhh/Desktop/计算机/go/agent-go/configure/agent.conf
	index := strings.Index(configure.ConfigPath, "/configure/agent.conf")
	configure.ProjectPath = configure.ConfigPath[:index]

	//os.Mkdir(configure.ProjectPath+"/logs", 777)
	//os.Mkdir(configure.ProjectPath+"/newScripts", 777)
	//service.ChangeDirPermission(configure.ProjectPath+"/newScripts")    win
	//service.Info("获取到agent-go项目的路径为" + configure.ProjectPath)
}

func SetLogPath()  {
	logPath := commonUtil.Configure("LogPath")
	if logPath != ""{
		logPath = strings.Replace(logPath, "\\", "/", -1)
		configure.LogPath = logPath
	} else {
		configure.LogPath = configure.ProjectPath + "/logs"
	}

	service.Info("Log path :" + configure.LogPath)
	go service.ClearLog(configure.LogPath + "/log_info.txt")
}

func SetDebug()  {
	isdebug := commonUtil.Configure("Debug")

	configure.IsDebug, _ = strconv.ParseBool(isdebug)
}

func run()  {

	if commonUtil.Configure("RunScripts")=="false" {
		return
	}

	service.RunAllScriptandPush()
}

func GetOsAndIp()  {
	osT := runtime.GOOS
	if strings.Contains(osT, "windows"){
		configure.OsType = 1
	} else if strings.Contains(osT, "linux"){
		configure.OsType = 2
	}
	service.Info("OsType: " + strconv.Itoa(configure.OsType))

	conn, _ := net.Dial("udp", "www.google.com.hk:80")

	defer conn.Close()
	configure.LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
	service.Info("Local Ip: " + configure.LocalIp)

	host, _ := os.Hostname()
	configure.HostName = host
}

func setPath()  {
	//scriptPath := configure.LogPath
	//
	//getScriptDstPath := ""
	//t := strings.Split(scriptPath, "/")
	//t = t[:len(t)-2]
	//for _,v := range t{
	//	getScriptDstPath = getScriptDstPath + v + "/"
	//}

	getScriptDstPath := configure.ProjectPath

	getScriptDstPath += "/script"
	configure.ScriptsPath = getScriptDstPath
	service.Info("下载脚本存放的路径： getScriptDstPath: " + configure.ScriptsPath)

}


/*
main函数运行的内容：
	1、goroutine 运行 HTTPSERVER
	2、goroutine 运行 run（脚本） --> 需要更新的部分
	3、主线程运行 start()
 */

func main()  {

	ProcessCongfigPath()
	SetDebug()
	SetLogPath()

	// 	记录agent开始运行的时间
	configure.StartTime = time.Now().Unix()

	// 记录agent版本信息
	configure.AgentVersion = commonUtil.Configure("AgentVersion")

	// 初始化更新脚本通道
	configure.ScriptUpdateChannal = make(chan bool, 1)

	// 初始化agent运行状态为true
	configure.AgentStatus = true


	pid := strconv.Itoa(os.Getpid())
	service.Info("主进程开始运行，pid = " + pid)
	configure.MainProcessPid = pid

	GetOsAndIp()
	setPath()

	go func() {
		service.Info("start Http Server...")
		time.Sleep(time.Second * 3)
		service.HttpServer()
	}()

	go run()

	//go func() {
	//	for{
	//		<- time.After(time.Second)
	//		fmt.Println("main")
	//	}
	//
	//}()

	service.Start()

}