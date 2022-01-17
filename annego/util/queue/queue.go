package queue

// Queue 可变长度的循环队列
type Queue struct {
	data []interface{}
	start int
	rear int
	count int
}

func NewQueue() *Queue {
	return &Queue {
		data: make([]interface{}, 4),
		start: 0,
		rear: 0,
		count: 0,
	}
}

func (q *Queue) Len() int {
	return q.count
}

func (q *Queue) PushFront(d interface{}) {
	q.grow()
	q.start = q.prev(q.start)
	q.data[q.start] = d
	q.count++
}

func (q *Queue) PushBack(d interface{}) {
	q.grow()
	q.data[q.rear] = d
	q.rear = q.next(q.rear)
	q.count++
}

func (q *Queue) PopFront() interface{} {
	if q.count == 0 {
		return nil
	}
	r := q.data[q.start]
	q.data[q.start] = nil
	q.start = q.next(q.start)
	q.count--
	q.shrink()
	return r
}

func (q *Queue) PopBack() interface{} {
	if q.count == 0 {
		return nil
	}
	q.rear = q.prev(q.rear)
	r := q.data[q.rear]
	q.data[q.rear] = nil
	q.count--
	q.shrink()
	return r
}

func (q *Queue) Front() interface{} {
	return q.data[q.start]
}

func (q *Queue) Back() interface{} {
	return q.data[q.prev(q.rear)]
}

func (q *Queue) Get(idx int) interface{} {
	if idx >= q.count {
		return nil
	}
	p := (q.start + idx) % len(q.data)
	return q.data[p]
}

func (q *Queue) next(idx int) int {
	return (idx + 1) % len(q.data)
}

func (q *Queue) prev(idx int) int {
	return (len(q.data) + idx - 1) % len(q.data)
}

func (q *Queue) grow() {
	if q.count == len(q.data) {
		q.reserve(q.count * 2)
	}
}

func (q *Queue) shrink() {
	if q.count < len(q.data) / 3 && len(q.data) >= 8 {
		l := len(q.data) / 2
		q.reserve(l)
	}
}

func (q *Queue) reserve(size int) {
	adjusted := make([]interface{}, size)
	if q.start < q.rear {
		copy(adjusted, q.data[q.start:q.rear])
	} else {
		n := copy(adjusted, q.data[q.start:])
		copy(adjusted[n:], q.data[:q.rear])
	}
	q.data = adjusted
	q.start = 0
	q.rear = q.count
}