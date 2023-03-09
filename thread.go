package utils

import "sync"

func RunSync(callBacks ...func()) {
	waitGroup := sync.WaitGroup{}
	for i, e := range callBacks {
		go callBackSync(e, i, &waitGroup)
	}
	waitGroup.Wait()
}

func RunChannel(callBacks ...func()) {
	chanel := make(chan int)
	for i, e := range callBacks {
		go callBackChanel(e, i, chanel)
	}
	<-chanel
}

func callBackSync(callBack func(), index int, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	defer waitGroup.Done()
	callBack()
}

func callBackChanel(callBack func(), index int, channel chan int) {
	callBack()
	channel <- index
}
