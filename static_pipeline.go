package pipeline

import (
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
)

type StaticPipeline struct {
	container contracts.Container

	passable interface{}

	pipes []contracts.MagicalFunc
}

func Static(container contracts.Container) *StaticPipeline {
	return &StaticPipeline{
		container: container,
	}
}

func (this *StaticPipeline) Send(passable interface{}) contracts.Pipeline {
	this.passable = passable
	return this
}

func (this *StaticPipeline) SendStatic(passable interface{}) *StaticPipeline {
	this.passable = passable
	return this
}

func (this *StaticPipeline) Through(pipes ...interface{}) contracts.Pipeline {
	for _, item := range pipes {
		pipe, isStaticFunc := item.(contracts.MagicalFunc)
		if !isStaticFunc {
			pipe = container.NewMagicalFunc(item)
		}
		if pipe.NumOut() != 1 {
			panic(PipeArgumentError)
		}
		this.pipes = append(this.pipes, pipe)
	}
	return this
}

func (this *StaticPipeline) ThroughStatic(pipes ...contracts.MagicalFunc) *StaticPipeline {
	this.pipes = append(this.pipes, pipes...)
	return this
}

func (this *StaticPipeline) Then(destination interface{}) interface{} {
	return this.ArrayReduce(
		this.reversePipes(), this.carry(), this.prepareDestination(destination),
	)(this.passable)
}

func (this *StaticPipeline) ThenStatic(destination contracts.MagicalFunc) interface{} {
	return this.ArrayReduce(
		this.reversePipes(), this.carry(), this.prepareStaticDestination(destination),
	)(this.passable)
}

func (this *StaticPipeline) carry() Callback {
	return func(stack Pipe, next contracts.MagicalFunc) Pipe {
		return func(passable interface{}) interface{} {
			return this.container.StaticCall(next, passable, stack)[0]
		}
	}
}

func (this *StaticPipeline) ArrayReduce(pipes []contracts.MagicalFunc, callback Callback, current Pipe) Pipe {
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (this *StaticPipeline) reversePipes() []contracts.MagicalFunc {
	for from, to := 0, len(this.pipes)-1; from < to; from, to = from+1, to-1 {
		this.pipes[from], this.pipes[to] = this.pipes[to], this.pipes[from]
	}
	return this.pipes
}

func (this *StaticPipeline) prepareDestination(destination interface{}) Pipe {
	return func(passable interface{}) interface{} {
		pipe, isStaticFunc := destination.(contracts.MagicalFunc)
		if !isStaticFunc {
			pipe = container.NewMagicalFunc(destination)
		}
		return this.container.StaticCall(pipe, passable)
	}
}

func (this *StaticPipeline) prepareStaticDestination(destination contracts.MagicalFunc) Pipe {
	return func(passable interface{}) interface{} {
		return this.container.StaticCall(destination, passable)
	}
}
