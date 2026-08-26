package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/Mondal-Prasun/w2w/worker"
	"github.com/redis/go-redis/v9"
)

const (
	serverPort           string = "SERVER_ADDR"
	redisPort            string = "REDIS_ADDR"
	uploadDirectory      string = "./upload"
	processJobsDirectory string = "./processJobs"
	staticJSFolderPath   string = "ui/static"
	templateHTMLFilePath string = "ui/template/index.html"
)

//go:embed ui/*
var uiFS embed.FS

func main() {

	sPort := os.Getenv(serverPort)
	rAddr := os.Getenv(redisPort)

	checkDirectories()

	httpServerMux := http.NewServeMux()

	worker := worker.Worker{
		Ctx:              context.Background(),
		Rdb:              initRedis(rAddr),
		UploadFolder:     uploadDirectory,
		ProcessJobFolder: processJobsDirectory,
	}

	err := worker.CheckRedisServer()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Redis server connected on addr : %v \n", rAddr)

	worker.ProcessJobs()

	tmpl, err := template.ParseFS(uiFS, templateHTMLFilePath)

	if err != nil {
		log.Panicf("Cannot parse template : %v", err)
	}

	staticFs, err := fs.Sub(uiFS, staticJSFolderPath)

	if err != nil {
		log.Panicf("Cannot parse static folder : %v", err)
	}

	httpServerMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFs))))

	httpServerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	httpServerMux.HandleFunc("/checkHealth", worker.Checkhealth)
	httpServerMux.HandleFunc("/acceptJob", worker.AppendJob)
	httpServerMux.HandleFunc("/getJobDetail/{jobId}", worker.GetJobStatus)
	httpServerMux.HandleFunc("/download/{jobId}", worker.DowloadDoneJob)

	log.Printf("Server is listening on %v \n", sPort)
	http.ListenAndServe(sPort, httpServerMux)

}

func initRedis(redisAddr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}

func checkDirectories() {

	directories := []string{uploadDirectory, processJobsDirectory}

	for _, dir := range directories {

		_, err := os.Stat(dir)

		if os.IsNotExist(err) {
			fmt.Printf("%v is does not exsist \n", dir)
			fmt.Printf("Creating %v directory \n", dir)
			err = os.Mkdir(dir, os.ModeDir)
			if err != nil {
				panic(err)
			}

		}
	}

}
