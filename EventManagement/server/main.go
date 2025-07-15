package main

import (
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/routes"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	// DB import
	utils.DBConnect()

	router := gin.Default()

	r := gin.Default()
	r.Use(middleware.CustomLogger())

	// function calls
	public.UserCollect()
	public.AdminCollect()
	private.UserAccessCollect()
	private.AdminAccessCollect()
	private.EventsCollect()
	private.FunctionCollect()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "Hello World From Gin",
		})
	})

	// routes register
	routes.PublicRoutes(router)
	routes.PrivateRoutes(router)

	router.Run(fmt.Sprintf(":%d", config.AppConfig.Port))
}
