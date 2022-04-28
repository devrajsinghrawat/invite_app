package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/devrajsinghrawat/invite_app/src/db"
	"github.com/devrajsinghrawat/invite_app/src/logging"
	"github.com/devrajsinghrawat/invite_app/src/middleware"
	"github.com/devrajsinghrawat/invite_app/src/model"
	"github.com/devrajsinghrawat/invite_app/src/repo"
	"github.com/devrajsinghrawat/invite_app/src/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tempDir string = "./localdB"

func main() {
	logger := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error("error loading .env file")
	}

	dbCfg := model.Database{
		DataSource: filepath.Join(tempDir, "invite.db"),
		Debug:      len(os.Getenv("DEBUG")) != 0,
		Schema:     "./scripts",
		Type:       "sqlite3",
	}

	conn, err := db.ConnectToDb(dbCfg)
	if err != nil {
		panic(err)
	}
	repoUserProc := repo.NewUserProcessor(conn, logging.GetLogger())
	s := service.NewUserService(repoUserProc)

	repoAppTokenProc := repo.NewAppTokenProcessor(conn)
	at := service.NewAppTokenService(repoAppTokenProc)
	r := mux.NewRouter()
	publicRoute := r.PathPrefix("/public/api/v1").Subrouter()
	publicRoute.Path("/login").HandlerFunc(s.Login).Methods(http.MethodPost)
	publicRoute.Path("/validatetoken/{appToken}").HandlerFunc(at.ValidateToken).Methods(http.MethodGet)
	// rate limiter is used from this link https://stackoverflow.com/questions/60406965/rate-limit-http-requests-based-on-host-address
	publicRoute.Use(middleware.Limit)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.Authenticate)
	api.Path("/genToken").HandlerFunc(at.GenToken).Methods(http.MethodGet)
	api.Path("/getAllToken").HandlerFunc(at.GetAllAppToken).Methods(http.MethodGet)
	api.Path("/invalidateToken").HandlerFunc(at.InvalidateToken).Methods(http.MethodPatch)
	port := os.Getenv("PORT")
	log.Println("Running local on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r))
}
