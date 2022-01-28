package tests

import (
	"fmt"
	"github.com/goal-web/container"
	"github.com/goal-web/pipeline"
	"github.com/pkg/errors"
	"testing"
)

type User struct {
	Id   int
	Name string
}

func TestPipeline(t *testing.T) {
	pipe := pipeline.New(container.New())
	pipe.Send(User{Id: 1, Name: "goal"}).
		Through(
			func(user User, next pipeline.Pipe) interface{} {
				fmt.Println("中间件1-start")
				result := next(user)
				fmt.Println("中间件1-end")
				return result
			},
			func(user User, next pipeline.Pipe) interface{} {
				fmt.Println("中间件2-start")
				result := next(user)
				fmt.Println("中间件2-end")
				return result
			},
		).
		Then(func(user User) {
			fmt.Println("then", user)
		})
}

// 测试异常情况
func TestPipelineException(t *testing.T) {
	defer func() {
		recover()
	}()
	pipe := pipeline.New(container.New())
	pipe.Send(User{Id: 1, Name: "goal"}).
		Through(
			func(user User, next pipeline.Pipe) interface{} {
				fmt.Println("中间件1-start")
				result := next(user)
				fmt.Println("中间件1-end", result)
				return result
			},
			func(user User, next pipeline.Pipe) interface{} {
				fmt.Println("中间件2-start")
				result := next(user)
				fmt.Println("中间件2-end", result)
				return result
			},
		).
		Then(func(user User) {
			panic(errors.New("报个错"))
		})
}
