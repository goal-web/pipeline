package pipeline

import (
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
)

type StaticPipeline struct {
	container contracts.Container

	passable any

	pipes []contracts.MagicalFunc
}

func Static(container contracts.Container) *StaticPipeline {
	return &StaticPipeline{
		container: container,
	}
}

func (pipeline *StaticPipeline) Send(passable any) contracts.Pipeline {
	pipeline.passable = passable
	return pipeline
}

func (pipeline *StaticPipeline) SendStatic(passable any) *StaticPipeline {
	pipeline.passable = passable
	return pipeline
}

func (pipeline *StaticPipeline) Through(pipes ...any) contracts.Pipeline {
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

func (pipeline *StaticPipeline) ThroughStatic(pipes ...contracts.MagicalFunc) *StaticPipeline {
	pipeline.pipes = append(pipeline.pipes, pipes...)
	return pipeline
}

func (pipeline *StaticPipeline) Then(destination any) any {
	return pipeline.ArrayReduce(
		pipeline.reversePipes(), pipeline.carry(), pipeline.prepareDestination(destination),
	)(pipeline.passable)
}

func (pipeline *StaticPipeline) ThenStatic(destination contracts.MagicalFunc) any {
	return pipeline.ArrayReduce(
		pipeline.reversePipes(), pipeline.carry(), pipeline.prepareStaticDestination(destination),
	)(pipeline.passable)
}

func (pipeline *StaticPipeline) carry() Callback {
	return func(stack contracts.Pipe, next contracts.MagicalFunc) contracts.Pipe {
		return func(passable any) any {
			return pipeline.container.StaticCall(next, passable, stack)[0]
		}
	}
}

func (pipeline *StaticPipeline) ArrayReduce(pipes []contracts.MagicalFunc, callback Callback, current contracts.Pipe) contracts.Pipe {
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (pipeline *StaticPipeline) reversePipes() []contracts.MagicalFunc {
	for from, to := 0, len(pipeline.pipes)-1; from < to; from, to = from+1, to-1 {
		pipeline.pipes[from], pipeline.pipes[to] = pipeline.pipes[to], pipeline.pipes[from]
	}
	return pipeline.pipes
}

func (pipeline *StaticPipeline) prepareDestination(destination any) contracts.Pipe {
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

func (pipeline *StaticPipeline) prepareStaticDestination(destination contracts.MagicalFunc) contracts.Pipe {
	if destination.NumOut() != 1 {
		panic(PipeArgumentError)
	}
	return func(passable any) any {
		return pipeline.container.StaticCall(destination, passable)[0]
	}
}
