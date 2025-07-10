package routes

import (
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/controllers/public"
	"github.com/gin-gonic/gin"
)


func PublicRoutes(r*gin.Engine){
	publicGroup := r.Group("/api/public")


	// users api routing
	publicGroup.POST("/user/signup", public.UserSignup)
	publicGroup.POST("/user/signin", public.UserSignIn)
	publicGroup.GET("/user/emailverify/:token", public.UserEmailVerify)
	publicGroup.POST("/user/change-password", public.UserChangePass)
	publicGroup.POST("/user/forgot-password", public.UserForgotPass)

	// admin apis routing
	publicGroup.POST("/admin/signup", public.AdminSignup)
	publicGroup.POST("/admin/signin", public.AdminSignIn)
	publicGroup.GET("/admin/emailverify/:token", public.EmailVerify)
	publicGroup.POST("/admin/change-password", public.AdminChangePass)
	publicGroup.POST("/admin/forgot-password", public.AdminForgotPass)
}