package private

import (
	"context"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var functionCollection *mongo.Collection

func FunctionCollect() {
	functionCollection = utils.MongoClient.Database("Event_Booking").Collection("functions")
}

// Create Function
func CreateFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	var newFunction models.Function
	newFunction.ID = primitive.NewObjectID()
	newFunction.UserId = userId
	newFunction.FuncName = funcName
	newFunction.FuncType = funcType
	newFunction.FuncDesc = funcDesc
	newFunction.ImageUrl = imageUrl
	newFunction.IsPublic = isPublic
	newFunction.Status = status
	newFunction.Location = location
	newFunction.CreatedAt = time.Now()
	newFunction.UpdatedAt = time.Now()

	_, err = functionCollection.InsertOne(ctx, newFunction)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "New Function Created✨", "functionDetails": newFunction})
}

// Get All Functions
func GetAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	cursor, err := functionCollection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var allFunctions []models.Function
	err = cursor.All(ctx, &allFunctions)
	if err != nil {
		c.JSON(400, gin.H{"msg": "decoding error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Functions are here✨", "functions": allFunctions})
}

// Get One Function
func GetOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)

	var oneFunction models.Function
	err = functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneFunction)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found❌"})
		return
	}

	c.JSON(200, gin.H{"msg": "One Function is here✨", "function": oneFunction})
}

// Edit Function
func EditFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)

	var oldFunc models.Function
	err = functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oldFunc)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found to update"})
		return
	}

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	update := bson.M{
		"$set": bson.M{
			"funcname":    funcName,
			"functype":    funcType,
			"funcdes":     funcDesc,
			"ispublic":    isPublic,
			"status":      status,
			"location":    location,
			"imageUrl":    imageUrl,
			"updated_at":  time.Now(),
		}}

	_, err = functionCollection.UpdateByID(ctx, oldFunc.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Function Updated Successfully!✅", "updatedFunction": oldFunc})
}

// Delete One Function
func DeleteOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	_, err = functionCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found to delete"})
		return
	}

	c.JSON(200, gin.H{"msg": "One Function Deleted✅"})
}

// Delete All Functions
func DeleteAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	_, err := functionCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Functions Deleted✅"})
}
