package private

import (
	"context"
	"time"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/models"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection
var mySnapCollection *mongo.Collection
var myStoriesCollection * mongo.Collection

func AllCollect(){

	userCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("users")
	mySnapCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("snaps")
   myStoriesCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("stories")
}

// func getall users and snaps and stories 

func GetAllUsers(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Users Found!",
		})
		return
	}

	defer cursor.Close(ctx)

	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "error while decoding",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "All Users Data","Data" : users})
}

// get one user by id 
func GetOneUserById(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId,err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid id",
		})
		return
	}

	// find id in db 
	var oneUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneUser)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "no users found!",
		})
		return
	}

	c.JSON(200, gin.H{

		"msg": "One user is here", "User": oneUser})
}


// GET ALL SNAPS
func GetAllAdminSnaps(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mySnapCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No snaps found!",
		})
		return
	}

	defer cursor.Close(ctx)

	var snaps []models.Snap
	err = cursor.All(ctx, &snaps)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "error decoding",
		})
		return
	}

	c.JSON(200, gin.H{
			"msg": "All Snaps✨", "Snaps": snaps})

}

// Get one snap 
func GetOneAdminSnap(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "error id",
		})
		return
	}

	// find in db 
	var oneSnap models.Snap
	err = mySnapCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneSnap)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "no snap found",
		})
		return
	}

	c.JSON(200, gin.H{
			"msg": "one snap is here✨","Snap": oneSnap})
		
}


// GET ALL ADMIN STORIES
func GetAllAdminStories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := myStoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No stories found!",
		})
		return
	}

	defer cursor.Close(ctx)

	var stories []models.Story
	err = cursor.All(ctx, &stories)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "error decoding",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":    "All Stories✨",
		"Stories": stories,
	})
}

// GET ONE ADMIN STORY BY ID
func GetOneAdminStory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid ID",
		})
		return
	}

	var oneStory models.Story
	err = myStoriesCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneStory)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No story found with this ID",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":   "One story is here✨",
		"Story": oneStory,
	})
}
