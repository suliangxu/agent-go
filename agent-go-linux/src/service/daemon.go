package service

import (
	"bufio"
	"configure"
	"fmt"
	"io"
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
		Error("发送结果失败！ err: " + err.Error())
		configure.AgentStatus = false
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	Info("返回结果 " + string(body))

	configure.DaemonPid = string(body)

}

func killDaemon()  {
	//fmt.Println(configure.DaemonPid)
	//kill(configure.DaemonPid)

	filepath := configure.ProjectPath + "/daemon/kill.sh"

	cmd := exec.Command("sh", filepath, configure.DaemonPid)

	stdout, _ := cmd.StdoutPipe()

	defer stdout.Close()   // 保证关闭输出流


	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		Error("kill daemon error, " + err.Error())
		return
	}

	reader := bufio.NewReader(stdout)
	for{
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
		Info("kill daemon res " + line)
	}

	Info("kill " + configure.DaemonPid + " success!")


	fmt.Println("kill Daemon success!")
}