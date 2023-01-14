package mongodb

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "news"
	collectionName = "posts"
)

type MongoStorage struct {
	client *mongo.Client
}

func New(constr string) (*MongoStorage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}

	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	s := MongoStorage{
		client: client,
	}

	return &s, nil
}

func (s *MongoStorage) Posts() ([]storage.Post, error) {
	collection := s.client.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, cur.Err()
}

func (s *MongoStorage) AddPost(p storage.Post) error {
	collection := s.client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)

	return err
}

func (s *MongoStorage) UpdatePost(p storage.Post) error {
	collection := s.client.Database(dbName).Collection(collectionName)

	filter := bson.D{{"id", p.ID}}
	update := bson.D{{"$set", p}}
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (s *MongoStorage) DeletePost(p storage.Post) error {
	collection := s.client.Database(dbName).Collection(collectionName)
	filter := bson.D{{"id", p.ID}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}
