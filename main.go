package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"livebug.dev/remoterun/routers"
	"livebug.dev/remoterun/routers/tasks"
)

// Task 定义任务结构
type Task struct {
	ID        string    `json:"id"`
	Script    string    `json:"script"`
	Args      []string  `json:"args"`
	Status    string    `json:"status"` // pending, running, completed, failed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LogPath   string    `json:"log_path"`
}

// 创建一个内存数据库，用于存储用户和其关联的值
var db = make(map[string]string)

// setupRouter 初始化并配置 Gin 路由器，包含各种路由和中间件。
func setupRouter() *gin.Engine {
	// 创建一个默认的 Gin 路由器实例
	r := gin.Default()

	// 定义一个 GET 路由 /ping，用于测试服务器是否正常运行
	r.GET("/ping", func(c *gin.Context) {
		// 返回字符串 "pong"
		c.String(http.StatusOK, "pong")
	})

	// 定义一个 GET 路由 /user/:name，用于获取指定用户的值
	r.GET("/user/:name", func(c *gin.Context) {
		// 获取 URL 参数中的用户名
		user := c.Params.ByName("name")
		// 从内存数据库中查找该用户的值
		value, ok := db[user]
		if ok {
			// 如果用户存在，返回用户和对应的值
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			// 如果用户不存在，返回状态 "no value"
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// 开发一个api服务，这个服务的功能就是，接受客户端发来的tasks，执行tasks，tasks的主要内容就是执行服务器的某些脚本；并最好能实时返回执行日志；另外服务端对tasks请求落库持久化，能实现管理，对执行日志进行归档处理，路径落库
	// 这里可以使用一个简单的 POST 路由来接收任务
	r.POST("/tasks", func(c *gin.Context) {
		// 定义一个结构体，用于解析 JSON 请求体
		var json struct {
			Task string `json:"task" binding:"required"` // 任务字段，必填
		}
		// 绑定并解析 JSON 数据
		if c.Bind(&json) == nil {
			// 这里可以执行任务，比如调用外部脚本
			// 这里只是模拟执行任务，实际应用中可以使用 os/exec 包来执行命令
			// 例如：cmd := exec.Command("sh", "-c", json.Task)
			// 这里可以将任务的执行结果返回给客户端
			// 这里只是返回一个模拟的执行结果
			c.JSON(http.StatusOK, gin.H{"status": "task executed", "task": json.Task})
		} else {
			// 如果绑定失败，返回错误信息
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid task"})
		}
	})

	// 定义一个受保护的路由组，使用 BasicAuth 中间件
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // 用户名: foo, 密码: bar
		"manu": "123", // 用户名: manu, 密码: 123
	}))

	// 定义一个 POST 路由 /admin，允许授权用户更新其值
	authorized.POST("admin", func(c *gin.Context) {
		// 获取当前授权用户的用户名
		user := c.MustGet(gin.AuthUserKey).(string)

		// 定义一个结构体，用于解析 JSON 请求体
		var json struct {
			Value string `json:"value" binding:"required"` // 值字段，必填
		}

		// 绑定并解析 JSON 数据
		if c.Bind(&json) == nil {
			// 将用户的值存储到内存数据库中
			db[user] = json.Value
			// 返回操作成功的状态
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	// 返回配置完成的路由器实例
	return r
}

func main() {
	// 调用 setupRouter 函数，初始化路由器
	r := routers.Init(tasks.Router)
	// 启动 HTTP 服务器，监听 0.0.0.0:8080
	r.Run(":8080")
}
