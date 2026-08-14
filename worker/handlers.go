package worker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Mondal-Prasun/w2w/job"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	Ctx context.Context
	Rdb *redis.Client
}

const (
	maxFromSize int64  = 100 * 1024 * 1024
	jobType     string = "jobType"
	jobFile     string = "jobFile"
)

func (w *Worker) Checkhealth(res http.ResponseWriter, req *http.Request) {

	response := &Response{
		w: res,
	}

	if req.Method != http.MethodGet {
		response.error(502, "Bad GateWay")
		return
	}

	response.success(200, struct {
		Msg string `json:"Msg"`
	}{
		Msg: "Server is live ",
	})

}

func (w *Worker) AppendJob(res http.ResponseWriter, req *http.Request) {
	response := Response{w: res}

	if req.Method != http.MethodPost {
		response.error(502, "Bad GateWay")
		return
	}

	if err := req.ParseMultipartForm(maxFromSize); err != nil {
		response.error(402, "Exceded max from size")
		return
	}

	jobT := req.FormValue(jobType)

	if jobT == "" {
		response.error(301, "Job Type is needed")
		return
	}

	jf, jh, err := req.FormFile(jobFile)
	if err != nil {
		response.error(301, "Something went wrong with uploaded file")
		return
	}

	defer jf.Close()

	switch jobT {
	case job.Rescale.ToString():
		fmt.Println("rescale job")
	case job.Segmentaion.ToString():
		fmt.Println("segmentation job")
	case job.Thumbnail.ToString():
		fmt.Println("thumbnail job")
	}

	fmt.Println(jh.Filename)

	response.success(200, struct {
		JobType  string `json:"JobType"`
		FileName string `json:"FileName"`
	}{
		JobType:  jobT,
		FileName: jh.Filename,
	})

}
