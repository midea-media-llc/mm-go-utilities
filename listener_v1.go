package utils

type FinalwareV1[P1 any] func(P1)

type IListenerV1[P1 any] interface {
	Push(*func(P1))
	Pop(*func(P1))
	Invoke(P1, ...FinalwareV1[P1])
	InvokeAll(P1, ...Finalware)
}

// ListenerV1 is a generic event listener for one parameter events.
// P1 can be any type of parameters
type ListenerV1[P1 any] struct {
	events    []*func(P1)
	isInvoked bool
}

// NewListenerV1 returns a new instance of ListenerV1 with initial values of p1 and p2.
func NewListenerV1[P1 any]() IListenerV1[P1] {
	return &ListenerV1[P1]{}
}

// Push adds an event to the listener.
func (l *ListenerV1[P1]) Push(event *func(a P1)) {
	if event == nil {
		return
	}
	// Check if the event already exists in the slice.
	for _, e := range l.events {
		if e == event {
			return
		}
	}
	l.events = append(l.events, event)
}

// Pop removes an event from the listener.
func (l *ListenerV1[P1]) Pop(event *func(a P1)) {
	if event == nil {
		return
	}
	// Find the index of the event in the slice.
	for i, e := range l.events {
		if e == event {
			l.removeAt(i)
			return
		}
	}
}

// Invoke calls all the events with the provided parameter.
func (l *ListenerV1[P1]) Invoke(a P1, wares ...FinalwareV1[P1]) {
	if l.isInvoked {
		return
	}

	l.isInvoked = true

	for i := range l.events {
		event := *l.events[i]
		event(a)
	}

	for _, e := range wares {
		e(a)
	}
}

// Invoke calls all the events with the provided parameter.
func (l *ListenerV1[P1]) InvokeAll(a P1, wares ...Finalware) {
	if l.isInvoked {
		return
	}

	l.isInvoked = true

	for i := range l.events {
		event := *l.events[i]
		event(a)
	}

	for _, e := range wares {
		e()
	}
}

// removeAt removes an event at the specified index from the listener's event queue.
func (l *ListenerV1[P1]) removeAt(i int) {
	l.events = RemoveAt(l.events, i)
}
