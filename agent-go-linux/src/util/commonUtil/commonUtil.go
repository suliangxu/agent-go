package commonUtil

import (
	"configure"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//var (
//	//CONFIG_PATH string = "./configure/agent.conf"
//	CONFIG_PATH string = "/Users/hhh/Desktop/计算机/go/agent-go/configure/agent.conf"
//)

func Error(message string) {
	hostname, _ := os.Hostname()
	log.SetPrefix(hostname + " ERROR ")
	log.Println(message)
}

func Configure(key string) string {
	data, err := ioutil.ReadFile(configure.ConfigPath)
	//log.Println(configure.ConfigPath)
	//fmt.Println("configpath", configure.ConfigPath)
	if err != nil {
		Error("配置文件读取错误 " + err.Error())
		os.Exit(1)
	}
	strData := string(data)
	datas := strings.Split(strData, "\n")
	for _, d := range datas {

		if strings.Contains(d, "#") {
			continue
		}
		args := strings.Split(d, "=")
		if args[0] == key {
			index := strings.Index(d, "=")
			return strings.TrimSpace(strings.TrimSpace(d[index+1:]))
		}
	}
	return ""
}