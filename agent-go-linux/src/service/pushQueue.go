package service

import (
	"bytes"
	"configure"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
	"util/commonUtil"
)

type Item map[string]interface{}

type Queue struct {
	Items []Item
}

type IQueue interface {
	New() Queue
	Enqueue(t Item)
	Dequeue(t Item)
	IsEmpty() bool
	Size() int
	EmptyQueue() []Item
}

var(
	q Queue
	q2 Queue
	ResultQueue *Queue
	lock sync.Mutex
)

func JudgePushQueue()  {
	ResultQueue = initQueue()
	//ResultQueue = q

	// 每隔3s，将队列中的执行结果向服务器发送
	for {
		<- time.After(time.Second * 3)

		if ResultQueue.Size() != 0{
			result := ResultQueue.EmptyQueue()
			PushQueueResult(result)
		}
	}
}

func (q *Queue) New() *Queue  {
	q.Items = []Item{}
	return q
}

func (q *Queue) Enqueue(data Item)  {
	lock.Lock()
	//fmt.Println("Enqueue: ", data)
	q.Items = append(q.Items, data)
	lock.Unlock()
}

func (q *Queue) Dequeue() *Item  {
	item := q.Items[0]
	q.Items = q.Items[1: len(q.Items)]
	return &item
}

func (q *Queue) IsEmpty() bool {
	return len(q.Items) == 0
}

func (q *Queue) Size() int {
	return len(q.Items)
}

func initQueue() *Queue {
	var t *Queue
	if q.Items == nil{
		q = Queue{}
		q.New()
		t = &q
	} else if q2.Items == nil{
		q2 = Queue{}
		q2.New()
		t = &q2
	}

	return t
}

func (q *Queue)EmptyQueue() []Item {
	lock.Lock()
	length := q.Size()
	return_slice := q.Items[:length]
	//fmt.Println("脚本运行结果：\n", return_slice)
	q.New()
	lock.Unlock()
	return return_slice
}

func PushQueueResult(Script_return []Item)  {
	My_url := commonUtil.Configure("RunScriptsDst")

	id := getScriptsId(Script_return)
	bytesData, _ := json.Marshal(Script_return)
	//Info("send info " + string(bytesData))
	resp, err := http.Post(My_url,"application/json", bytes.NewReader(bytesData))
	if err != nil{
		Error(id + "发送结果失败！ err: " + err.Error())
		configure.AgentStatus = false
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	Info(id+ "发送结果 " + string(body))
}

func getScriptsId(result []Item) string {
	info := ""
	//fmt.Println(result)
	for _, v := range(result){
		info = info + " " + strconv.Itoa(v["itemId"].(int))
	}

	return info
}