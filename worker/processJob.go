package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Mondal-Prasun/w2w/job"
)

const (
	queueKey     string = "Queue:Jobs"
	jobHashKey   string = "Job:"
	failedJobKey string = "Job:Failed:"
	workerNumber uint8  = 3
)

func (w *Worker) ProcessJobs() {

	jobChan := make(chan job.Job, 5)

	for i := range workerNumber {
		go func(workerId uint8) {
			for job := range jobChan {

				log.Printf("Woker number is :%v \n", workerId)
				w.workOnJob(&job)
			}
		}(i)
	}

	go func() {

		for {
			jobSlice, err := w.Rdb.BLPop(w.Ctx, 0, queueKey).Result()

			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Context canceled.. Stopping server..")
					close(jobChan)
					break
				}
				log.Println("Reddis error.. trying in 2 seconds.. ")
				time.Sleep(2 * time.Second)
				continue
			}

			jobId := jobSlice[1]

			log.Println(jobHashKey + jobId)

			result, err := w.Rdb.Get(w.Ctx, (jobHashKey + jobId)).Bytes()

			if err != nil {
				log.Println("cannot get job from key..." + jobId)
				continue
			}

			var currentJob job.Job

			if err = json.Unmarshal(result, &currentJob); err != nil {
				log.Println("Cannot unmarshal :" + jobId)
				continue
			} else {
				jobChan <- currentJob
			}

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

	log.Println("Setting key : " + jobHashKey + job.JobUniqueId)

	rdbPipe.Set(w.Ctx, (jobHashKey + job.JobUniqueId), payload, 24*time.Hour)

	rdbPipe.LPush(w.Ctx, queueKey, job.JobUniqueId)

	_, err = rdbPipe.Exec(w.Ctx)

}

func (w *Worker) addFiledJob(jobDet *job.Job) {

	w.Rdb.LPush(w.Ctx, failedJobKey, jobDet.JobUniqueId)
	w.changeJobStatus(jobDet, job.Failed)

}

func (w *Worker) changeJobStatus(jobDet *job.Job, status job.JobStatus) {

	jobDet.Status = status

	payLoad, err := json.Marshal(jobDet)

	if err != nil {
		log.Printf("Cannot add failed job: %v \n", err)
		return
	}

	w.Rdb.SetXX(w.Ctx, (jobHashKey + jobDet.JobUniqueId), payLoad, 24*time.Hour)
}

func (w *Worker) JobDetails(jobId string) (*job.Job, error) {

	var reqJob job.Job

	jobBytes, err := w.Rdb.Get(w.Ctx, (jobHashKey + jobId)).Bytes()

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jobBytes, &reqJob); err != nil {
		return nil, err
	}

	return &reqJob, nil

}

func (w *Worker) workOnJob(jobDet *job.Job) {
	// log.Println(job.JobUniqueId + " .....done")

	switch jobDet.JobType {
	case job.Rescale:

		w.generateRescale(jobDet)

	case job.Segmentaion:

		w.generateHLSSegmentation(jobDet)

	case job.Thumbnail:

		w.generateThumbnail(jobDet)

	}

}
