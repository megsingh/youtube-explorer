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

// Storage interface defines methods for interacting with the underlying storage system.
type Storage interface {
	InsertVideos(videos []models.Video) error
	GetPaginatedVideos(token string) (models.PaginationResponse, error)
	SearchVideos(query, token string) (models.PaginationResponse, error)
}

// storage implements the Storage interface
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

// InsertVideos inserts an array of videos into the DB collection.
func (s *storage) InsertVideos(videos []models.Video) error {

	// converting videos to an array of interfaces for insertion.
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

// GetPaginatedVideos retrieves paginated video data from the database.
func (s *storage) GetPaginatedVideos(searchToken string) (models.PaginationResponse, error) {

	// configuring the parameters for aggregated search
	searchByField := "publishedat"
	limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_SEARCH"))
	requiredField := searchByField
	searchIndex := os.Getenv("SEARCH_INDEX_PUBLISH_DATE")
	sortMap := bson.M{
		searchByField: -1,
	}

	// build aggregation pipeline
	pipeline := buildSearchPipeline(searchIndex, "", searchToken, requiredField, nil, sortMap, limit)

	// perform the aggregation
	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.PaginationResponse{}, err
	}
	defer cursor.Close(context.TODO())

	response, err := createPaginatedResponse(cursor, limit)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	return response, nil
}

// SearchVideos searches for videos in the repository based on the query.
func (s *storage) SearchVideos(query string, searchToken string) (models.PaginationResponse, error) {

	limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_SEARCH"))
	searchFields := []string{"title", "description"}
	requiredField := "title"
	searchIndex := os.Getenv("SEARCH_INDEX_TEXT")

	// build a dynamic aggregate pipeline
	pipeline := buildSearchPipeline(searchIndex, query, searchToken, requiredField, searchFields, nil, limit)

	// Perform the aggregation
	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.PaginationResponse{}, err
	}
	defer cursor.Close(context.TODO())

	response, err := createPaginatedResponse(cursor, limit)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	return response, nil
}

// buildSearchPipeline constructs a dynamic aggregate pipeline for search
func buildSearchPipeline(searchIndex, searchQuery, token, requiredField string, searchFields []string, sortMap map[string]interface{}, limit int) []bson.M {

	var pipeline []bson.M

	// define the $search stage of the pipeline
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
	} else {
		searchStage["$search"].(bson.M)["exists"] = bson.M{
			"path": requiredField,
		}
	}

	if token != "" {
		searchStage["$search"].(bson.M)["searchAfter"] = token
	}

	pipeline = append(pipeline, searchStage)

	// define the $sort stage of the pipeline
	if sortMap != nil {
		sortStage := bson.M{
			"$sort": sortMap,
		}
		pipeline = append(pipeline, sortStage)
	}

	// define the $project stage of the pipeline
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

	// define the $limit stage of the pipeline
	limitStage := bson.M{
		"$limit": limit,
	}

	pipeline = append(pipeline, limitStage, projectStage)
	log.Println(pipeline)
	return pipeline
}

// createPaginatedResponse creates a paginated response from the cursor.
func createPaginatedResponse(cursor *mongo.Cursor, limit int) (models.PaginationResponse, error) {

	var videos []models.Video

	no_of_results := 0
	hasNext := true
	nextToken := ""

	//iterate over the response using the cursor and store the videos in the video array
	for cursor.Next(context.TODO()) {
		log.Println("iterating through responses")
		var video models.Video
		err := cursor.Decode(&video)
		if err != nil {
			return models.PaginationResponse{}, err
		}
		videos = append(videos, video)
		no_of_results++
	}

	// check for errors from iterating over the cursor
	if err := cursor.Err(); err != nil {
		return models.PaginationResponse{}, err
	}

	// check whether there are more results
	if no_of_results == 0 || no_of_results < limit {
		hasNext = false
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
