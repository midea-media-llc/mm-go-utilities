package utils

type IListenerV2[P1 any, P2 any] interface {
	Push(event *func(P1, P2))
	Pop(event *func(P1, P2))
	New(p1 P1, p2 P2) IListenerV2[P1, P2]
	Invoke(p1 P1, p2 P2)
	FindFirstParameter() *P1
	FindSecondParameter() *P2
}

type ListenerV2[P1 any, P2 any] struct {
	events []*func(P1, P2)
	p1     P1
	p2     P2
}

func NewListenerV2[P1 any, P2 any](p1 P1, p2 P2) *ListenerV2[P1, P2] {
	result := &ListenerV2[P1, P2]{}
	return result
}

func (this *ListenerV2[P1, P2]) Push(event *func(a P1, b P2)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

func (this *ListenerV2[P1, P2]) Pop(event *func(a P1, b P2)) {
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

func (this *ListenerV2[P1, P2]) Invoke(a P1, b P2) {
	for _, e := range this.events {
		event := *e
		event(a, b)
	}
}

func (this *ListenerV2[P1, P2]) FindFirstParameter() *P1 {
	return &this.p1
}

func (this *ListenerV2[P1, P2]) FindSecondParameter() *P2 {
	return &this.p2
}

func (this *ListenerV2[P1, P2]) removeAt(i int) {
	this.events = append(this.events[:i], this.events[i+1:]...)
}
