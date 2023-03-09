package utils

import "sync"

type Function func()

func RunSync(callBacks ...func()) {
	waitGroup := sync.WaitGroup{}
	for i, e := range callBacks {
		go callBackSync(e, i, &waitGroup)
	}
	waitGroup.Wait()
}

func RunChannel(callBacks ...func()) {
	chanel := make(chan int, len(callBacks))
	for i, e := range callBacks {
		go callBackChanel(e, i, chanel)
	}
	<-chanel
}

func RunFuncThreads(sources []Function, maxChannel int) {
	totalInstance := len(sources)
	if totalInstance < maxChannel {
		maxChannel = totalInstance
	}

	var ch = make(chan Function, totalInstance)
	var wg sync.WaitGroup
	wg.Add(maxChannel)

	for i := 0; i < maxChannel; i++ {
		go func(chanel chan Function) {
			for {
				callBack, ok := <-chanel
				if !ok {
					defer wg.Done()
					return
				}
				callBack()
			}
		}(ch)
	}

	for i := 0; i < totalInstance; i++ {
		ch <- sources[i]
	}

	close(ch)
	wg.Wait()
}

func RunValueThreads[T any](sources []T, maxChannel int, completedFunc func(T) error) []error {
	totalInstance := len(sources)
	if totalInstance < maxChannel {
		maxChannel = totalInstance
	}

	var ch = make(chan T, totalInstance)
	var chErrs = make(chan error, totalInstance)
	var wg sync.WaitGroup
	wg.Add(maxChannel)

	for i := 0; i < maxChannel; i++ {
		go func(errs chan error, chanel chan T) {
			for {
				elem, ok := <-chanel
				if !ok {
					defer wg.Done()
					return
				}

				errs <- completedFunc(elem)
			}
		}(chErrs, ch)
	}

	for i := 0; i < totalInstance; i++ {
		ch <- sources[i]
	}
	close(ch)
	wg.Wait()

	result := make([]error, totalInstance)
	for i := 0; i < totalInstance; i++ {
		if e, ok := <-chErrs; ok {
			result[i] = e
		}
	}

	return result
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
