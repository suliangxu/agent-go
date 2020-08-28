package service

import (
	"configure"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
)

func runDaemon()  {
	Info("开启监控程序")
	My_url := "http://localhost:10010/daemonInfo"

	v := url.Values{}

	v.Add("mainPid", configure.MainProcessPid)

	//fmt.Println("上传结果：", v)

	resp, err := http.PostForm(My_url, v)
	if err != nil{
		Error("向监控程序发送结果失败！ err: " + err.Error())
		configure.AgentStatus = false
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	Info("返回结果 " + string(body))

	configure.DaemonPid = string(body)

}

func killDaemon()  {
	fmt.Println(configure.DaemonPid)

	fmt.Println("taskkill /f /pid " + configure.DaemonPid)
	cmd := exec.Command("cmd", "/C", "taskkill", "/F", "/pid", configure.DaemonPid)

	stdout, _ := cmd.StdoutPipe()

	defer stdout.Close()   // 保证关闭输出流


	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		Error("kill daemon error, " + err.Error())
		return
	}

	fmt.Println("kill Daemon success!")
}