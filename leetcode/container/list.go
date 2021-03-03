package container

// List to store
type List struct {
	Head   *Element
	Tail   *Element
	Length int
}

// NewList new list
func NewList() *List {
	return &List{}
}

// PushBack push one at back
func (l *List) PushBack(v interface{}) *Element {
	e := &Element{
		Prev:  l.Tail,
		Next:  nil,
		Value: v,
	}
	l.Tail = e
	if l.Head == nil {
		l.Head = e
	}
	if l.Tail == nil {
		l.Tail = e
	}
	l.Length++
	return e
}

// PushFront push one at front
func (l *List) PushFront(v interface{}) *Element {
	e := &Element{
		Prev:  nil,
		Next:  l.Head,
		Value: v,
	}
	if l.Head != nil {
		l.Head.Prev = e
	}
	l.Head = e
	if l.Tail == nil {
		l.Tail = e
	}
	l.Length++
	return e
}

// Remove element
func (l *List) Remove(e *Element) {
	if e == nil {
		return
	}
	if e.Prev != nil {
		e.Prev = e.Next
	}
	if e.Next != nil {
		e.Next.Prev = e.Prev
	}
	if e == l.Head {
		l.Head = l.Head.Next
	}
	if e == l.Tail {
		l.Tail = l.Tail.Prev
	}
	l.Length--
}

// SetFront set element as front
func (l *List) SetFront(e *Element) {
	if e.Prev == nil {
		return
	}
	e.Prev.Next = e.Next
	if e.Next != nil {
		e.Next.Prev = e.Prev
	} else {
		l.Tail = e.Prev
	}
	e.Prev = nil
	e.Next = l.Head
	l.Head.Prev = e
}

// Element of list
type Element struct {
	Prev  *Element
	Next  *Element
	Value interface{}
}
