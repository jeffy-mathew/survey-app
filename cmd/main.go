// main package contains the driver code for running the application
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"survey-platform/internal/app"
	"survey-platform/internal/db/jsondb"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories/responserepo"
	"survey-platform/internal/repositories/surveyrepo"
	"survey-platform/internal/services/surveyservice"
	"syscall"
	"time"
)

const (
	AppPortEnv  = "APP_PORT"
	AppPort     = ":8000"
	surveyAppDB = "survey_app.json"
)

// serve handles the logic of running  server in a goroutine and waiting for signal to gracefully stop the server
// on ctx.Done signal a request to shut down the server is sent, so that no new requests will be served
// after that the data is dumped to the file
func serve(ctx context.Context, surveyApp *app.SurveyApp) {
	router := surveyApp.SetupRoutes()
	port := os.Getenv(AppPortEnv)
	if port == "" {
		port = AppPort
	}
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	log.Printf("server started on port %s", port)

	<-ctx.Done()

	log.Printf("graceful shutdown request received")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%s", err.Error())
	}
	log.Println("application stopped accepting requests, dumping data")
	if err := surveyApp.Dump(); err != nil {
		log.Fatalln("dumping data failed", err.Error())
	}
	log.Println("dumping data complete. app exiting!!")
}

// main initiates new app and calls serve to start the server
// it also spawns a goroutine to listen to os signals SIGINT or SIGTERM
// once the os signal is received the cancel func of ctx passed to serve is called
// notifying it to initiate a graceful shutdown
func main() {
	jsonDB, err := jsondb.NewJsonDB(surveyAppDB)
	if err != nil {
		log.Fatalln("error while initiating db")
	}
	var dbEntry = models.DBEntry{}
	err = jsonDB.Load(&dbEntry)
	if err != nil {
		log.Fatalln("error while loading persisted entries")
	}
	surveyRepo := surveyrepo.NewSurveyRepo(dbEntry.Surveys)
	responseRepo := responserepo.NewResponseRepo(dbEntry.Responses)
	surveyService := surveyservice.NewSurveyService(3, surveyRepo, responseRepo)
	surveyApp := app.NewSurveyApp(jsonDB, surveyService)
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovering from panic, dumping data")
			dumpErr := surveyApp.Dump()
			if dumpErr != nil {
				log.Fatalln("dumping data failed", dumpErr.Error())
			}
			log.Println("dumping data complete. app exiting!!")
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		log.Printf("system call received")
		cancel()
	}()
	serve(ctx, surveyApp)
}
