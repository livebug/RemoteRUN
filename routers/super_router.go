package routers

import "github.com/gin-gonic/gin"

// 定义一个类型 Register，表示一个函数类型，接收一个 *gin.Engine 参数
type Register func(*gin.Engine)

// Init 函数用于初始化 gin 引擎并注册路由
func Init(routers ...Register) *gin.Engine {
	// 将传入的路由注册函数追加到一个新的切片中
	rs := append([]Register{}, routers...)

	// 创建一个默认的 gin 引擎实例
	r := gin.Default()

	// 遍历所有的路由注册函数，并依次调用它们
	for _, register := range rs {
		register(r) // 调用路由注册函数，将 gin 引擎实例传入
	}

	// 返回初始化完成的 gin 引擎实例
	return r
}
