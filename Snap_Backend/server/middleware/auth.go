package middleware

import (
	"strings"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte(config.AppConfig.JWTKEY)

func AuthMiddleware() gin.HandlerFunc{
	return  func(c*gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(400, gin.H{
				"msg": "No token provided",
			})
			c.Abort()
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts)!= 2 || parts[0] != "Bearer"{
			c.JSON(400, gin.H{
				"msg": "Invalid Token Format",
			})
			c.Abort()
			return 
		}

		myToken := parts[1]

		// token verify
		token, err := jwt.Parse(myToken, func(token*jwt.Token)(interface{}, error){
               return  jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(400, gin.H{
				"msg": "Invalid Token",
			})
			c.Abort()
			return 
		}

		// get data fromt token into claims
		claims , ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "no data found in token",
			})
			c.Abort()
			return 
		}

		// get userId from token
		userStrId, ok := claims["id"].(string)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "no userId data found in token",
			})
			c.Abort()
			return 
		}

		userId, err := primitive.ObjectIDFromHex(userStrId)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "invalid user id format",
			})
			c.Abort()
			return 
		}

		// role nikalo token se 
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "no role data found in token",
			})
			c.Abort()
			return 
		}

		// now set krdo userId, role context variable m
		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()

	}
}