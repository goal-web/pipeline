# Goal-web/pipeline
[goal-web/pipeline](https://github.com/goal-web/pipeline)  
这是一个管道库，实现了 和 `laravel` 一样的管道功能，如果你很熟悉 `laravel` 的管道或者中间件，那你一定对这个库很有亲切感。

## 安装 - install
```bash
go get github.com/goal-web/pipeline
```

## 使用 - usage
得益于 goal 强大的容器，你可以在管道(pipe)和目的地(destination)任意注入容器中存在的实例
> 对管道不熟悉的同学，可以把 pipe 理解为中间件，destination 就是控制器方法
> 
```go
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
	/**
	中间件1-start
	中间件2-start
	then {1 goal}
	中间件2-end [1]
	中间件1-end [1]
	穿梭结果 [1]
	 */
}
```

### 在 goal 之外的框架使用 - use in frameworks other than goal
这个库并不会限制你在哪个框架使用它，所以你可以在任意 go 环境使用这个管道库

[goal-web](https://github.com/goal-web/goal)  
[goal-web/pipeline](https://github.com/goal-web/pipeline)  
qbhy0715@qq.com
