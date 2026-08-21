package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Mondal-Prasun/w2w/job"
)

const (
	queueKey   string = "Queue:Jobs"
	jobHashKey string = "Job:"
)

func (w *Worker) ProcessJobs() {

	go func() {
		for {
			jobSlice, err := w.Rdb.BLPop(w.Ctx, 0, queueKey).Result()

			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Context canceled.. Stopping server..")
					break
				}
				log.Println("Reddis error.. trying in 2 seconds.. ")
				time.Sleep(2 * time.Second)
				continue
			}

			fmt.Println(jobSlice[1])

		}
	}()

}

func (w *Worker) AddJobs(job *job.Job) {
	rdbPipe := w.Rdb.Pipeline()

	payload, err := json.Marshal(job)

	if err != nil {
		log.Println("Cannot marshal Job struct...")
		return
	}

	rdbPipe.HSet(w.Ctx, (jobHashKey + job.JobUniqueId), payload)
	rdbPipe.Expire(w.Ctx, (jobHashKey + job.JobUniqueId), 24*time.Hour)

	rdbPipe.LPush(w.Ctx, queueKey, job.JobUniqueId)

	_, err = rdbPipe.Exec(w.Ctx)

}
