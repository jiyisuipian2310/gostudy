package system_router

import (
	"github.com/gin-gonic/gin"
)

type SysRoutersManager struct {
	SysWebUserRouter // 用户路由
	SysWebRoleRouter // 角色路由
	SysRedirectRouter
}

var SysRoutersMng SysRoutersManager

func InitRouters(Router *gin.RouterGroup) {
	SysRoutersMng.InitSysWebRoleRouter(Router)
	SysRoutersMng.InitSysWebUserRouter(Router)
	SysRoutersMng.InitSysRedirectRouter(Router)
}
