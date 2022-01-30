package tests

import (
	"fmt"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
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
	result := pipe.Send(User{Id: 1, Name: "goal"}).
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
		Then(func(user User) interface{} {
			fmt.Println("then", user)
			return user.Id
		})

	fmt.Println("穿梭结果：", result)
	/**
	中间件1-start
	中间件2-start
	then {1 goal}
	中间件2-end
	中间件1-end
	穿梭结果： 1
	*/
}

// TestPipelineException 测试异常情况
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
	/**
	中间件1-start
	中间件2-start
	*/
}

// TestStaticPipeline 测试调用magical函数
func TestStaticPipeline(t *testing.T) {
	// 应用启动时就准备好的中间件和控制器函数，在大量并发时用 StaticPipeline 可以提高性能
	middlewares := []contracts.MagicalFunc{
		container.NewMagicalFunc(func(user User, next pipeline.Pipe) interface{} {
			fmt.Println("中间件1-start")
			result := next(user)
			fmt.Println("中间件1-end", result)
			return result
		}),
		container.NewMagicalFunc(func(user User, next pipeline.Pipe) interface{} {
			fmt.Println("中间件2-start")
			result := next(user)
			fmt.Println("中间件2-end", result)
			return result
		}),
	}
	controller := container.NewMagicalFunc(func(user User) int {
		fmt.Println("then", user)
		return user.Id
	})

	pipe := pipeline.Static(container.New())
	result := pipe.SendStatic(User{Id: 1, Name: "goal"}).
		ThroughStatic(middlewares...).
		ThenStatic(controller)

	fmt.Println("穿梭结果", result)

	pipe = pipeline.Static(container.New())
	result = pipe.SendStatic(User{Id: 1, Name: "goal"}).
		ThroughStatic(middlewares...).
		ThenStatic(controller)
	fmt.Println("穿梭结果", result)
	/**
	中间件1-start
	中间件2-start
	then {1 goal}
	中间件2-end 1
	中间件1-end 1
	穿梭结果 1
	*/
}

// TestPurePipeline 测试纯净的 pipeline
func TestPurePipeline(t *testing.T) {
	// 如果你的应用场景对性能要求极高，不希望反射影响你，那么你可以试试下面这个纯净的管道
	pipe := pipeline.Pure()
	result := pipe.SendPure(User{Id: 1, Name: "goal"}).
		ThroughPure(
			func(user interface{}, next pipeline.Pipe) interface{} {
				fmt.Println("中间件1-start")
				result := next(user)
				fmt.Println("中间件1-end", result)
				return result
			},
			func(user interface{}, next pipeline.Pipe) interface{} {
				fmt.Println("中间件2-start")
				result := next(user)
				fmt.Println("中间件2-end", result)
				return result
			},
		).
		ThenPure(func(user interface{}) interface{} {
			fmt.Println("then", user)
			return user.(User).Id
		})
	fmt.Println("穿梭结果", result)
	/**
	中间件1-start
	中间件2-start
	then {1 goal}
	中间件2-end 1
	中间件1-end 1
	穿梭结果 1
	*/
}
