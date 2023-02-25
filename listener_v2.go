package utils

type FinalwareV2[P1 any, P2 any] func(P1, P2)

type IListenerV2[P1 any, P2 any] interface {
	Push(*func(P1, P2))
	Pop(*func(P1, P2))
	Invoke(P1, P2, ...FinalwareV2[P1, P2])
}

// ListenerV2 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 can be any type of parameters
type ListenerV2[P1 any, P2 any] struct {
	events []*func(P1, P2) // slice of function pointers that accept P1 and P2 as parameters
}

// NewListenerV2 returns a new instance of ListenerV2 with initial values of p1 and p2.
func NewListenerV2[P1 any, P2 any]() IListenerV2[P1, P2] {
	return &ListenerV2[P1, P2]{}
}

// Push adds an event to the events slice
func (l *ListenerV2[P1, P2]) Push(event *func(a P1, b P2)) {
	if event == nil {
		return
	}

	l.events = append(l.events, event)
}

// Pop removes an event from the events slice
func (l *ListenerV2[P1, P2]) Pop(event *func(a P1, b P2)) {
	if event == nil {
		return
	}

	for i, item := range l.events {
		if item == event {
			l.removeAt(i)
			return
		}
	}
}

// Invoke calls all the events in the events slice and passes a and b as parameters to them
func (l *ListenerV2[P1, P2]) Invoke(a P1, b P2, wares ...FinalwareV2[P1, P2]) {
	for _, e := range l.events {
		event := *e
		event(a, b)
	}

	for _, e := range wares {
		e(a, b)
	}
}

// removeAt removes an element at a given index from the events slice
func (l *ListenerV2[P1, P2]) removeAt(i int) {
	l.events = RemoveAt(l.events, i)
}
