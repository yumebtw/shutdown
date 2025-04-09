package domain

import "context"

type Job struct {
	Payload string
}

type Processor interface {
	Process(job Job, ctx context.Context) error
}
