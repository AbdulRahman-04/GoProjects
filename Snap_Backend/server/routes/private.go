package routes

import (
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/middleware"
	"github.com/gin-gonic/gin"
)


func PrivateRoutes(r*gin.Engine){
	privateGroup := r.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware())


	// snaps route 
	privateGroup.POST("/snap/create", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.CreateSnap)
	privateGroup.GET("/snap/getall",  middleware.AuthMiddleware(), middleware.OnlyUsers(),private.GetAllSnaps)
	privateGroup.GET("/snap/:id", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.GetOneSnap)
	privateGroup.PUT("/editsnap/:id", middleware.AuthMiddleware(),  middleware.OnlyUsers(),private.EditSnap)
	privateGroup.DELETE("/deletesnap/:id",  middleware.AuthMiddleware(), middleware.OnlyUsers(),private.DeleteOneSnap)
	privateGroup.DELETE("/deleteall",  middleware.AuthMiddleware(), middleware.OnlyUsers(),private.DeleteAllSnaps)

	// stories apis 
	privateGroup.POST("/story/create", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.CreateStory)
	privateGroup.GET("/stories/getall", middleware.AuthMiddleware(),  middleware.OnlyUsers(),private.GetAllStories)
	privateGroup.GET("/getonestory/:id", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.GetOneStory)
	privateGroup.PUT("/editonestory/:id", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.EditStory)
	privateGroup.DELETE("/deleteonestory/:id", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.DeleteOneStory)
	privateGroup.DELETE("/deleteallstories", middleware.AuthMiddleware(), middleware.OnlyUsers(),private.DeleteAllStories)
 
	// only admin access apis 
	privateGroup.GET("/admin/getallusers", middleware.AuthMiddleware(), middleware.OnlyAdmin(), private.GetAllUsers)
	privateGroup.GET("/admin/getoneuser/:id", middleware.AuthMiddleware(),middleware.OnlyAdmin(), private.GetOneUserById)
	privateGroup.GET("/admin/getallsnaps", middleware.AuthMiddleware(),middleware.OnlyAdmin(), private.GetAllAdminSnaps)
	privateGroup.GET("/admin/getonesnap/:id", middleware.AuthMiddleware(),middleware.OnlyAdmin(), private.GetOneAdminSnap)
	privateGroup.GET("/admin/getallstories", middleware.AuthMiddleware(),middleware.OnlyAdmin(), private.GetAllAdminStories)
	privateGroup.GET("/admin/getonestory/:id", middleware.AuthMiddleware(),middleware.OnlyAdmin(), private.GetOneAdminStory)

}