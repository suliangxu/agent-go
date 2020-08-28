package service

import (
	"configure"
	"os"
	"strings"
	"time"
)

func Info(info string)  {
	info = strings.Split(time.Now().String(),".")[0] +" INFO " + info + "\n"
	f, _ := os.OpenFile(configure.LogPath + "/log_info.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f.Write([]byte(info))
	f.Close()
}

func Error(info string)  {
	info = strings.Split(time.Now().String(),".")[0] +" ERROR " + info + "\n"
	f, _ := os.OpenFile(configure.LogPath + "/log_error.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f.Write([]byte(info))
	f.Close()
}

func Debug(info string)  {
	info = strings.Split(time.Now().String(),".")[0] +" DEBUG " + info + "\n"
	if configure.IsDebug {
		f, _ := os.OpenFile(configure.LogPath + "/log_info.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		f.Write([]byte(info))
		f.Close()
	}
}

func ClearLog(LOG_PATH string) {

	for{
		<- time.After(time.Hour * 1)

		file, err := os.Stat(LOG_PATH)
		if err != nil {
			continue
		}

		logSize := os.FileInfo(file).Size()
		if logSize >= configure.LogSize {
			os.Truncate(LOG_PATH, 0)
			Info("清空日志文件: " + configure.LogPath + "/log_info.txt")
		}
	}
}