package main

import "github.com/gin-gonic/gin"


func main(){
	router := gin.Default()

	router.GET("/", func (c*gin.Context){
		c.JSON(200, gin.H{
			"msg": "Gin server live at port 6040",
		})
	})

	router.Run(":6040")
}