package queue

type Queue[T any] struct {
	data []T
}

func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue[T]) Len() int {
	return len(q.data)
}

func (q *Queue[T]) Push(val T) {
	q.data = append(q.data, val)
}

func (q *Queue[T]) Front() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.data[0], true
}

func (q *Queue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	val := q.data[0]
	q.data = q.data[1:]
	return val, true
}
