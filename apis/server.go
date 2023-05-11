package apis

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"lottery-engine/docs"
)

func NewServer(server Server) error {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API end point
	r.GET("/api/v1/winners", server.getEventWinners)
	r.POST("api/v1/winners", server.generateWinners)

	return r.Run(server.Addr)
}
