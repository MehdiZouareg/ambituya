package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tuya/tuya-connector-go/example/service"
)

func NewGinEngin() *gin.Engine {
	engine := gin.New()
	initRouter(engine)
	return engine
}

func initRouter(r *gin.Engine) {
	deviceGroup := r.Group("/devices")
	deviceGroup.GET("/:device_id", service.GetDevice)
	deviceGroup.POST("/:device_id/commands", service.PostDeviceCmd)
}
