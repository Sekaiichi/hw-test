package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length    int
	frontItem *ListItem
	backItem  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.length == 0 {
		l.frontItem = newItem
		l.backItem = newItem
	} else {
		newItem.Next = l.frontItem
		l.frontItem.Prev = newItem
		l.frontItem = newItem
	}
	l.length++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.length == 0 {
		l.frontItem = newItem
		l.backItem = newItem
	} else {
		newItem.Prev = l.backItem
		l.backItem.Next = newItem
		l.backItem = newItem
	}
	l.length++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if l.length == 0 {
		panic("List is empty")
	}
	prev := i.Prev
	next := i.Next

	if i == l.frontItem {
		l.frontItem = next
	} else {
		prev.Next = next
	}

	if i == l.backItem {
		l.backItem = prev
	} else {
		next.Prev = prev
	}

	i.Next, i.Prev = nil, nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.length == 0 {
		panic("List is empty")
	}
	if i == l.frontItem {
		return
	}

	prev := i.Prev
	next := i.Next

	prev.Next = next
	if i == l.backItem {
		l.backItem = prev
	} else {
		next.Prev = prev
	}

	i.Prev = nil
	i.Next = l.frontItem
	l.frontItem.Prev = i

	l.frontItem = i
}

func NewList() List {
	return new(list)
}
