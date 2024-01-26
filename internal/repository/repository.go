package repository

import (
	"context"
	"log"
	"os"
	"strconv"
	"youtube_project/internal/models"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	InsertVideos(videos []models.Video) error
	GetPaginatedVideos(token string) (models.PaginationResponse, error)
	SearchVideos(query, token string) (models.PaginationResponse, error)
}

type storage struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewStorage(db *mongo.Client, collection *mongo.Collection) Storage {
	return &storage{
		db:         db,
		collection: collection,
	}
}

func (s *storage) InsertVideos(videos []models.Video) error {

	insertInput := make([]interface{}, len(videos))

	for i := range videos {
		insertInput[i] = videos[i]
	}

	_, err := s.collection.InsertMany(context.Background(), insertInput)
	if err != nil {
		return err
	}

	return nil
}

// GetPaginatedVideos retrieves paginated video data from the repository.
func (s *storage) GetPaginatedVideos(searchToken string) (models.PaginationResponse, error) {
	// findOptions := options.Find().SetSort(bson.D{{"publishedat", -1}, {"_id", 1}}) // Sort by descending order of publish_date

	// limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_SEARCH"))
	// findOptions.SetSkip(int64((page - 1) * limit))
	// findOptions.SetLimit(int64(limit))

	// cur, err := s.collection.Find(context.TODO(), bson.D{}, findOptions)
	// if err != nil {
	// 	return nil, err
	// }
	// defer cur.Close(context.TODO())

	// var videos []models.Video
	// for cur.Next(context.TODO()) {
	// 	var video models.Video
	// 	err := cur.Decode(&video)
	// 	if err != nil {
	// 		log.Println("Error decoding video")
	// 		return nil, err
	// 	}
	// 	videos = append(videos, video)
	// }

	// return videos, nil

	limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_SEARCH"))

	requiredField := "publishedat"
	searchIndex := "publishedAt"
	sortMap := bson.M{
		"publishedat": -1,
	}

	// Build a dynamic aggregate pipeline based on your parameters
	pipeline := buildSearchPipeline(searchIndex, "", searchToken, requiredField, nil, sortMap, limit)

	// Perform the aggregation
	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		return models.PaginationResponse{}, err
	}
	defer cursor.Close(context.TODO())

	response, err := createPaginatedResponse(cursor, limit)
	if err != nil {
		log.Fatal(err)
		return models.PaginationResponse{}, err
	}

	return response, nil
}

func (s *storage) SearchVideos(query string, searchToken string) (models.PaginationResponse, error) {

	limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_SEARCH"))

	searchFields := []string{"title", "description"}
	searchIndex := "text_search"

	// Build a dynamic aggregate pipeline based on your parameters
	pipeline := buildSearchPipeline(searchIndex, query, searchToken, "", searchFields, nil, limit)

	// Perform the aggregation
	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		return models.PaginationResponse{}, err
	}
	defer cursor.Close(context.TODO())

	response, err := createPaginatedResponse(cursor, limit)
	if err != nil {
		log.Fatal(err)
		return models.PaginationResponse{}, err
	}

	return response, nil
}

// buildSearchPipeline constructs a dynamic aggregate pipeline for search
func buildSearchPipeline(searchIndex, searchQuery, token, requiredField string, searchFields []string, sortMap map[string]interface{}, limit int) []bson.M {

	var pipeline []bson.M
	// Define the aggregation pipeline
	searchStage := bson.M{
		"$search": bson.M{
			"index": searchIndex,
		},
	}

	if searchQuery != "" {
		searchStage["$search"].(bson.M)["text"] = bson.M{
			"query": searchQuery,
			"path":  searchFields,
		}
	}

	if token != "" {
		searchStage["$search"].(bson.M)["searchAfter"] = token
	}

	if requiredField != "" {
		searchStage["$search"].(bson.M)["exists"] = bson.M{
			"path": requiredField,
		}
	}

	pipeline = append(pipeline, searchStage)

	if sortMap != nil {
		sortStage := bson.M{
			"$sort": sortMap,
		}
		pipeline = append(pipeline, sortStage)
	}

	projectStage := bson.M{
		"$project": bson.M{
			"_id":             1,
			"title":           1,
			"description":     1,
			"publishedat":     1,
			"thumbnailurl":    1,
			"channelId":       1,
			"channelTitle":    1,
			"paginationToken": bson.M{"$meta": "searchSequenceToken"},
		},
	}

	limitStage := bson.M{
		"$limit": limit,
	}

	pipeline = append(pipeline, limitStage, projectStage)
	log.Println(pipeline)
	return pipeline
}

func createPaginatedResponse(cursor *mongo.Cursor, limit int) (models.PaginationResponse, error) {

	var videos []models.Video
	no_of_results := 0
	var hasNext bool = true
	var nextToken string
	for cursor.Next(context.TODO()) {
		var video models.Video
		err := cursor.Decode(&video)
		if err != nil {
			log.Fatal(err)
			return models.PaginationResponse{}, err
		}
		videos = append(videos, video)
		no_of_results++
	}

	// Check for errors from iterating over the cursor
	if err := cursor.Err(); err != nil {
		return models.PaginationResponse{}, err
	}

	if no_of_results == 0 || no_of_results < limit {
		hasNext = false
		nextToken = ""
	} else {
		nextToken = videos[0].PaginationToken
	}

	response := models.PaginationResponse{
		Videos:          videos,
		HasNext:         hasNext,
		PaginationToken: nextToken,
	}

	return response, nil
}
