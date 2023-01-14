package mongodb

import (
	"GoNews/pkg/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := connect()
	if err != nil {
		log.Fatal(err)
	}
}

func TestMongodb_Posts(t *testing.T) {
	st, _ := connect()
	data, err := st.Posts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestMongodb_New(t *testing.T) {
	post := storage.Post{
		Title:       "post",
		Content:     "my best post",
		AuthorID:    0,
		CreatedAt:   1,
		PublishedAt: 1,
	}

	st, _ := connect()

	err := st.AddPost(post)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMongodb_UpdatePost(t *testing.T) {
	st, _ := connect()
	post := storage.Post{
		Title:       "updated post",
		Content:     "my best updated post",
		AuthorID:    0,
		CreatedAt:   2,
		PublishedAt: 2,
	}

	err := st.UpdatePost(post)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_DeletePost(t *testing.T) {
	st, _ := connect()
	post := storage.Post{
		ID: 0,
	}

	err := st.DeletePost(post)
	if err != nil {
		log.Fatal(err)
	}
}

func connect() (*MongoStorage, error) {
	err := godotenv.Load("./../../../.env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	pass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	server := os.Getenv("MONGO_SERVER_ADDRESS")
	port := os.Getenv("MONGO_PORT")

	if user == "" || pass == "" || port == "" || server == "" {
		os.Exit(1)
	}

	connstr := "mongodb://" + user + ":" + pass + "@" + server + ":" + port + "/?authSource=admin"
	st, err := New(connstr)

	return st, err
}
