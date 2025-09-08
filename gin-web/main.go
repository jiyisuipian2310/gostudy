package main

import (
	"gin-web/routers/system_router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	group := r.Group("")
	system_router.InitRouters(group)
	r.Run(":9999")
}
