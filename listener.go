package utils

import "reflect"

type Finalware = func()

type Listener struct {
	events    []func(args ...interface{})
	isInvoked bool
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
	if l.isInvoked {
		return
	}

	l.isInvoked = true

	for _, e := range l.events {
		event := e
		event(args...)
	}
}

func (l *Listener) removeAt(i int) {
	l.events = RemoveAt(l.events, i)
}
