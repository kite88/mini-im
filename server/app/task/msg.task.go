package task

import (
	"fmt"
	"miniIM/server/app/service"
	"time"
)

func Start() {
	go saveMsgData() //协程开启
}

func saveMsgData() {
	fmt.Println("保存消息任务启动")
	for ; ; {
		timer := time.NewTimer(100 * time.Millisecond)
		t := <-timer.C
		isOk, data := service.ReadInList()
		if isOk == true {
			ok := service.SavePersistentData(data)
			if !ok {
				service.SaveToErrorList(data)
			}
		}
		fmt.Printf("保存消息任务正在执行:%v\n", t.Format("2006-01-02 15:04:05"))
	}
}
