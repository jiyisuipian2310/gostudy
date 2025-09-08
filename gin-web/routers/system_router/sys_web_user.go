package system_router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SysWebUserRouter struct {
}

// InitSysWebUserRouter 定义初始化用户路由
func (s *SysWebUserRouter) InitSysWebUserRouter(Router *gin.RouterGroup) {
	user := Router.Group("/users")
	user.GET("/userLogin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户登录",
		})
	})
	user.GET("/userRegister", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户注册",
		})
	})

	//路由参数之 :路由： http://localhost:9999/users/1000 (id is 1000)
	user.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "User id is "+id)
	})

	//路由参数之 查询参数： http://localhost:9999/users?name=zhangsan
	user.GET("", func(c *gin.Context) {
		name := c.Query("name") // 获取查询参数name
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取用户信息",
			"name": name,
		})
	})
}
