package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Database *mongo.Database

func ConnectDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://playerchess423_db_user:i3Y4TImkk8cqhvbc@cluster0.z5hbwwd.mongodb.net/?appName=Cluster0&tlsInsecure=true").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	Database = client.Database("boxpark")
	fmt.Println("Connected to MongoDB!")
}
