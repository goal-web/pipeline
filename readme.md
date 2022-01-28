# Goal-web/pipeline
[goal-web/pipeline](https://github.com/goal-web/pipeline)  
这是一个管道库，实现了 和 `laravel` 一样的管道功能，如果你很熟悉 `laravel` 的管道或者中间件，那你一定对这个库很有亲切感。

## 安装 - install
```bash
go get github.com/goal-web/pipeline
```

## 使用 - usage
```go
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

```

### 在 goal 之外的框架使用 - use in frameworks other than goal
这个库并不会限制你在哪个框架使用它，所以你可以在任意 go 环境使用这个管道库

[goal-web](https://github.com/goal-web/goal)  
[goal-web/pipeline](https://github.com/goal-web/pipeline)  
qbhy0715@qq.com
