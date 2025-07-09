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

var snapCollection *mongo.Collection

func SnapCollect() {
	snapCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("snaps")
}

// Create snap api
func CreateSnap(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userid
	userId := c.MustGet("userId").(primitive.ObjectID)
	// form se data lo
	Text := c.PostForm("snaptext")
	Song := c.PostForm("snapsong")
	Filtername := c.PostForm("filtername")
	Location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	if Text == "" || Song == "" || Filtername == "" || Location == "" || imageUrl == "" {
		c.JSON(400, gin.H{"msg": "Required fields missing"})
		return
	}

	// create a new var and push it into db
	var newSnap models.Snap
	newSnap.ID = primitive.NewObjectID()
	newSnap.UserId = userId
	newSnap.Text = Text
	newSnap.Song = Song
	newSnap.FilterName = Filtername
	newSnap.Location = Location
	newSnap.ImageUrl = imageUrl
	newSnap.CreatedAt = time.Now()
	newSnap.UpdatedAt = time.Now()

	

	_, err = snapCollection.InsertOne(ctx, newSnap)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Snap Created Successfully!🎉✅",
	})
}

// Get all snaps
func GetAllSnaps(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	cursor, err := snapCollection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Snaps found for this user⚠️",
		})
		return
	}
	defer cursor.Close(ctx)

	var snaps []models.Snap
	err = cursor.All(ctx, &snaps)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Error while decoding snaps!",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":   "Snaps fetched successfully✅",
		"snaps": snaps,
	})
}

// get one snap
func GetOneSnap(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userId from context
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")

	// check if param id is type of mongoDb id
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid ID format",
		})
		return
	}

	// find the id in db
	var snap models.Snap
	err = snapCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&snap)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Snap Found of this id for this user",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Your Snap Is Here✨", "Snap": snap})
}

// edit snap
func EditSnap(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userid
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")

	// check paramid mongo ud type m h ya nai
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Id format in query",
		})
		return
	}

	// check if id exists or not in db
	// new var
	var updateSnap models.Snap
	err = snapCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&updateSnap)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Snap Found for this user",
		})
		return
	}

	// take input
	Text := c.PostForm("snaptext")
	Song := c.PostForm("snapsong")
	Filtername := c.PostForm("filtername")
	Location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	update := bson.M{
		"$set": bson.M{
			"snaptext":   Text,
			"snapsong":   Song,
			"filtername": Filtername,
			"location":   Location,
			"imageUrl":   imageUrl,
			"updated_at": time.Now(),
		}}
	// update into db
	_, err = snapCollection.UpdateByID(ctx, updateSnap.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Snap Edited Successfully!✅"})

}

// delete one snap api
func DeleteOneSnap(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId lalo
	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")

	// check type of param id
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Param id type",
		})
		return
	}

	// find one and delete
	_, err = snapCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "snap not found or unauthorized",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Snap Deleted Successfully!",
	})
}

// delete all
func DeleteAllSnaps(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId lalo
	userId := c.MustGet("userId").(primitive.ObjectID)

	res, err := snapCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No snaps found or unauthorized",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":     "All Snaps deleted Successfully!😊",
		"deleted": res.DeletedCount,
	})

}
