package main

import (
	"context"
	"fmt"
	"github.com/Mondal-Prasun/w2w/worker"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

const (
	serverPort           string = "SERVER_ADDR"
	redisPort            string = "REDIS_ADDR"
	uploadDirectory      string = "./upload"
	processJobsDirectory string = "./processJobs"
)

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
