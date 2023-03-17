package utils

import "sync"

type Function func()

type IThreadIndex[T any] interface {
	GetIndex() int
	GetData() T
}

type threadIndex[T any] struct {
	Index int
	Data  T
}

func (v *threadIndex[T]) GetIndex() int {
	return v.Index
}

func (v *threadIndex[T]) GetData() T {
	return v.Data
}

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

// RunValueIndexThreads is a function that receives a slice of elements of any type 'T',
// and an integer 'maxChannel' that represents the number of goroutines to be executed simultaneously.
// It also takes a function 'completedFunc' that receives an object of type IThreadIndex[T]
// and returns an error, which is the function that will be executed for each element.
//
// The function iterates through the slice, and for each element, it creates an object of type
// 'threadIndex[T]' that holds the index and data of the element, and sends it to the channel 'ch'.
// Then it creates 'maxChannel' goroutines that read from the channel 'ch', and for each element,
// execute the function 'completedFunc' passing the element, and send the returned error to the
// channel 'chErrs'.
//
// After all elements are processed, the function closes the channel 'ch' and waits for all goroutines
// to finish with the WaitGroup 'wg'. Finally, it reads from the channel 'chErrs' to get the errors
// returned by the function 'completedFunc' for each element, and returns them in a slice of errors.
//
// The function returns a slice of errors, where each index corresponds to the index of the element
// in the input slice, and the value is the error returned by the function 'completedFunc' for that element.
func RunValueIndexThreads[T any](sources []T, maxChannel int, completedFunc func(v IThreadIndex[T]) error) []error {
	totalInstance := len(sources)
	if totalInstance < maxChannel {
		maxChannel = totalInstance
	}

	datas := SelectIndex(sources, func(i int, v T) IThreadIndex[T] {
		return &threadIndex[T]{Index: i, Data: v}
	})

	var ch = make(chan IThreadIndex[T], totalInstance)
	var chErrs = make(chan error, totalInstance)

	var wg sync.WaitGroup
	wg.Add(maxChannel)
	for i := 0; i < maxChannel; i++ {
		go func(errs chan error, chanel chan IThreadIndex[T]) {
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
		ch <- datas[i]
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
