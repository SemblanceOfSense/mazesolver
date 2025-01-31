package queue

type QueueInterface interface {
    enqueue(int)
    dequeue() int
}

type Queue[T any] []T

func (q Queue[T]) Enqueue(v T) {
    q = append(q, v)
}

func (q Queue[T]) Dequeue() T {
    x := q[0]
    q = q[1:]
    return x
}

func (q Queue[T]) IsEmpty() bool {
    if len(q) == 0 {
        return true
    }
    return false;
}
