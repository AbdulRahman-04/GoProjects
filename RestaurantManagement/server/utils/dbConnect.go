package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient*mongo.Client

func DBconnect() error {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connection establish 
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.DBURI))
	if err != nil {
		fmt.Println(err)
	}

	// test request to db 
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println(err)
	}

	MongoClient = client
	fmt.Println("MONGO DB CONNECTED SUCCESSFULLY!✅")
	return  nil
}