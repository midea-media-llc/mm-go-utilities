package utils

type FinalwareV3[P1 any, P2 any, P3 any] func(P1, P2, P3)

type IListenerV3[P1 any, P2 any, P3 any] interface {
	Push(*func(P1, P2, P3))
	Pop(*func(P1, P2, P3))
	Invoke(P1, P2, P3, ...FinalwareV3[P1, P2, P3])
}

// ListenerV3 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 and P3 can be any type of parameters
type ListenerV3[P1 any, P2 any, P3 any] struct {
	events []*func(P1, P2, P3) // a slice of function pointers that can take P1, P2, P3 parameters
}

// NewListenerV3 returns a new instance of ListenerV3 with initial values of p1 and p2.
func NewListenerV3[P1 any, P2 any, P3 any]() IListenerV3[P1, P2, P3] {
	return &ListenerV3[P1, P2, P3]{}
}

// Push adds an event to the listener's event queue.
func (l *ListenerV3[P1, P2, P3]) Push(event *func(a P1, b P2, c P3)) {
	if event == nil {
		return
	}

	l.events = append(l.events, event)
}

// Pop removes an event from the listener's event queue.
func (l *ListenerV3[P1, P2, P3]) Pop(event *func(a P1, b P2, c P3)) {
	if event == nil {
		return
	}

	for i := 0; i < len(l.events); i++ {
		item := l.events[i]
		if item == event {
			l.removeAt(i)
			i--
		}
	}
}

// Invoke calls all the events in the listener's event queue with the given parameters.
func (l *ListenerV3[P1, P2, P3]) Invoke(a P1, b P2, c P3, wares ...FinalwareV3[P1, P2, P3]) {
	for _, e := range l.events {
		event := *e
		event(a, b, c)
	}

	for _, e := range wares {
		e(a, b, c)
	}
}

// removeAt removes an event at the specified index from the listener's event queue.
func (l *ListenerV3[P1, P2, P3]) removeAt(i int) {
	l.events = RemoveAt(l.events, i)
}
