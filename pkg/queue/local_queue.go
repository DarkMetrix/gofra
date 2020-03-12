package queue

import (
	"errors"
	"time"
)

var (
	ErrChannelFull = errors.New("Channel full")
	ErrChannelPopTimeout = errors.New("Channel pop timeout")
	ErrChannelClosed = errors.New("Channel closed")
	ErrImproperType = errors.New("interface{} is not this type")
)

// local message queue to buffer message
type LocalMessageQueue struct {
	queueChannel chan interface{}
}

// new local message queue function
func NewLocalMessageQueue(bufferSize uint32) *LocalMessageQueue {
	return &LocalMessageQueue{
		queueChannel: make(chan interface{}, bufferSize),
	}
}

// push message to queue
func (queue *LocalMessageQueue) Push(item interface{}) error {
	select {
	case queue.queueChannel <- item:
		return nil
	default:
		return ErrChannelFull
	}
}

// pop message from queue
func (queue *LocalMessageQueue) Pop(ms time.Duration) (interface{}, error) {
	select {
	case item, ok := <- queue.queueChannel:
		if !ok {
			return nil, ErrChannelClosed
		}

		return item, nil
	case <- time.After(ms):
		return nil, ErrChannelPopTimeout
	}
}

// pop message from queue as []byte
func (queue *LocalMessageQueue) PopAsBytes(ms time.Duration) ([]byte, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return nil ,err
	}

	switch item.(type) {
	case []byte:
		return item.([]byte), nil
	default:
		return nil, ErrImproperType
	}
}

// pop message from queue as string
func (queue *LocalMessageQueue) PopAsString(ms time.Duration) (string, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return "" ,err
	}

	switch item.(type) {
	case string:
		return item.(string), nil
	default:
		return "", ErrImproperType
	}
}

// pop message from queue as int
func (queue *LocalMessageQueue) PopAsInt(ms time.Duration) (int, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return 0 ,err
	}

	switch item.(type) {
	case int:
		return item.(int), nil
	default:
		return 0, ErrImproperType
	}
}

// pop message from queue as int32
func (queue *LocalMessageQueue) PopAsInt32(ms time.Duration) (int32, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return 0 ,err
	}

	switch item.(type) {
	case int32:
		return item.(int32), nil
	default:
		return 0, ErrImproperType
	}
}

// pop message from queue as int64
func (queue *LocalMessageQueue) PopAsInt64(ms time.Duration) (int64, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return 0 ,err
	}

	switch item.(type) {
	case int64:
		return item.(int64), nil
	default:
		return 0, ErrImproperType
	}
}

// pop message from queue as float32
func (queue *LocalMessageQueue) PopAsFloat32(ms time.Duration) (float32, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return 0 ,err
	}

	switch item.(type) {
	case float32:
		return item.(float32), nil
	default:
		return 0, ErrImproperType
	}
}

// pop message from queue as float64
func (queue *LocalMessageQueue) PopAsFloat64(ms time.Duration) (float64, error) {
	item, err := queue.Pop(ms)

	if err != nil {
		return 0 ,err
	}

	switch item.(type) {
	case float64:
		return item.(float64), nil
	default:
		return 0, ErrImproperType
	}
}

// get message channel
func (queue *LocalMessageQueue) Chan() chan interface{} {
	return queue.queueChannel
}