package utils

type IListenerV1[P1 any] interface {
	Push(event *func(P1))
	Pop(event *func(P1))
	New(p1 P1) IListenerV1[P1]
	Invoke(p1 P1)
	FindFirstParameter() *P1
}

type ListenerV1[P1 any] struct {
	events []*func(P1)
	p1     P1
}

func NewListenerV1[P1 any](p1 P1) *ListenerV1[P1] {
	result := &ListenerV1[P1]{}
	return result
}

func (this *ListenerV1[P1]) Push(event *func(a P1)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

func (this *ListenerV1[P1]) Pop(event *func(a P1)) {
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

func (this *ListenerV1[P1]) Invoke(a P1) {
	for _, e := range this.events {
		event := *e
		event(a)
	}
}

func (this *ListenerV1[P1]) FindFirstParameter() *P1 {
	return &this.p1
}

func (this *ListenerV1[P1]) removeAt(i int) {
	this.events = append(this.events[:i], this.events[i+1:]...)
}
