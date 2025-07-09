package private

import (
	"context"
	"time"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/models"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var storiesCollection *mongo.Collection

func StoryCollect(){
	storiesCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("stories")
}

// create story api 
func CreateStory(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userid 
	userId := c.MustGet("userId").(primitive.ObjectID)
	// form se data lo
	Text := c.PostForm("storytext")
	Song := c.PostForm("storysong")
	Filtername := c.PostForm("filtername")
	Location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	if Text == "" || Song == "" || Filtername == "" || Location == "" || imageUrl == "" {
		c.JSON(400, gin.H{
			"msg": "fill all fields",
		})
		return
	}

	// create new var and push into db 
	var newStory models.Story
	newStory.ID = primitive.NewObjectID()
	newStory.UserId = userId
	newStory.Text = Text
	newStory.Song = Song
	newStory.Location = Location
	newStory.FilterName = Filtername
	newStory.ImageUrl = imageUrl
	newStory.CreatedAt = time.Now()
	newStory.UpdatedAt = time.Now()

	// db m push 
	_, err = storiesCollection.InsertOne(ctx, newStory)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Db error",
		})
		return
	}

	c.JSON(200 , gin.H{
		"msg": "Story CreatedSuccessfully!✅",
	})
}

// func getall stories 
func GetAllStories(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId 
	userId := c.MustGet("userId").(primitive.ObjectID)

	// cursor connection
	cursor, err := storiesCollection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Stories not found or unauthorized",
		})
		return
	}

	defer cursor.Close(ctx)

	var stories []models.Story
	err = cursor.All(ctx, &stories)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "decode error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "All Stories are here✨",
	   "Stories": stories})

}

// get one story 
func GetOneStory(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// user lalo 
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")

	// check kro type param id ka 
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	// find kro user k story ku db m 
	var oneStory models.Story
	err = storiesCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneStory)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "sTORY NOT FOUND",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Your Story is here✨", "Story": oneStory})
}

// edit api 
func EditStory(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userId 
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")

	// check id type 
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "id type wrong",
		})
		return
	}

	// find the collection in db first 
	var updateStory models.Story
	err = storiesCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&updateStory)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	// inpu lelo update k liye 
   // form se data lo
	Text := c.PostForm("storytext")
	Song := c.PostForm("storysong")
	Filtername := c.PostForm("filtername")
	Location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	// update in db 
	update := bson.M{
		"$set": bson.M{
			"storytext": Text,
			"storysong": Song,
			"filtername": Filtername,
			"location": Location,
			"imageUrl": imageUrl,
			"updated_at": time.Now(),
		}}

	// update in db 
	_, err = storiesCollection.UpdateByID(ctx, updateStory.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	} 

	c.JSON(200, gin.H{
			"msg": "Story Updated✨✅",
		})

}

// delet one api 
func DeleteOneStory(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userid lalo 
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param id",
		})
		return
	}

	// delete krdo db m se 
	_ , err = storiesCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Story Deleted Successfully!✅",
	})
}

// delete all api 
func DeleteAllStories(c*gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	// delete all stories 
	_, err := storiesCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No stories found or unauthorized",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "All Stories Deleted successfully!✅",
	})
}