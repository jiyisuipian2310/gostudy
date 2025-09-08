package system_router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SysWebRoleRouter struct {
}

// InitSysWebRoleRouter 定义初始化用户路由
func (s *SysWebUserRouter) InitSysWebRoleRouter(Router *gin.RouterGroup) {
	role := Router.Group("/roles")
	role.GET("/roleLogin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "角色登录",
		})
	})

	role.GET("/roleRegister", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "角色注册",
		})
	})

}
