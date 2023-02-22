package utils

type IListenerV3[P1 any, P2 any, P3 any] interface {
	Push(event *func(P1, P2, P3))
	Pop(event *func(P1, P2, P3))
	New(p1 P1, p2 P2, p3 P3) IListenerV3[P1, P2, P3]
	Invoke(p1 P1, p2 P2, p3 P3)
	FindFirstParameter() *P1
	FindSecondParameter() *P2
	FindThirdParameter() *P3
}

// ListenerV3 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 and P3 can be any type of parameters
type ListenerV3[P1 any, P2 any, P3 any] struct {
	events []*func(P1, P2, P3) // a slice of function pointers that can take P1, P2, P3 parameters
	p1     P1                  // a placeholder for P1
	p2     P2                  // a placeholder for P2
}

// NewListenerV3 returns a new instance of ListenerV3 with initial values of p1 and p2.
func NewListenerV3[P1 any, P2 any, P3 any](p1 P1, p2 P2) *ListenerV3[P1, P2, P3] {
	result := &ListenerV3[P1, P2, P3]{
		p1: p1,
		p2: p2,
	}
	return result
}

// Push adds an event to the listener's event queue.
func (this *ListenerV3[P1, P2, P3]) Push(event *func(a P1, b P2, c P3)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

// Pop removes an event from the listener's event queue.
func (this *ListenerV3[P1, P2, P3]) Pop(event *func(a P1, b P2, c P3)) {
	if event == nil {
		return
	}

	for i := 0; i < len(this.events); i++ {
		item := this.events[i]
		if item == event {
			this.removeAt(i)
			i--
		}
	}
}

// Invoke calls all the events in the listener's event queue with the given parameters.
func (this *ListenerV3[P1, P2, P3]) Invoke(a P1, b P2, c P3) {
	for _, e := range this.events {
		event := *e
		event(a, b, c)
	}
}

// FindFirstParameter returns a pointer to the listener's first parameter placeholder.
func (this *ListenerV3[P1, P2, P3]) FindFirstParameter() *P1 {
	return &this.p1
}

// FindSecondParameter returns a pointer to the listener's second parameter placeholder.
func (this *ListenerV3[P1, P2, P3]) FindSecondParameter() *P2 {
	return &this.p2
}

// removeAt removes an event at the specified index from the listener's event queue.
func (this *ListenerV3[P1, P2, P3]) removeAt(i int) {
	copy(this.events[i:], this.events[i+1:])
	this.events[len(this.events)-1] = nil
	this.events = this.events[:len(this.events)-1]
}
