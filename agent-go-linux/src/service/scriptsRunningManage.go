package service

import (
	"configure"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//ch chan bool
	cpuMax float64 = 5
	sciptAllow string
	cpuNow float64 = 0
	addScriptChannel chan float64
	removeScriptChannel chan float64
	scriptQueue *Queue
	cpuNowLock sync.RWMutex
	scriptAllowLock sync.RWMutex
)

/*
	处理逻辑：
	1、每个脚本运行前，访问变量scriptAllow，查看脚本是否被允许执行
	2、若脚本被允许，则cpuNow增加，向通道发送信号
	3、每个脚本运行结束时，向通道发送信号，cpuNow减少
	4、对于scriptAllow，加读写锁
*/

func manageScripts()  {
	addScriptChannel = make(chan float64, 1)
	removeScriptChannel = make(chan float64)
	scriptQueue = initQueue()
	go addScript()
	go removeScript()
}

func addScript()  {
	for{
		<- time.After(time.Millisecond * 10)
		if scriptQueue.IsEmpty() {
			continue
		}

		tmp_map := scriptQueue.Items[0]

		tmp_cpu := tmp_map["cpu_persent"].(float64)
		
		cpuNowLock.Lock()
		if tmp_cpu+cpuNow <= cpuMax{
			scriptAllowLock.Lock()
			sciptAllow = tmp_map["filename"].(string)
			scriptAllowLock.Unlock()
			tmp_cpu = <- addScriptChannel
			cpuNow += tmp_cpu
			cpuNowLock.Unlock()
			scriptQueue.Dequeue()
			Debug("add script")
			continue
		}
		cpuNowLock.Unlock()
	}
}

func removeScript()  {
	for{
		tmp_cpu := <- removeScriptChannel
		Debug("remove script")
		cpuNowLock.Lock()
		cpuNow -= tmp_cpu
		cpuNowLock.Unlock()
	}
}

func checkAllow(info_map map[string]interface{})  {
	Debug(info_map["filename"].(string) + " wait ")
	filename := info_map["filename"].(string)
	cpu_persent := info_map["cpu_persent"].(float64)
	scriptQueue.Enqueue(info_map)

	for {
		scriptAllowLock.RLock()
		if sciptAllow == filename{
			scriptAllowLock.RUnlock()
			addScriptChannel <- cpu_persent
			Debug(info_map["filename"].(string) + " start ")
			break
		}
		scriptAllowLock.RUnlock()
		<- time.After(time.Millisecond * 10)
	}
}

var cpuMapLock sync.RWMutex

// 用于判断该脚本所需的cpu
func getScriptSource()  {

	num := len(configure.ScriptsInfo)

	for i:=0; i<num; i++{
		if i!=0 && i%5==0{
			<- time.After(time.Second * 5)
		}

		runEachScriptTest(i)
		Info(configure.ScriptsInfo[i]["filename"].(string) + " cpu use " + strconv.FormatFloat(configure.ScriptsInfo[i]["cpu_persent"].(float64), 'E', -1, 32))
	}
	Info("get scripts cpu percent")
}

func runEachScriptTest(i int) {
	ch := make(chan bool)
	filename := configure.ScriptsInfo[i]["filename"].(string)

	go func() {
		t1, t2 := getEachScriptSource(filename, ch)
		configure.ScriptsInfo[i]["cpu_persent"] = t1
		configure.ScriptsInfo[i]["exit_times"] = t2
	}()

	go func() {
		<- time.After(time.Millisecond * 100)
		cmd := exec.Command("python", filename)

		cmd.Output()
	}()

	<- ch

}

func getEachScriptSource(filename string, ch chan bool) (float64, int) {

	count := 0
	times := 0
	max := 0.0

	com := "ps axu | grep " + filename + " | awk '{print $3}'"

	for {
		times ++
		cmd := exec.Command("/bin/sh", "-c", com)

		out, _ := cmd.Output()

		t := strings.Split(string(out), "\n")

		var total float64 = 0.0
		for _, v := range(t){
			if v != ""{
				v_f, _ := strconv.ParseFloat(v, 64)
				total += v_f
			}
		}

		if total > max {
			max = total
		}
		if max>1 && total<1{
			break
		}

		if total == 0{
			count ++
		}

		if count == 10{
			break
		}
		<- time.After(time.Millisecond * 200)
	}

	ch <- true

	return max, times
}

func waitChanExit(infoMap map[string]interface{})  {
	t := time.Duration(200 * infoMap["exit_times"].(int))
	<- time.After(time.Millisecond * t)
	removeScriptChannel <- infoMap["cpu_persent"].(float64)
}