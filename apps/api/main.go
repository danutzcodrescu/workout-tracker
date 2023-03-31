package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	api_repositories "workout-tracker/libs/api/repositories"
	api_utils "workout-tracker/libs/api/utils"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Settings struct {
	db_user     string
	db_password string
	db_host     string
	db_name     string
}

const PORT = 8080

func setupRoutes(app *api_utils.Application) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload-activity", uploadHandler(app))
	return mux
}

func setupDB(user string, password string, host string, dbName string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=verify-full", user, password, dbName, host))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("db connection successful")
	return db
}

func setupEnv() Settings {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err, "failed to load .env")
	}

	return Settings{
		db_user:     os.Getenv("DB_USER"),
		db_password: os.Getenv("DB_PASSWORD"),
		db_host:     os.Getenv("DB_HOST"),
		db_name:     os.Getenv("DB_NAME"),
	}

}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	env := setupEnv()
	db := setupDB(env.db_user, env.db_password, env.db_host, env.db_name)
	app := &api_utils.Application{ErrorLog: errorLog, InfoLog: infoLog, Repositories: api_repositories.SetupRepositories(db)}
	mux := setupRoutes(app)
	log.Printf("Starting server at port %d\n", PORT)
	if err := http.ListenAndServe(":"+fmt.Sprintf("%d", PORT), mux); err != nil {
		log.Fatal(err)
	}
}
