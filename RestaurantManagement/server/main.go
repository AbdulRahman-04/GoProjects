package main

import (
	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/utils"
	"github.com/gin-gonic/gin"
)


func main(){

	// db import 
	utils.DBconnect()

	router := gin.Default()

	router.GET("/", func (c*gin.Context){
		c.JSON(200, gin.H{
			"msg": "Gin server live at port 6040",
		})
	})

	router.Run(":6040")
}