package route

import (
	"mvt-demo/internal/handler/httpd"
	"mvt-demo/internal/utils/response"
	"github.com/gin-gonic/gin"

	//_ "mvt-demo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(engine *gin.Engine) {

	//404
	engine.NoRoute(func(c *gin.Context) {
		utilGin := response.Gin{Ctx: c}
		utilGin.Response(404, "请求方法不存在", nil)
	})

	api := engine.Group("/test")
	{

		api.GET("/ping", func(c *gin.Context) {
			utilGin := response.Gin{Ctx: c}
			utilGin.Response(1, "pong", nil)
		})
		api.GET("/mvt/:z/:x/:y", httpd.MVTGet)
	}

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
