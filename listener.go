package utils

import "reflect"

type Finalware = func()

type Listener struct {
	events []func(args ...interface{})
}

// NewListenerV1 returns a new instance of ListenerV1 with initial values of p1 and p2.
func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Push(event func(args ...interface{})) {
	if event == nil {
		return
	}

	l.events = append(l.events, event)
}

func (l *Listener) Pop(event func(args ...interface{})) {
	if event == nil {
		return
	}

	eventPtr := reflect.ValueOf(event).Pointer()
	for i := 0; i < len(l.events); i++ {
		itemPtr := reflect.ValueOf(l.events[i]).Pointer()
		if itemPtr == eventPtr {
			l.removeAt(i)
			i--
		}
	}
}

func (l *Listener) Invoke(args ...interface{}) {
	for _, e := range l.events {
		event := e
		event(args...)
	}
}

func (l *Listener) removeAt(i int) {
	copy(l.events[i:], l.events[i+1:])
	l.events[len(l.events)-1] = nil
	l.events = l.events[:len(l.events)-1]
}
