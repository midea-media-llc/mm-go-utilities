package utils

import "reflect"

type Listener struct {
	events []func(args ...interface{})
}

func (me *Listener) Push(event func(args ...interface{})) {
	if event == nil {
		return
	}

	me.events = append(me.events, event)
}

func (me *Listener) Pop(event func(args ...interface{})) {
	if event == nil {
		return
	}

	eventPtr := reflect.ValueOf(event).Pointer()
	for i := 0; i < len(me.events); i++ {
		itemPtr := reflect.ValueOf(me.events[i]).Pointer()
		if itemPtr == eventPtr {
			me.removeAt(i)
			i--
		}
	}
}

func (me *Listener) Invoke(args ...interface{}) {
	for _, e := range me.events {
		event := e
		event(args...)
	}
}

func (me *Listener) removeAt(i int) {
	copy(me.events[i:], me.events[i+1:])
	me.events[len(me.events)-1] = nil
	me.events = me.events[:len(me.events)-1]
}
