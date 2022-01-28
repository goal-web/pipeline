package pipeline

import (
	"errors"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
)

type Pipeline struct {
	container contracts.Container

	passable interface{}

	pipes []contracts.MagicalFunc
}

var PipeArgumentError = errors.New("pipe parameters must have a return value")

type Callback func(stack Pipe, next contracts.MagicalFunc) Pipe
type Pipe func(passable interface{}) interface{}

func New(container contracts.Container) *Pipeline {
	return &Pipeline{
		container: container,
	}
}

func (this *Pipeline) Send(passable interface{}) *Pipeline {
	this.passable = passable
	return this
}

func (this *Pipeline) Through(pipes ...interface{}) *Pipeline {
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

func (this *Pipeline) Then(destination interface{}) interface{} {
	return this.ArrayReduce(
		this.reversePipes(), this.carry(), this.prepareDestination(destination),
	)(this.passable)
}

func (this *Pipeline) carry() Callback {
	return func(stack Pipe, next contracts.MagicalFunc) Pipe {
		return func(passable interface{}) interface{} {
			return this.container.StaticCall(next, passable, stack)[0]
		}
	}
}

func (this *Pipeline) ArrayReduce(pipes []contracts.MagicalFunc, callback Callback, initial Pipe) Pipe {
	var current = initial
	for _, magicalFunc := range pipes {
		current = callback(current, magicalFunc)
	}
	return current
}

func (this *Pipeline) reversePipes() []contracts.MagicalFunc {
	for from, to := 0, len(this.pipes)-1; from < to; from, to = from+1, to-1 {
		this.pipes[from], this.pipes[to] = this.pipes[to], this.pipes[from]
	}
	return this.pipes
}

func (this *Pipeline) prepareDestination(destination interface{}) Pipe {
	return func(passable interface{}) interface{} {
		pipe, isStaticFunc := destination.(contracts.MagicalFunc)
		if !isStaticFunc {
			pipe = container.NewMagicalFunc(destination)
		}
		return this.container.StaticCall(pipe, passable)
	}
}
