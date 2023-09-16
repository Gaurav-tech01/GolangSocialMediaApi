package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinit(colName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error Loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DB_NAME")
	// client,err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOption := options.Client().ApplyURI(MongoDb).SetServerAPIOptions(serverAPI)
	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")
	var collection *mongo.Collection = client.Database(dbName).Collection(colName)

	return collection
}
