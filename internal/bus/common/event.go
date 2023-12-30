package common

import "context"

type Consumer func(context.Context, Event)
type Event struct {
	Topic string
	Body  []byte
}
