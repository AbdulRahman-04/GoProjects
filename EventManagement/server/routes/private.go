package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/gin-gonic/gin"
)


func PrivateRoutes(r*gin.Engine){
	privateGroup := r.Group("/api/private")
    privateGroup.Use(middleware.AuthMiddleware())
	 
	{
		// events routes 
		privateGroup.POST("/events/create", middleware.OnlyUsers(), private.CreateEvent)
		privateGroup.GET("/getallevents", middleware.OnlyUsers() ,private.GetAllEvents)
		privateGroup.GET("/getoneevent/:id", middleware.OnlyUsers(),private.GetOneEvent)
		privateGroup.PUT("/updateevent/:id",  middleware.OnlyUsers(),private.EditEventApi)
		privateGroup.DELETE("/deleteoneevent/:id",  middleware.OnlyUsers(),private.DeleteOneEvent)
		privateGroup.DELETE("/deleteallevents", middleware.OnlyUsers(), private.DeleteAllEvents)

		// function routes 
		privateGroup.POST("/func/create", middleware.OnlyUsers(), private.CreateFunction)
		privateGroup.GET("/getallfunc", middleware.OnlyUsers(), private.GetAllFunctions)
		privateGroup.GET("/getonefunc/:id", middleware.OnlyUsers(), private.GetOneFunction)
		privateGroup.PUT("/updatefunc/:id", middleware.OnlyUsers(), private.EditFunction)
		privateGroup.DELETE("/deleteonefunc/:id", middleware.OnlyUsers(), private.DeleteOneFunction)
		privateGroup.DELETE("/deleteallfuncs", middleware.OnlyUsers(), private.DeleteAllFunctions)

		// Admins access routes
		privateGroup.GET("/admins/getallevents", middleware.OnlyAdmins(), private.GetAllEventsAdmin)
		privateGroup.GET("/admins/getone/:id", middleware.OnlyAdmins(), private.GetOneEventAdmin)
		privateGroup.GET("/admins/getallusers", middleware.OnlyAdmins(), private.GetAllUsersAdmin)
		privateGroup.GET("/admins/getoneuser/:id", middleware.OnlyAdmins(), private.GetOneUser)
		privateGroup.GET("/admins/getallfuncs", middleware.OnlyAdmins(), private.GetAllFunctionsAdmin)
		privateGroup.GET("/admins/getonefunc/:id", middleware.OnlyAdmins(), private.GetOneFunctionAdmin)
	}

}