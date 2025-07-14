package main

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
)

func main(){

	// DB import
	utils.DBConnect()

	router := gin.Default()

	router.GET("/", func(c*gin.Context){
		c.JSON(200, gin.H{
			"msg": "Hello World From Gin",
		})
	})

	router.Run(":4040")
}