package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongodb"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db1 := memdb.New()
	// Реляционная БД PostgreSQL.
	db2, err := postgresDb()
	if err != nil {
		log.Fatal(err)
	}
	// Документная БД MongoDB.
	db3, err := mongoDb()
	if err != nil {
		log.Fatal(err)
	}
	_, _, _ = db1, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db2

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}

func postgresDb() (*postgres.DbStorage, error) {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbService := os.Getenv("POSTGRES_DB_SERVICE")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pwd == "" || dbService == "" || dbPort == "" || dbName == "" {
		os.Exit(1)
	}

	connstr := "postgres://" + user + ":" + pwd +
		"@" + dbService + ":" + dbPort + "/" + dbName

	return postgres.New(connstr)
}

func mongoDb() (*mongodb.MongoStorage, error) {
	err := godotenv.Load("./../../.env")
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
	return mongodb.New(connstr)
}
