package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mondal-Prasun/w2w/worker"
	"github.com/redis/go-redis/v9"
)

const (
	serverPort      string = "SERVER_ADDR"
	redisPort       string = "REDIS_ADDR"
	uploadDirectory string = "./upload"
)

func main() {

	sPort := os.Getenv(serverPort)
	rAddr := os.Getenv(redisPort)

	checkDirectories()

	httpServerMux := http.NewServeMux()

	worker := worker.Worker{
		Ctx:          context.Background(),
		Rdb:          initRedis(rAddr),
		UploadFolder: uploadDirectory,
	}

	err := worker.CheckRedisServer()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Redis server connected on addr : %v \n", rAddr)

	worker.ProcessJobs()

	httpServerMux.HandleFunc("/checkHealth", worker.Checkhealth)
	httpServerMux.HandleFunc("/acceptJob", worker.AppendJob)

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
	_, err := os.Stat(uploadDirectory)

	if os.IsNotExist(err) {
		fmt.Printf("%v is does not exsist \n", uploadDirectory)
		fmt.Printf("Creating %v directory \n", uploadDirectory)
		err = os.Mkdir(uploadDirectory, os.FileMode(os.O_CREATE))

		if err != nil {
			panic(err)
		}

	}

}
