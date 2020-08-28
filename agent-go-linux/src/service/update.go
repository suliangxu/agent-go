package service

import (
	"bufio"
	"configure"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func kill(id string) bool {
	fmt.Println("run kill...")
	Info("run kill..." + id)
	time.Sleep(time.Second * 2)

	fmt.Println("kill " + id)
	cmd := exec.Command("kill", id)

	stdout, _ := cmd.StdoutPipe()

	defer stdout.Close()   // 保证关闭输出流


	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		Error("kill error, " + err.Error())
		return false
	}

	reader := bufio.NewReader(stdout)
	for{
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
		Info("kill res " + line)
	}

	Info("kill " + id + " success!")

	return true
}

//func runNewScript(filename string)  {
//
//	os.Chmod(filename, 0777)
//	Info(filename + " set permission 0777 success!")
//
//	command := filename
//	Info("run new script command:" + command)
//	cmd := exec.Command(command)
//
//	stdout, _ := cmd.StdoutPipe()
//
//	defer stdout.Close()   // 保证关闭输出流
//
//	if err := cmd.Start(); err != nil {
//		fmt.Println(err)
//		Error("run new script err: " + err.Error())
//		return
//	}
//
//	configure.ProcessPid = strconv.Itoa(cmd.Process.Pid)
//	Info("新脚本开始运行，pid = " + configure.ProcessPid)
//
//	reader := bufio.NewReader(stdout)
//	for{
//		line, err2 := reader.ReadString('\n')
//		if err2 != nil || io.EOF == err2 {
//			break
//		}
//		fmt.Println(line)
//	}
//}

func download(url string, filename string)  {
	res, err := http.Get(url)
	if err != nil {
		Info("get file error: " + err.Error())
		panic(err)
	}
	f, err := os.Create(filename)
	os.Chmod(filename, 0777)
	defer f.Close()
	Info("创建文件 " + filename)
	if err != nil {
		Error("写入文件" + filename + "失败！")
		panic(err)
	}
	io.Copy(f, res.Body)
}

//func scriptUpdate(url string, filename string)  {
//
//	filename = configure.ProjectPath + "/newScripts/" + filename
//
//	download(url, filename)
//
//	if ! kill(configure.ProcessPid){
//		return
//	}
//
//	runNewScript(filename)
//
//}

func scriptUpdate()  {
	configure.UpdateScriptSignal = true

	for i := 0; i<configure.ScriptsNum; i++{
		configure.ScriptUpdateChannal <- true
	}

	Info("old scripts have stopped")

	RemoveDir(configure.ScriptsPath + "/scripts")
	Info("old dir has removed")

	Info("start new scripts running")

	RunAllScriptandPush()
}

func CopyFile(src string, dst string)  {
	srcFile, err := os.Open(src)
	defer srcFile.Close()

	if err != nil {
		Error("open file " + src + "err: " + err.Error())
		panic(err)
	}

	f, err := os.Create(dst)
	defer f.Close()
	Info("创建/清空 文件 " + dst)
	if err != nil {
		Error("写入文件" + dst + "失败！")
		panic(err)
	}
	io.Copy(f, srcFile)

	Info("copy file success. src: " + src + "   dst: " + dst)
}
