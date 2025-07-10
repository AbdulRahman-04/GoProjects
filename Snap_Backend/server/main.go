package main

import (
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/config"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/routes"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/utils"
	"github.com/gin-gonic/gin"
)


func main(){
 
	// DB IMPORT
	utils.DbConnect()

	// func call 
	public.UserCollect()
	public.AdminCollect()
	private.AllCollect()
	private.SnapCollect()
	private.StoryCollect()

	router := gin.Default()

	router.GET("/", func(c*gin.Context){
		c.JSON(200, gin.H{
			"msg": "Hello from gin!✨",
		})
	})

	// all routes groups 
	routes.PrivateRoutes(router)
	routes.PublicRoutes(router)

	router.Run(fmt.Sprintf(":%d", config.AppConfig.Port))

}