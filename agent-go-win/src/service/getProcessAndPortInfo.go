package service

import (
	"configure"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func JudgeProcessPara(w http.ResponseWriter, FormMap sync.Map) bool {
	_, ok := FormMap.Load("orderType")
	if !ok{
		return false
	}

	_, ok = FormMap.Load("orderDesc")
	if !ok{
		return false
	}

	_, ok = FormMap.Load("size")
	if !ok{
		return false
	}

	return true
}

/*
	orderType   1: cpu   2: 内存
	orderDesc   1: 正序  2: 倒序
	size   需要的数量
*/

func DealProcessInfo(w http.ResponseWriter, FormMap sync.Map) {
	var res interface{}

	res, _ = FormMap.Load("orderType")
	orderType := res.(string)

	res, _ = FormMap.Load("orderDesc")
	orderDesc := res.(string)

	res, _ = FormMap.Load("size")
	size := res.(string)

	push_info := RunProcessCommand(orderType, orderDesc, size)

	fmt.Fprintln(w, push_info)
}

func RunProcessCommand(orderType string, orderDesc string, size string) string {
	getProcessScriptPath := configure.ScriptsPath + "/getProcessInfo.py"

	cmd := exec.Command("python", getProcessScriptPath, orderType, orderDesc, size)
	//Error("python " + getProcessScriptPath + " " + orderType + " " + orderDesc + " " + size)

	out, err := cmd.Output()
	if err != nil {
		Error("run get process script err" + err.Error())
	}

	//Error("Process error " + string(out))
	return string(out)
}


func delete_extra_space(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}



func DealPortInfo(w http.ResponseWriter, FormMap sync.Map) {
	var res interface{}
	push_map_slice := make([]map[string]interface{}, 0, 100)

	res, _ = FormMap.Load("size")
	size := res.(string)
	size_int, _ := strconv.Atoi(size)

	//Info("加载参数size：" + size)

	push_map_slice = RunPortCommand(size_int)

	mjson,_ :=json.Marshal(push_map_slice)
	mString :=string(mjson)

	fmt.Fprintln(w, mString)
}


func RunPortCommand(size int) []map[string]interface{} {
	push_map_slice := make([]map[string]interface{}, 0, 100)
	push_map := make(map[string]interface{})

	cmd := exec.Command("netstat", "-ano")

	out, err := cmd.Output()
	if err != nil {
		Error("run port command err" + err.Error())
	}

	//Info("run command result: " + string(out))

	a := strings.Split(string(out), "\n")
	//fmt.Println("len(a) : ", len(a))
	//fmt.Println(a)
	for i:=4; i<len(a)-1; i++{
		//fmt.Println(i)
		t := strings.Split(delete_extra_space(a[i]), " ")
		//Info("i = " + strconv.Itoa(i) + " " + delete_extra_space(a[i]))
		push_map = makePortPushMap(t)
		push_map_slice = append(push_map_slice, push_map)
	}

	if len(push_map_slice) <= size{
		return push_map_slice
	}

	//Info("push_map_slice: ")

	return push_map_slice[:size]
}


func makePortPushMap(info []string) map[string]interface{} {
	push_map := make(map[string]interface{})

	//fmt.Println(info)

	push_map["protocol"] = info[1]
	push_map["localIPPort"] = info[2]
	push_map["remoteIPPort"] = info[3]

	if len(info)==5 {
		push_map["state"] = ""
		push_map["process"] = info[4]
	} else if len(info)==6 {
		push_map["state"] = info[4]
		push_map["process"] = info[5][:len(info[5])-1]
	}

	return push_map
}