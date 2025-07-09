package middleware

import "github.com/gin-gonic/gin"

// only admins
func OnlyAdmin() gin.HandlerFunc{
	return  func(c*gin.Context){
		role, exists := c.Get("role")
		if !exists || role != "admin"{
			c.JSON(400, gin.H{
				"msg": "Access Denied⚠️, Only Admins can access this route!",
			})
			c.Abort()
			return 
		}
		c.Next()
	}
}


// only users 
func OnlyUsers() gin.HandlerFunc{
	return func(c*gin.Context){
		role, exists := c.Get("role")
		if !exists || role != "user"{
			c.JSON(400, gin.H{
				"msg": "Access Denied⚠️, Only Users can access this route!",
			})
			c.Abort()
			return 
		}
		c.Next()
	}
}

// only user and admins 
func OnlyUserOrAdmin() gin.HandlerFunc{
	return func(c*gin.Context){
		role, exists := c.Get("role")
		if !exists || (role != "admin" && role != "user"){
			c.JSON(400, gin.H{
				"msg": "Access Denied⚠️, Only Users and Admins can access this route!",
			})
			c.Abort()
			return 
		} 
	}
}