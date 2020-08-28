package service

import (
	"archive/zip"
	"configure"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	newUrl string = ""
	count int = 0  // 用于接收内容中 横杠（-） 的个数
	//newProjectName string  // 新项目的名称
	newProjectLocalPath string  // 项目文件夹搬移之后的本机路径
	readyUpdateProject chan bool  //杀掉原项目的daemon时，管道传值，开始新项目的运行
)

func projectUpdate(url string, filename string)  {

	newDir := configure.ProjectPath + "/newTempDir"

	t := strings.Index(filename, ".zip")
	if t==-1 {
		Error("更新压缩包不是zip类型，停止更新！")
		return
	}
	configure.NewProjectName = filename[:t]

	//folder := configure.ProjectPath
	//err := os.Rename(folder, folder+"-old")
	//if err != nil {
	//	Error("rename old project error " + err.Error())
	//} else {
	//	Info("文件夹重命名成功")
	//}


	//DirCreatAndSetPermission(newDir)

	//tmpfilename := newDir + "/" + "agent-go-win.zip"
	filename = newDir + "/" + filename
	download(url, filename)

	//t := strings.Index(filename, ".zip")
	//newProjectName = filename[:t]
	//
	//downloadProject(newDir, filename)
	//
	
	/* 1、将新项目解压到当前项目目录下
	  2、kill 当前项目线程
	  3、将logs/log.txt 、 configure/agent.conf 覆盖到新项目对应文件中
	  4、开启新项目线程
	  5、删除原项目
	*/
	setNewProjectLocalPath()
	UnzipProject(filename, newProjectLocalPath+"/"+configure.NewProjectName)
	newProjectLocalPath = newProjectLocalPath + "/" + configure.NewProjectName
	Info("newProjectLocalPath: " + newProjectLocalPath)

	//ChangeDirPermission(newProjectLocalPath)   //win
	//MoveProject()
	//Movefile("configure/agent.conf")
	////ChangeConfigContent()  // win
	////
	//SignalHandler()
	//kill(configure.MainProcessPid)

	ProcessKill()
}

func setNewProjectLocalPath()  {
	t := strings.Split(configure.ProjectPath, "/")

	dst := ""
	t = t[0:len(t)-1]
	for _,v := range t{
		dst = dst + v + "/"
	}
	//dst += configure.NewProjectName
	dst = dst[:len(dst)-1]
	newProjectLocalPath = dst

	Debug("setNewProjectLocalPath: " + newProjectLocalPath)
}

func UnzipProject(zipFile string, unzipPath string) {
	//打开要解包的文件，tarFile是要解包的 .tar 文件的路径
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		Error("打开压缩文件失败！" + err.Error())
		return
	}
	//用 tr.Next() 来遍历包中的文件，然后将文件的数据保存到磁盘中
	for _, file := range reader.File {
		rc, err := file.Open()

		t := file.Name
		index := strings.Index(t, "/")
		//fmt.Println(file.Name, index)
		fileName := unzipPath + "/" + file.Name[index+1:]

		if file.FileInfo().IsDir(){
			//先创建目录
			dir := path.Dir(fileName)
			_, err = os.Stat(dir)
			//如果err 为空说明文件夹已经存在，就不用创建
			if err != nil {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					Error("创建文件夹失败！" + err.Error())
					return
				}
				//Info("创建文件夹 " + dir)
			}
			continue
		}

		//创建空文件，准备写入解压后的数据
		fw, er := os.Create(fileName)
		if er != nil {
			Error("创建文件失败！" + er.Error())
			return
		}
		//Info("创建文件 " + fileName)
		defer fw.Close()
		// 写入解压后的数据
		_, er = io.CopyN(fw, rc, int64(file.UncompressedSize64))
		if er != nil {
			Error("写入解压数据至 " + fileName + " 失败！ " + er.Error())
			return
		}
		//Info("向文件 " + fileName + " 写入数据成功!")
	}

	Info("解压项目至 " + unzipPath + " 成功！")
}

//func downloadProject(dirPath string, filename string)  {
//	data, err := ioutil.ReadFile(filename)
//
//	if err != nil {
//		Error("文件" + filename + "读取错误 " + err.Error())
//		return
//	}
//	strData := string(data)
//	datas := strings.Split(strData, "\n")
//
//	var (
//		projectName string
//	)
//	for i, d := range datas {
//		if i==0 {
//			projectName = d[12:len(d)-1]
//			Info("projectName:" + projectName)
//			configure.NewProjectName = projectName
//			configure.NewProjectPath = dirPath + "/" + projectName
//			DirCreatAndSetPermission(configure.NewProjectPath)
//			Info("NewProjectPath: " + configure.NewProjectPath)
//		} else if i==1 {
//			configure.UpdateUrl = d[4:len(d)-1]
//			Info("UpdateUrl: " + configure.UpdateUrl)
//		} else {
//			if strings.Contains(d, "content:") {
//				continue
//			}
//			processFileContent(d)
//		}
//	}
//	Info("download project success!")
//}

//func processFileContent(info string) {
//	if len(info) <= 3 {
//		return
//	}
//
//	var tempCount int = 0
//	for i, ch := range info{
//		if ch=='-' {
//			continue
//		} else {
//			tempCount = i
//			break
//		}
//	}
//	filename := info[tempCount:len(info)-4]
//	Debug("get filename: " + filename)
//
//	// 通过count得到newUrl
//	getNewUrl(tempCount, filename)
//
//	urlSrc := configure.UpdateUrl+newUrl
//	urlDst := configure.NewProjectPath+newUrl
//
//	if info[len(info)-2]=='d'{
//		DirCreatAndSetPermission(urlDst)
//	}else if info[len(info)-2]=='f'{
//		download(urlSrc, urlDst)
//	}
//}

//func getNewUrl(tempCount int, filename string)  {
//	/* count, newUrl
//	   newUrl = http://192.168.0.103:10010/agent-go-mac2
//	   count = 0
//	   tempcount = 1
//	   filename = bin
//	   dst = http://192.168.0.103:10010/agent-go-mac2/bin
//
//
//	   newUrl = http://192.168.0.103:10010/agent-go-mac2/src/util/commonUtil/commonUtil.go
//	   count = 4
//	   tempcount = 2
//	   filename = script
//	   dst = http://192.168.0.103:10010/agent-go-mac2/src/script
//	*/
//
//	a := strings.Split(newUrl, "/")
//	dis := tempCount - count
//	flag := true
//
//	if dis==1 {
//		newUrl = newUrl + "/" + filename
//	} else if dis==0 {
//		a[len(a)-1] = filename
//		flag = false
//	} else if dis<0 {
//		a[len(a)+dis-1] = filename
//		a = a[:len(a)+dis]
//		flag = false
//	}
//
//	if flag==false{
//		newUrl = ""
//		for _,v := range a{
//			newUrl = newUrl + v + "/"
//		}
//		newUrl = newUrl[:len(newUrl)-1]
//	}
//
//	count = tempCount
//}

//func MoveProject()  {
//	/*
//		src: /Users/hhh/Desktop/计算机/go/agent-go-mac/newTempDir/agent-go-mac2
//		dst: /Users/hhh/Desktop/计算机/go/agent-go-mac2
//	*/
//
//	// INFO NewProjectPath: /Users/hhh/Desktop/计算机/go/agent-go-mac/newTempDir/agent-go-mac2
//	// INFO 获取到agent-go项目的路径为/Users/hhh/Desktop/计算机/go/agent-go-mac
//
//	t := strings.Split(configure.ProjectPath, "/")
//
//	dst := ""
//	t = t[:len(t)-1]
//	for _,v := range t{
//		dst = dst + v + "/"
//	}
//	dst += configure.NewProjectName
//	newProjectLocalPath = dst
//
//	os.Rename(configure.NewProjectPath, dst)
//	Info("move project success!  from " + configure.NewProjectPath + " to " + dst)
//}

//func ChangeConfigContent()  {
//	configPath := newProjectLocalPath + "/configure/agent.conf"
//	data, err := ioutil.ReadFile(configPath)
//
//	if err != nil {
//		Error("文件" + configPath + "读取错误 " + err.Error())
//		return
//	}
//	strData := string(data)
//	datas := strings.Split(strData, "\n")
//
//	content := ""
//	for _,v := range(datas){
//		if strings.Contains(v, "LogPath="){
//			v = "LogPath=" + newProjectLocalPath + "/logs/"
//		}
//
//		content = content + "\n" + v
//	}
//
//	f, err := os.Create(configPath)
//	defer f.Close()
//	Info("重写文件 " + configPath)
//	if err != nil {
//		Error("写入文件" + configPath + "失败！")
//		panic(err)
//	}
//
//	if _, err := f.WriteString(content); err!=nil{
//		Error("写入文件" + configPath + "失败！")
//		panic(err)
//	}
//
//}

//func SignalHandler() {
//	c := make(chan os.Signal, 2)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	go func() {
//		Info("等待 os.Interrupt 信号......")
//		<-c //阻塞等待
//		ProcessKill() //工作
//		os.Exit(0)
//	}()
//}

func ProcessKill()  {
	readyUpdateProject = make(chan bool, 10)
	go func() {
		//path := newProjectLocalPath + "/bin/main.exe"
		//configpath := newProjectLocalPath + "/configure/agent.conf"
		//command := "start /min " + path + configpath
		//Info("command: " + command)
		//cmd := exec.Command("cmd", "/C", "start", "/min", path, configpath)

		<- readyUpdateProject
		path := newProjectLocalPath + "/run.exe"
		command := "start /min " + path
		Info("command: " + command)
		cmd := exec.Command("cmd", "/C", "start", "/min", path)

		Info("run new main process")

		if err := cmd.Start(); err != nil {
			Info("run new process error: " + err.Error())
		}
		Info("new main process pid: " + strconv.Itoa(cmd.Process.Pid))

	}()

	time.Sleep(5 * time.Second)
	//Movefile("logs/log.txt")
	//go RemoveOldProject()
	RemoveOldProjectAndKill()
}

//func Movefile(file string)  {
//	// 将logs/log.txt 、 configure/agent.conf 覆盖到新项目对应文件中
//	/*
//			object: logs/log.txt 、 configure/agent.conf
//		    src: /Users/hhh/Desktop/计算机/go/agent-go-mac
//			dst: /Users/hhh/Desktop/计算机/go/agent-go-mac/newTempDir/agent-go-mac2
//	*/
//
//	// INFO NewProjectPath: /Users/hhh/Desktop/计算机/go/agent-go-mac/newTempDir/agent-go-mac2
//	// INFO 获取到agent-go项目的路径为/Users/hhh/Desktop/计算机/go/agent-go-mac
//
//	var (
//		src string
//		dst string
//	)
//
//	src = configure.ProjectPath + "/" + file
//	dst = newProjectLocalPath + "/" + file
//	CopyFile(src, dst)
//}

func RemoveOldProjectAndKill()  {
	for i := 0; i<configure.ScriptsNum; i++{
		configure.ScriptUpdateChannal <- true
	}

	killDaemon()

	filename := configure.ProjectPath + "/uninstall.bat"

	if !configure.UninstallProjectSignal{
		readyUpdateProject <- true
	}

	fmt.Println("cmd /C " + filename)
	cmd := exec.Command("cmd", "/C",filename)
	fmt.Println("remove old project...")

	if err := cmd.Start(); err != nil {
		fmt.Println("remove old project error: " + err.Error())
	}
	fmt.Println("remove old project success! ")
}