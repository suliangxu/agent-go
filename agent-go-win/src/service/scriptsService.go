package service

import (
	"bytes"
	"configure"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"util/commonUtil"
)


func RunAllScriptandPush()  {
	Info("开始运行批量接收脚本运行返回的任务")

	for{
		if getScripsInfo(){
			break
		}
		<- time.After(time.Minute * 5)
	}

	go runAllScripts()
	go JudgePushQueue()
}



func getScripsInfo() bool {

	result_map_slice := make([]map[string]interface{}, 0, 100)

	// 构建url
	pars := "?osType=" + strconv.Itoa(configure.OsType) + "&targetIp=" + configure.LocalIp
	My_Url := commonUtil.Configure("RunScriptsSrc")
	My_Url += pars

	resp,err := http.PostForm(My_Url, nil)
	if err!=nil{
		Error("从接口获取脚本数据失败！err: " + err.Error())
		configure.AgentStatus = false
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &result_map_slice)
	if err != nil{
		Error("解析脚本接口内容失败！" + err.Error())
		configure.AgentStatus = false
		return false
	}

	configure.ScriptsNum = len(result_map_slice)
	if configure.ScriptsNum == 0{
		Info("script num = 0")
		return false
	}

	os.Mkdir(configure.ScriptsPath + "/scripts", 777)
	for i:=0; i<configure.ScriptsNum; i++{
		//fmt.Println(result_map_slice[i])
		processScriptInfo(result_map_slice[i], i)
	}

	return true
}


/*
	result_map 的数据结构：
	id  float64 --> int
	source  float64 --> int  1、python 2、shell
	osScript string
	excVal  float --> int
	excUnit  float  -->  int  1、秒  2、分  3、时

	InfoMap 的数据结构：
	id  int
	filename  string
	scriptType  python | sh
	intervalTime  int
*/
func processScriptInfo(result_map map[string]interface{}, i int) bool {
	InfoMap := make(map[string]interface{})

	result_map["id"] = int(result_map["id"].(float64))
	result_map["source"] = int(result_map["source"].(float64))
	result_map["excUnit"] = int(result_map["excUnit"].(float64))
	result_map["excVal"] = int(result_map["excVal"].(float64))

	InfoMap["id"] = result_map["id"]

	var filename string
	if result_map["source"] == 1{
		filename = configure.ScriptsPath + "/scripts/" + strconv.Itoa(result_map["id"].(int)) + ".py"
		InfoMap["filename"] = filename
		InfoMap["scriptType"] = "python"
	} else if result_map["source"] == 2{
		filename = configure.ScriptsPath + "/scripts/" + strconv.Itoa(result_map["id"].(int)) + ".sh"
		InfoMap["filename"] = filename
		InfoMap["scriptType"] = "sh"
	}

	f,_ := os.Create(filename)
	Info("创建文件 " + filename)
	defer f.Close()

	//fmt.Println(result_map["osScript"])
	_, err := f.Write([]byte(result_map["osScript"].(string)))

	if err != nil {
		Error("写脚本" + "(id=" + strconv.Itoa(result_map["id"].(int)) + ")内容至文件失败！" + err.Error())
		return false
	}
	Debug("将脚本"+ "(id=" + strconv.Itoa(result_map["id"].(int)) + ")内容写入文件")
	f.Close()

	if result_map["excUnit"].(int) == 1{
		InfoMap["intervalTime"] = result_map["excVal"]
	} else if result_map["excUnit"].(int) == 2{
		InfoMap["intervalTime"] = result_map["excVal"].(int) * 60
	} else if result_map["excUnit"].(int) == 3{
		InfoMap["intervalTime"] = result_map["excVal"].(int) * 360
	}
	configure.ScriptsInfo = append(configure.ScriptsInfo, InfoMap)
	//fmt.Println(configure.ScriptsInfo)

	return true
}


func runAllScripts()  {
	configure.UpdateScriptSignal = false

	Info("scriptNum: " + strconv.Itoa(configure.ScriptsNum))
	num := configure.ScriptsNum
	for i:=0; i<num; i++{
		if i % 5 == 0 {
			<- time.After(time.Second * 5)
		}
		go runEachScript(configure.ScriptsInfo[i])
	}
}


func runEachScript(infoMap map[string]interface{}) {

	t := infoMap["intervalTime"].(int)
	durationTime = time.Duration(t) * time.Second

	for{
		select {
		case data := <- configure.ScriptUpdateChannal:
			if data{
				Info(strconv.Itoa(infoMap["id"].(int)) + "stopped")
				return
			}
		default:
			cmd := exec.Command(infoMap["scriptType"].(string), infoMap["filename"].(string))
			Debug("运行脚本" + infoMap["filename"].(string))

			out, err := cmd.Output()
			if err != nil {
				Error(strconv.Itoa(infoMap["id"].(int)) + "脚本运行出错 " + err.Error())
				infoMap["code"] = 500
				configure.ScriptsNum --
				scriptReturn := MatchErrorResult(infoMap)
				//PushEachResult(scriptReturn)
				ResultQueue.Enqueue(scriptReturn)
				return
			}
			infoMap["code"] = 200

			scriptReturn := MatchEachResult(string(out), infoMap)
			//PushEachResult(scriptReturn)
			ResultQueue.Enqueue(scriptReturn)

			<- time.After(durationTime)
		}
	}
}


/*
	脚本执行结果返回所需的参数
		osType string
		source int   1、python脚本   2、命令行指令
 		itemId int
		targetIp  string
		code int  200、成功    500、失败
		results []
*/

func MatchErrorResult(infomap map[string]interface{}) map[string]interface{} {
	push_map := make(map[string]interface{})

	push_map["osType"] = configure.OsType
	push_map["itemId"] = infomap["id"]
	push_map["targetIp"] = configure.LocalIp
	push_map["code"] = infomap["code"]
	push_map["results"] = make([]map[string]interface{}, 0, 100)

	if infomap["scriptType"] == "python"{
		push_map["source"] = 1
	} else if infomap["scriptType"] == "sh"{
		push_map["source"] = 2
	}

	return push_map
}

func MatchEachResult(Script_result string, infomap map[string]interface{}) map[string]interface{} {

	push_map := make(map[string]interface{})

	push_map["osType"] = configure.OsType
	push_map["itemId"] = infomap["id"]
	push_map["targetIp"] = configure.LocalIp
	push_map["code"] = infomap["code"]

	if infomap["scriptType"] == "python"{
		push_map["source"] = 1
		push_map = MatchPythonResult(Script_result, push_map)
	} else if infomap["scriptType"] == "sh"{
		push_map["source"] = 2
		push_map = MatchShellResult(Script_result, push_map)
	}

	return push_map
}


func MatchPythonResult(info string, push_map map[string]interface{}) map[string]interface{} {
	body := []byte(info)
	result_map_slice := make([]map[string]interface{}, 0, 100)

	err := json.Unmarshal(body, &result_map_slice)
	if err != nil{
		Error("匹配结果失败， id = " + strconv.Itoa(push_map["itemId"].(int)))
	}

	push_map["results"] = result_map_slice

	return push_map
}


func MatchShellResult(info string, push_map map[string]interface{}) map[string]interface{} {
	result_map_slice := make([]map[string]interface{}, 0, 100)

	info = strings.Trim(info, "\n")
	body := strings.Split(info, "\n")

	for _, value := range body{
		result_map := make(map[string]interface{})
		result_map["name"] = ""
		result_map["value"], _ = strconv.Atoi(value)
		result_map_slice = append(result_map_slice, result_map)
	}

	push_map["results"] = result_map_slice

	return push_map
}


func PushEachResult(Script_return map[string]interface{})  {

	My_url := commonUtil.Configure("RunScriptsDst")
	Debug("向接口" + commonUtil.Configure("RunScriptDst") + "发送脚本" + strconv.Itoa(Script_return["itemId"].(int)) + "运行结果")
	Info("")

	bytesData, _ := json.Marshal(Script_return)
	resp, err := http.Post(My_url,"application/json", bytes.NewReader(bytesData))
	if err != nil{
		Error(strconv.Itoa(Script_return["itemId"].(int)) + "发送结果失败！ err: " + err.Error())
		configure.AgentStatus = false
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	Debug(strconv.Itoa(Script_return["itemId"].(int)) + "发送结果 " + string(body))

}