package service

import (
	"configure"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"util/commonUtil"
)

var (
	durationTime time.Duration
)

func Start()  {

	Time := commonUtil.Configure("GetInfoScriptIntervalTime")
	if Time[len(Time)-1] == 's'{
		intervalTime,_ := strconv.Atoi(Time[:len(Time)-1])
		durationTime = time.Duration(intervalTime) * time.Second
	} else if Time[len(Time)-1] == 'm'{
		intervalTime,_ := strconv.Atoi(Time[:len(Time)-1])
		durationTime = time.Duration(intervalTime) * time.Minute
	} else if Time[len(Time)-1] == 'h'{
		intervalTime,_ := strconv.Atoi(Time[:len(Time)-1])
		durationTime = time.Duration(intervalTime) * time.Hour
	}

	//setPath()
	Info("开始运行接收系统信息脚本运行返回的任务")
	scriptName := getInfo()
	if scriptName==""{
		return
	}
	for{
		RunScriptandPush(scriptName)
		<- time.After(durationTime)
	}
}


func RunScriptandPush(scriptName string)  {
	scriptResult := runScript(scriptName)
	if scriptResult==""{
		return
	}

	sciptReturn := MatchResult(scriptResult)

	PushResult(sciptReturn)
}


func getInfo() string {
	My_Url := commonUtil.Configure("GetInfoScriptSrc")

	resp,err := http.PostForm(My_Url, nil)
	if err!=nil{
		Error("从接口获取数据失败！err: " + err.Error())
		configure.AgentStatus = false
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)

	result_map := make(map[string]string)
	Info("获取到内容！ from " + commonUtil.Configure("GetInfoScriptSrc"))

	err = json.Unmarshal([]byte(body), &result_map)
	if err != nil{
		Error("解析接口内容失败！" + err.Error())
		configure.AgentStatus = false
		return ""
	}

	filename := configure.ScriptsPath + "/" + result_map["scriptName"]
	f,_ := os.Create(filename)
	Info("创建文件 " + filename)
	defer f.Close()

	_, err = f.Write([]byte(result_map["scriptContent"]))

	if err != nil {
		Error("写脚本内容至文件失败！" + err.Error())
		configure.AgentStatus = false
		return ""
	}
	Info("将脚本内容写入文件")
	f.Close()

	return filename
}


func runScript(scriptname string) string {
	cmd := exec.Command("python",  scriptname)
	Debug("运行python脚本" + scriptname)

	out, err := cmd.Output()
	if err != nil {
		Error("脚本运行出错 " + err.Error())
		configure.AgentStatus = false
		return ""
	}

	return string(out)
}


func MatchResult(Script_result string) map[string]string {

	push_map := make(map[string]string)

	push_list := [...]string{"cpuNumber", "cpuType", "disk", "hostIP", "hostName", "macAddress", "memory", "osArch", "osName", "osVersion"}

	for i:=0; i<len(push_list); i++{
		pattern := ".*" + push_list[i] + ": (.*).*"

		reg := regexp.MustCompile(pattern)

		res := reg.FindAllStringSubmatch(Script_result, -1)

		result := res[0][1]

		push_map[push_list[i]] = strings.Replace(result, "\r", "", -1)
	}

	return push_map
}


func PushResult(Script_return map[string]string)  {

	My_url := commonUtil.Configure("GetInfoScriptDst")
	Debug("向接口" + commonUtil.Configure("GetInfoScriptDst") + "发送脚本运行结果")

	v := url.Values{}

	for key, value := range(Script_return){
		v.Add(key, value)
	}
	//fmt.Println("上传结果：", v)

	resp, err := http.PostForm(My_url, v)
	if err != nil{
		Error("发送结果失败！ err: " + err.Error())
		configure.AgentStatus = false
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	Debug("发送结果 " + string(body))

}

func DirCreatAndSetPermission(path string)  {

	if err := os.Mkdir(path, 777); err != nil {
		if os.IsExist(err) {
			Info(path + "  exits")
		} else{
			Error("creat dir " + path + " error :" + err.Error())
			return
		}

	} else {
		Info("creat dir " + path + " success!")
		os.Chmod(path, 0777)
		Info(path + " set permission 0777 success!")
	}
}

func ChangeDirPermission(filename string)  {
	cmd := exec.Command("chmod",  "-R", "777", filename)

	err :=cmd.Start()
	if err != nil{
		Error("改变文件夹 " + filename + " err: " + err.Error())
		return
	}
	Info("递归改变文件夹权限为777: " + filename)
}

//UninstallDst=http://172.18.22.21/../UAC/agentList/reciveAgentUninstall
//UninstallDst=http://192.168.210.128:8080/getIP
func sendIP()  {
	Myurl := commonUtil.Configure("UninstallDst")
	ip := configure.LocalIp

	v := url.Values{}
	v.Add("ip", ip)

	http.PostForm(Myurl, v)
}