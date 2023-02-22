package utils

type IListenerV2[P1 any, P2 any] interface {
	Push(event *func(P1, P2))
	Pop(event *func(P1, P2))
	New(p1 P1, p2 P2) IListenerV2[P1, P2]
	Invoke(p1 P1, p2 P2)
	FindFirstParameter() *P1
	FindSecondParameter() *P2
}

// ListenerV2 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 can be any type of parameters
type ListenerV2[P1 any, P2 any] struct {
	events []*func(P1, P2) // slice of function pointers that accept P1 and P2 as parameters
	p1     P1              // the first parameter
	p2     P2              // the second parameter
}

// Push adds an event to the events slice
func (this *ListenerV2[P1, P2]) Push(event *func(a P1, b P2)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

// Pop removes an event from the events slice
func (this *ListenerV2[P1, P2]) Pop(event *func(a P1, b P2)) {
	if event == nil {
		return
	}

	for i, item := range this.events {
		if item == event {
			this.removeAt(i)
			return
		}
	}
}

// Invoke calls all the events in the events slice and passes a and b as parameters to them
func (this *ListenerV2[P1, P2]) Invoke(a P1, b P2) {
	for _, e := range this.events {
		event := *e
		event(a, b)
	}
}

// FindFirstParameter returns a pointer to the first parameter
func (this *ListenerV2[P1, P2]) FindFirstParameter() *P1 {
	return &this.p1
}

// FindSecondParameter returns a pointer to the second parameter
func (this *ListenerV2[P1, P2]) FindSecondParameter() *P2 {
	return &this.p2
}

// removeAt removes an element at a given index from the events slice
func (this *ListenerV2[P1, P2]) removeAt(i int) {
	copy(this.events[i:], this.events[i+1:])
	this.events[len(this.events)-1] = nil
	this.events = this.events[:len(this.events)-1]
}
