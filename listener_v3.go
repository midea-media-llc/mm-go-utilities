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

type ListenerV3[P1 any, P2 any, P3 any] struct {
	events []*func(P1, P2, P3)
	p1     P1
	p2     P2
}

func NewListenerV3[P1 any, P2 any, P3 any](p1 P1, p2 P2) *ListenerV3[P1, P2, P3] {
	result := &ListenerV3[P1, P2, P3]{}
	return result
}

func (this *ListenerV3[P1, P2, P3]) Push(event *func(a P1, b P2, c P3)) {
	if event == nil {
		return
	}

	this.events = append(this.events, event)
}

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

func (this *ListenerV3[P1, P2, P3]) Invoke(a P1, b P2, c P3) {
	for _, e := range this.events {
		event := *e
		event(a, b, c)
	}
}

func (this *ListenerV3[P1, P2, P3]) FindFirstParameter() *P1 {
	return &this.p1
}

func (this *ListenerV3[P1, P2, P3]) FindSecondParameter() *P2 {
	return &this.p2
}

func (this *ListenerV3[P1, P2, P3]) removeAt(i int) {
	this.events = append(this.events[:i], this.events[i+1:]...)
}
