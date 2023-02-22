package utils

type IListenerV4[P1 any, P2 any, P3 any, P4 any] interface {
	Push(event *func(P1, P2, P3, P4))
	Pop(event *func(P1, P2, P3, P4))
	New(p1 P1, p2 P2, p3 P3, p4 P4) IListenerV4[P1, P2, P3, P4]
	Invoke(p1 P1, p2 P2, p3 P3, p4 P4)
	FindFirstParameter() *P1
	FindSecondParameter() *P2
	FindThirdParameter() *P3
	FindFourthParameter() *P4
}

// ListenerV3 is a generic type that represents an event listener that accepts two parameters
// P1 and P2 and P3 can be any type of parameters
type ListenerV4[P1 any, P2 any, P3 any, P4 any] struct {
	events []*func(P1, P2, P3, P4)
	p1     P1
	p2     P2
	p3     P3
}

// Push adds an event to the event slice
func (this *ListenerV4[P1, P2, P3, P4]) Push(event *func(a P1, b P2, c P3, d P4)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

// Pop removes an event from the event slice
func (this *ListenerV4[P1, P2, P3, P4]) Pop(event *func(a P1, b P2, c P3, d P4)) {
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

// Invoke calls all events in the event slice with the given parameters
func (this *ListenerV4[P1, P2, P3, P4]) Invoke(a P1, b P2, c P3, d P4) {
	for _, e := range this.events {
		event := *e
		event(a, b, c, d)
	}
}

// FindFirstParameter returns a pointer to the first parameter
func (this *ListenerV4[P1, P2, P3, P4]) FindFirstParameter() *P1 {
	return &this.p1
}

// FindSecondParameter returns a pointer to the second parameter
func (this *ListenerV4[P1, P2, P3, P4]) FindSecondParameter() *P2 {
	return &this.p2
}

// FindThirdParameter returns a pointer to the third parameter
func (this *ListenerV4[P1, P2, P3, P4]) FindThirdParameter() *P3 {
	return &this.p3
}

// removeAt removes an event from the event slice at the given index
func (this *ListenerV4[P1, P2, P3, P4]) removeAt(i int) {
	copy(this.events[i:], this.events[i+1:])
	this.events[len(this.events)-1] = nil
	this.events = this.events[:len(this.events)-1]
}
