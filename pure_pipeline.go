package pipeline

import (
	"github.com/goal-web/contracts"
)

type PurePipeline struct {
	passable interface{}

	pipes []NextFunc
}

type PureCallback func(stack contracts.Pipe, next NextFunc) contracts.Pipe
type NextFunc func(passable interface{}, pipe contracts.Pipe) interface{}
type PureDestination func(passable interface{}) interface{}

func Pure() *PurePipeline {
	return &PurePipeline{}
}

func (this *PurePipeline) Send(passable interface{}) contracts.Pipeline {
	this.passable = passable
	return this
}

func (this *PurePipeline) SendPure(passable interface{}) *PurePipeline {
	this.passable = passable
	return this
}

func (this *PurePipeline) Through(pipes ...interface{}) contracts.Pipeline {
	for _, item := range pipes {
		pipe, isNextFunc := item.(NextFunc)
		if !isNextFunc {
			panic(PipeArgumentError)
		}
		this.pipes = append(this.pipes, pipe)
	}
	return this
}

func (this *PurePipeline) ThroughPure(pipes ...NextFunc) *PurePipeline {
	this.pipes = append(this.pipes, pipes...)
	return this
}

func (this *PurePipeline) Then(destination interface{}) interface{} {
	return this.ArrayReduce(
		this.reversePipes(), this.carry(), this.prepareDestination(destination),
	)(this.passable)
}

func (this *PurePipeline) ThenPure(destination contracts.Pipe) interface{} {
	return this.ArrayReduce(
		this.reversePipes(), this.carry(), destination,
	)(this.passable)
}

func (this *PurePipeline) carry() PureCallback {
	return func(stack contracts.Pipe, next NextFunc) contracts.Pipe {
		return func(passable interface{}) interface{} {
			return next(passable, stack)
		}
	}
}

func (this *PurePipeline) ArrayReduce(pipes []NextFunc, callback PureCallback, current contracts.Pipe) contracts.Pipe {
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (this *PurePipeline) reversePipes() []NextFunc {
	for from, to := 0, len(this.pipes)-1; from < to; from, to = from+1, to-1 {
		this.pipes[from], this.pipes[to] = this.pipes[to], this.pipes[from]
	}
	return this.pipes
}

func (this *PurePipeline) prepareDestination(destination interface{}) contracts.Pipe {
	pipe, isPipeFunc := destination.(contracts.Pipe)
	if !isPipeFunc {
		panic(PipeArgumentError)
	}
	return pipe
}
