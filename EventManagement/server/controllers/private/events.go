package private

import (
	"context"
	"strconv"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventsCollection *mongo.Collection

func EventsCollect() {
	eventsCollection = utils.MongoClient.Database("Event_Booking").Collection("events")
}

// create even api
func CreateEvent(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	// take input from form
	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventAttendenceInt, err := strconv.Atoi(eventAttendenceStr)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Conversion error",
		})
		return
	}
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	// var bnake push in db
	var newEvent models.Event

	newEvent.ID = primitive.NewObjectID()
	newEvent.UserId = userId
	newEvent.EventName = eventName
	newEvent.EventtType = eventType
	newEvent.EventAttendence = eventAttendenceInt
	newEvent.EventDescription = eventDes
	newEvent.IsPublic = isPublic
	newEvent.Status = status
	newEvent.Location = location
	newEvent.ImageUrl = imageUrl
	newEvent.CreatedAt = time.Now()
	newEvent.UpdatedAt = time.Now()

	_, err = eventsCollection.InsertOne(ctx, newEvent)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "New Event Created✨", "event Details": newEvent})

}

// get all events 
func GetAllEvents(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	// cursor connection bnao
	cursor, err := eventsCollection.Find(ctx, bson.M{"userId": userId})
    if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	defer cursor.Close(ctx)

	var AllEvents []models.Event
	err = cursor.All(ctx, &AllEvents)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "decoding error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "All Events Are here✨", "Events": AllEvents})
}

// get one event 
func GetOneEvent(c*gin.Context){
		// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// param id 
	paramId := c.Param("id")

	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param Id",
		})
		return
	}

	// userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	// find one event in db
	var oneEvent models.Event
	err = eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneEvent)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No event such found❌",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "One event is here✨", "OneEvent": oneEvent})
}

// edit event api 
func EditEventApi(c*gin.Context){
		// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param Id",
		})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)

	// find it in db 
	var editEvent models.Event
	err = eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&editEvent)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid db error",
		})
		return
	}

	// take inputs 

	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventAttendenceInt, err := strconv.Atoi(eventAttendenceStr)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Conversion error",
		})
		return
	}
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	// var bnake push in db
	// var myEvent models.Event

	//  update db 
	update := bson.M{
		"$set": bson.M{
			"eventname": eventName,
			"eventtype": eventType,
			"attendence": eventAttendenceInt,
			"eventdesc": eventDes,
			"ispublic": isPublic,
			"status": status,
			"location": location,
			"imageUrl": imageUrl,
			"updated_at": time.Now(),
		}}
	// update the db
	_, err = eventsCollection.UpdateByID(ctx, mongoId, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	} 

	c.JSON(200, gin.H{
   "msg": "Event Updated Successfully!✅", "UpdatedEvent": editEvent})
}

// DeleteOne Event Api 
func DeleteOneEvent(c*gin.Context){
		// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param Id",
		})
		return
	}

	// find one id and delete
	// var deleteEvent models.Event
	_, err = eventsCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Event Found or userid not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":"One Event is deleted✅",
	})
}

// delete all events apis 
func DeleteAllEvents(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userid
	userId := c.MustGet("userId").(primitive.ObjectID)

	// find userId and delete all events of it 
	_, err := eventsCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
      c.JSON(400, gin.H{
		"msg": "DB error",
	  })
	  return
	}

	c.JSON(200, gin.H{
		"msg": "All Events Deleted✅",
	})
}