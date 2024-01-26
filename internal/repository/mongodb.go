package repository

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDatabase establishes a connection to MongoDB and returns the mongo client
func SetupDatabase() (*mongo.Client, *mongo.Collection, error) {

	// setting up MongoDb connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// create a new client
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{bson.E{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, nil, err
	}

	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))

	// We can create MongoDB "indexes" using the the below
	// I have decided the create Atlas "Search Index" as it is more efiicient and functional than normal db index

	// collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	// _, err = collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
	// 	{
	// 		Keys:    bson.D{{"publish_date", -1}}, // Descending order
	// 		Options: options.Index().SetName("publish_date_index"),
	// 	},
	// 	{
	// 		Keys:    bson.D{{"title", "text"}, {"description", "text"}},
	// 		Options: options.Index().SetName("text_search_index"),
	// 	},
	// })

	return client, collection, nil

}
