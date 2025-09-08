package system_router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SysRedirectRouter struct {
}

func costTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("begin cost time\n")
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("===> the request URL %s cost %v\n", url, costTime)
		fmt.Printf("end cost time\n")
	}
}

func (s *SysRedirectRouter) InitSysRedirectRouter(Router *gin.RouterGroup) {
	role := Router.Group("")
	role.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/redirect_test")
	})

	// 处理 /redirect_test
	role.GET("/redirect_test", costTime(), func(c *gin.Context) {
		fmt.Printf("begin cost redirect_test\n")
		time.Sleep(2000 * time.Millisecond) // 模拟业务处理
		c.JSON(http.StatusOK, gin.H{
			"message": "This is a redirected page",
		})
		fmt.Printf("end cost redirect_test\n")
	})
}
