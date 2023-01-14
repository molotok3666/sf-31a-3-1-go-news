package postgres

import (
	"GoNews/pkg/storage"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var st *DbStorage

func TestNew(t *testing.T) {
	_, err := connect()
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostgres_Posts(t *testing.T) {
	st, _ := connect()
	data, err := st.Posts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestPostgres_New(t *testing.T) {
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

func TestPostgres_UpdatePost(t *testing.T) {
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

func TestPostgres_DeletePost(t *testing.T) {
	st, _ := connect()
	post := storage.Post{
		ID: 0,
	}

	err := st.DeletePost(post)
	if err != nil {
		log.Fatal(err)
	}
}

func connect() (*DbStorage, error) {
	err := godotenv.Load("./../../../.env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbService := os.Getenv("POSTGRES_DB_SERVICE")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pwd == "" || dbService == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("Empty environment variables")
	}

	connstr := "postgres://" + user + ":" + pwd +
		"@" + dbService + ":" + dbPort + "/" + dbName
	st, err := New(connstr)

	return st, err
}
