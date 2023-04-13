package pipeline

import (
	"errors"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
)

type Pipeline struct {
	container contracts.Container

	passable any

	pipes []contracts.MagicalFunc
}

var PipeArgumentError = errors.New("pipe parameters must have a return value")

type Callback func(stack contracts.Pipe, next contracts.MagicalFunc) contracts.Pipe

func New(container contracts.Container) contracts.Pipeline {
	return &Pipeline{
		container: container,
	}
}

func (pipeline *Pipeline) Send(passable any) contracts.Pipeline {
	pipeline.passable = passable
	return pipeline
}

func (pipeline *Pipeline) Through(pipes ...any) contracts.Pipeline {
	for _, item := range pipes {
		pipe, isStaticFunc := item.(contracts.MagicalFunc)
		if !isStaticFunc {
			pipe = container.NewMagicalFunc(item)
		}
		if pipe.NumOut() != 1 {
			panic(PipeArgumentError)
		}
		pipeline.pipes = append(pipeline.pipes, pipe)
	}
	return pipeline
}

func (pipeline *Pipeline) Then(destination any) any {
	return pipeline.ArrayReduce(
		pipeline.reversePipes(), pipeline.carry(), pipeline.prepareDestination(destination),
	)(pipeline.passable)
}

func (pipeline *Pipeline) carry() Callback {
	return func(stack contracts.Pipe, next contracts.MagicalFunc) contracts.Pipe {
		return func(passable any) any {
			return pipeline.container.StaticCall(next, passable, stack)[0]
		}
	}
}

func (pipeline *Pipeline) ArrayReduce(pipes []contracts.MagicalFunc, callback Callback, current contracts.Pipe) contracts.Pipe {
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (pipeline *Pipeline) reversePipes() []contracts.MagicalFunc {
	for from, to := 0, len(pipeline.pipes)-1; from < to; from, to = from+1, to-1 {
		pipeline.pipes[from], pipeline.pipes[to] = pipeline.pipes[to], pipeline.pipes[from]
	}
	return pipeline.pipes
}

func (pipeline *Pipeline) prepareDestination(destination any) contracts.Pipe {
	pipe, isStaticFunc := destination.(contracts.MagicalFunc)
	if !isStaticFunc {
		pipe = container.NewMagicalFunc(destination)
	}
	if pipe.NumOut() != 1 {
		panic(PipeArgumentError)
	}
	return func(passable any) any {
		return pipeline.container.StaticCall(pipe, passable)[0]
	}
}
