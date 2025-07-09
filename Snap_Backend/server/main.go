package main

import (
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/utils"
	"github.com/gin-gonic/gin"
)


func main(){
 
	// DB IMPORT
	utils.DbConnect()

	router := gin.Default()

	router.GET("/", func(c*gin.Context){
		c.JSON(200, gin.H{
			"msg": "Hello from gin!✨",
		})
	})

	router.Run(":6060")

}