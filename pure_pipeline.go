package pipeline

import (
	"github.com/goal-web/contracts"
)

type PurePipeline struct {
	passable any

	pipes []NextFunc
}

type PureCallback func(stack contracts.Pipe, next NextFunc) contracts.Pipe
type NextFunc func(passable any, pipe contracts.Pipe) any
type PureDestination func(passable any) any

func Pure() *PurePipeline {
	return &PurePipeline{}
}

func (pipeline *PurePipeline) Send(passable any) contracts.Pipeline {
	pipeline.passable = passable
	return pipeline
}

func (pipeline *PurePipeline) SendPure(passable any) *PurePipeline {
	pipeline.passable = passable
	return pipeline
}

func (pipeline *PurePipeline) Through(pipes ...any) contracts.Pipeline {
	for _, item := range pipes {
		pipe, isNextFunc := item.(NextFunc)
		if !isNextFunc {
			panic(PipeArgumentError)
		}
		pipeline.pipes = append(pipeline.pipes, pipe)
	}
	return pipeline
}

func (pipeline *PurePipeline) ThroughPure(pipes ...NextFunc) *PurePipeline {
	pipeline.pipes = append(pipeline.pipes, pipes...)
	return pipeline
}

func (pipeline *PurePipeline) Then(destination any) any {
	return pipeline.ArrayReduce(
		pipeline.reversePipes(), pipeline.carry(), pipeline.prepareDestination(destination),
	)(pipeline.passable)
}

func (pipeline *PurePipeline) ThenPure(destination contracts.Pipe) any {
	return pipeline.ArrayReduce(
		pipeline.reversePipes(), pipeline.carry(), destination,
	)(pipeline.passable)
}

func (pipeline *PurePipeline) carry() PureCallback {
	return func(stack contracts.Pipe, next NextFunc) contracts.Pipe {
		return func(passable any) any {
			return next(passable, stack)
		}
	}
}

func (pipeline *PurePipeline) ArrayReduce(pipes []NextFunc, callback PureCallback, current contracts.Pipe) contracts.Pipe {
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (pipeline *PurePipeline) reversePipes() []NextFunc {
	for from, to := 0, len(pipeline.pipes)-1; from < to; from, to = from+1, to-1 {
		pipeline.pipes[from], pipeline.pipes[to] = pipeline.pipes[to], pipeline.pipes[from]
	}
	return pipeline.pipes
}

func (pipeline *PurePipeline) prepareDestination(destination any) contracts.Pipe {
	pipe, isPipeFunc := destination.(contracts.Pipe)
	if !isPipeFunc {
		panic(PipeArgumentError)
	}
	return pipe
}
