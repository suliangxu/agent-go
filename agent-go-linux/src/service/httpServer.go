package service

import (
	"configure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"util/commonUtil"
)

func ProcessRequest(w http.ResponseWriter, FormMap sync.Map) string {

	var res string
	url, ok := FormMap.Load("url")
	if !ok{
		res = "no url parameter"
		fmt.Fprintln(w, res)
		Info(res + " update project : no url parameter")
		return ""
	}

	fmt.Fprintln(w, "true")

	return url.(string)
}

func updateScript(w http.ResponseWriter, r *http.Request)  {

	fmt.Fprintln(w, "true")

	Info(r.RemoteAddr + " connect the server updateScript")

	if configure.UpdateScriptSignal {
		Info("Scripts updating, this request has be passed!")
	} else{
		Info("start update scripts...")
		go scriptUpdate()
	}

}

func updateProject(w http.ResponseWriter, r *http.Request)  {

	var FormMap =  sync.Map{}

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			FormMap.Store(k, strings.Join(v, ""))
		}
	}

	Info(r.RemoteAddr + " connect the server updateProject")

	url := ProcessRequest(w, FormMap)

	if url!=""{
		Info("start update project from" + url)
		go projectUpdate(url)
	}
}


func getAgentInfo(w http.ResponseWriter, r *http.Request)  {

	Debug(r.RemoteAddr + " connect the server getAgentInfo")

	push_map := make(map[string]interface{})

	now := time.Now().Unix()
	push_map["runtime"] = now - configure.StartTime  // 单位为秒

	push_map["runStatus"] = configure.AgentStatus
	push_map["versionInfo"] = configure.AgentVersion
	push_map["hostName"] = configure.HostName

	t := time.Unix(configure.StartTime, 0)
	a  := t.String()
	a = a[:19]
	push_map["startTime"] = a

	mjson,_ :=json.Marshal(push_map)
	mString :=string(mjson)

	fmt.Fprintln(w, mString)

}

func getProcessInfo(w http.ResponseWriter, r *http.Request)  {
	Debug(r.RemoteAddr + " connect the server getProcessInfo")

	var FormMap =  sync.Map{}

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			FormMap.Store(k, strings.Join(v, ""))
		}
	}

	ok := JudgeProcessPara(w, FormMap)
	if !ok {
		fmt.Fprintln(w, "参数不完整")
		Info(r.RemoteAddr + "getProcessInfo server 参数不完整")
		return
	}

	DealProcessInfo(w, FormMap)
}


func getPortInfo(w http.ResponseWriter, r *http.Request)  {
	Debug(r.RemoteAddr + " connect the server getProcessInfo")

	var FormMap =  sync.Map{}

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			FormMap.Store(k, strings.Join(v, ""))
		}
	}

	_, ok := FormMap.Load("size")
	if !ok{
		fmt.Fprintln(w, "参数不完整")
		Info(r.RemoteAddr + "getPortInfo server 参数不完整")
		return
	}

	DealPortInfo(w, FormMap)
}

func getUninstallSignal(w http.ResponseWriter, r *http.Request)  {
	Info(r.RemoteAddr + " connect the server getUninstallSignal")

	fmt.Fprintln(w, "true")

	go func() {
		killDaemon()
		sendIP()
		time.Sleep(time.Second * 3)
		RemoveDir(configure.ProjectPath)
		kill(configure.MainProcessPid)
	}()
}

func HttpServer()  {

	http.HandleFunc("/updateScript", updateScript)
	http.HandleFunc("/updateProject", updateProject)
	http.HandleFunc("/getAgentInfo", getAgentInfo)
	http.HandleFunc("/getProcessInfo", getProcessInfo)
	http.HandleFunc("/getPortInfo", getPortInfo)
	http.HandleFunc("/getUninstallSignal", getUninstallSignal)

	port := ":" + commonUtil.Configure("HttpServerPort")

	Info("建立监听，路径为 http://" + port + "/updataScript")
	Info("建立监听，路径为 http://" + port + "/updataProject")
	Info("建立监听，路径为 http://" + port + "/getAgentInfo")
	Info("建立监听，路径为 http://" + port + "/getProcessInfo")
	Info("建立监听，路径为 http://" + port + "/getPortInfo")
	Info("建立监听，路径为 http://" + port + "/getUninstallSignal")

	runDaemon()
	configure.UpdateScriptSignal = false

	for i:=1;i<=3;i++{
		err := http.ListenAndServe(port, nil)  // 建立监听

		if err==nil{
			break
		}else {
			Error(fmt.Sprintf("http server failed, err:%v\n", err))
			Info("try to start server again. count=" + strconv.Itoa(i))
			time.Sleep(time.Second * 3)
		}

		if i==3{
			Error("http server failed!")
			configure.AgentStatus = false
			return
		}
	}
}
