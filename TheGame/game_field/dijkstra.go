package game_field

import (
	"errors"
)

type queue struct {
	qData   chan *gameCell
	inQueue int
	qSize   int
}

func (q *queue) isEmpty() (result bool) {
	return q.inQueue == 0
}

func (q *queue) enqueue(cell *gameCell) (err error) {
	if q.inQueue > 64 {
		return errors.New("queue is full")
	}
	q.qData <- cell
	return
}

func (q *queue) dequeue() (cell *gameCell, err error) {
	if !q.isEmpty() {
		return <-q.qData, nil
	}
	return nil, errors.New("Queue is empty")
}

func CheckCellReachability(currenCell, goalCell *gameCell) {

}
