package apis

import (
	"lottery-engine/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(server Server) error {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API end point
	r.GET("/api/v1/eventinfo/Winners", server.GetEventWinners)
	r.POST("api/v1/event/addWinner", server.addNewWinner)

	return r.Run(server.Addr)
}
