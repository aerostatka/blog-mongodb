package objects

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerostatka/mongodb-example/blog/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	TestConnection()
	CreatePost(post *structs.Post) (*structs.Post, error)
	ReadPost(id string) (*structs.Post, error)
	DeletePost(id string) (bool, error)
	UpdatePost(post *structs.Post) (*structs.Post, error)
	ListPosts() ([]*structs.Post, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(c *mongo.Client) *MongoRepository {
	collect := c.Database("blog").Collection("posts")
	return &MongoRepository{
		collection: collect,
	}
}

func (r *MongoRepository) TestConnection() {
	fmt.Println("Connection works successfully")
}

func (r *MongoRepository) CreatePost(post *structs.Post) (*structs.Post, error) {
	res, err := r.collection.InsertOne(context.Background(), post)
	if err != nil {
		return post, err
	}

	objId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return post, errors.New("cannot convert to ObjectID")
	}

	post.ID = objId

	return post, nil
}

func (r *MongoRepository) UpdatePost(updatedPost *structs.Post) (*structs.Post, error) {
	post := &structs.Post{}
	filter := bson.M{"_id": updatedPost.ID}
	result := r.collection.FindOne(context.Background(), filter)
	err := result.Decode(post)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	post.AuthorId = updatedPost.AuthorId
	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	_, err = r.collection.ReplaceOne(context.Background(), filter, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *MongoRepository) ReadPost(id string) (*structs.Post, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	post := &structs.Post{}
	filter := bson.M{"_id": objId}

	result := r.collection.FindOne(context.Background(), filter)
	err = result.Decode(post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return post, nil
}

func (r *MongoRepository) DeletePost(id string) (bool, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": objId}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (r *MongoRepository) ListPosts() ([]*structs.Post, error) {
	var posts []*structs.Post
	cursor, err := r.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return posts, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		post := &structs.Post{}
		err = cursor.Decode(post)
		fmt.Println(post)
		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}

	if cursor.Err() != nil {
		return posts, cursor.Err()
	}

	return posts, nil
}
