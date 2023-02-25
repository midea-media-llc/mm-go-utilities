package utils

type FinalwareV4[P1 any, P2 any, P3 any, P4 any] func(P1, P2, P3, P4)

type IListenerV4[P1 any, P2 any, P3 any, P4 any] interface {
	Push(*func(P1, P2, P3, P4))
	Pop(*func(P1, P2, P3, P4))
	Invoke(P1, P2, P3, P4, ...FinalwareV4[P1, P2, P3, P4])
}

// ListenerV3 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 and P3 can be any type of parameters
type ListenerV4[P1 any, P2 any, P3 any, P4 any] struct {
	events []*func(P1, P2, P3, P4)
}

// NewListenerV4 returns a new instance of ListenerV3 with initial values of p1 and p2.
func NewListenerV4[P1 any, P2 any, P3 any, P4 any]() IListenerV4[P1, P2, P3, P4] {
	return &ListenerV4[P1, P2, P3, P4]{}
}

// Push adds an event to the event slice
func (l *ListenerV4[P1, P2, P3, P4]) Push(event *func(a P1, b P2, c P3, d P4)) {
	if event == nil {
		return
	}

	l.events = append(l.events, event)
}

// Pop removes an event from the event slice
func (l *ListenerV4[P1, P2, P3, P4]) Pop(event *func(a P1, b P2, c P3, d P4)) {
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

// Invoke calls all events in the event slice with the given parameters
func (l *ListenerV4[P1, P2, P3, P4]) Invoke(a P1, b P2, c P3, d P4, wares ...FinalwareV4[P1, P2, P3, P4]) {
	for _, e := range l.events {
		event := *e
		event(a, b, c, d)
	}

	for _, e := range wares {
		e(a, b, c, d)
	}
}

// removeAt removes an event from the event slice at the given index
func (l *ListenerV4[P1, P2, P3, P4]) removeAt(i int) {
	l.events = RemoveAt(l.events, i)
}
