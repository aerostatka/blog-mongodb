package structs

import (
	"github.com/aerostatka/mongodb-example/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (p *Post) ToBlogPbPost() *blogpb.Post {
	return &blogpb.Post{
		Id:       p.ID.Hex(),
		Title:    p.Title,
		Content:  p.Content,
		AuthorId: p.AuthorId,
	}
}

func CreatePostFromGrpcObject(post blogpb.Post) *Post {
	storagePost := &Post{
		AuthorId: post.GetAuthorId(),
		Title:    post.GetTitle(),
		Content:  post.GetContent(),
	}

	if post.Id != "" {
		objId, ok := primitive.ObjectIDFromHex(post.Id)
		if ok == nil {
			storagePost.ID = objId
			return storagePost
		} else {
			return nil
		}
	}

	return storagePost
}
