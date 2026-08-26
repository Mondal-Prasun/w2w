package worker

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mondal-Prasun/w2w/job"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	Ctx              context.Context
	Rdb              *redis.Client
	UploadFolder     string
	ProcessJobFolder string
}

const (
	maxFromSize int64  = 32 << 10
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

func (w *Worker) CheckRedisServer() error {
	rdbCmd := w.Rdb.Ping(w.Ctx)

	if rdbCmd.Err() != nil {
		return rdbCmd.Err()
	}

	return nil
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

	workerJob := &job.Job{}

	uuid, err := uuid.NewV7()

	if err != nil {
		response.error(502, "Server error")
		log.Println(err)
		return
	}

	workerJob.JobUniqueId = uuid.String()
	workerJob.Status = job.Pending

	switch jobT {
	case string(job.Rescale):
		workerJob.JobType = job.Rescale
	case string(job.Segmentaion):
		workerJob.JobType = job.Segmentaion
	case string(job.Thumbnail):
		workerJob.JobType = job.Thumbnail
	default:
		response.error(303, "Cannot make the request")
		return
	}

	// log.Println(jh.Filename)

	uploadPath, err := w.writeToUploadFolder(jf, jh, workerJob.JobUniqueId)

	if err != nil {
		response.error(502, "Server error")
		log.Println(err)
		return
	}

	workerJob.JobFileDestination = uploadPath

	w.AddJobs(workerJob)

	response.success(200, workerJob)

}

func (w *Worker) GetJobStatus(res http.ResponseWriter, req *http.Request) {
	response := Response{
		w: res,
	}
	if req.Method != http.MethodGet {
		response.error(502, "Bad Gate Way")
		return
	}

	jobId := req.PathValue("jobId")

	if jobId == "" {
		response.error(302, "Cannot find JodId")
		return
	}

	jobDetails, err := w.JobDetails(jobId)

	if err != nil {
		log.Println(err)
		response.error(500, "Internal server error")
		return
	}

	response.success(200, jobDetails)

}

func (w *Worker) writeToUploadFolder(file multipart.File, stat *multipart.FileHeader, jobId string) (path string, uploadErr error) {

	uploadFileDir := fmt.Sprintf("%v/%v", w.UploadFolder, jobId)

	err := os.MkdirAll(uploadFileDir, os.ModeDir)

	if err != nil {
		return "", err
	}

	uploadPath := fmt.Sprintf("%v/%v/%v", w.UploadFolder, jobId, stat.Filename)

	f, err := os.Create(uploadPath)

	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, file)

	if err != nil {
		return "", err
	}

	return uploadPath, nil

}

func (w *Worker) DowloadDoneJob(res http.ResponseWriter, req *http.Request) {
	response := Response{
		w: res,
	}

	if req.Method != http.MethodGet {
		response.error(502, "Bad gateway")
		return
	}

	jobId := req.PathValue("jobId")

	if jobId == "" {
		response.error(302, "cannot find jod Id")
		return
	}

	outPutFolder := w.ProcessJobFolder + "/" + jobId

	if _, err := os.Stat(outPutFolder); os.IsNotExist(err) {
		response.error(302, "cannot find the download folder")
		return
	}

	res.Header().Set("Content-Type", "application/zip")
	res.Header().Set("Content-Disposition", "attachment; filename=\"job_"+jobId+".zip\"")

	zipWriter := zip.NewWriter(res)

	defer zipWriter.Close()

	err := filepath.Walk(outPutFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		relPath, err := filepath.Rel(outPutFolder, path)

		if err != nil {
			return err
		}

		zipFileEntry, err := zipWriter.Create(relPath)

		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileEntry, file)

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		log.Printf("Cannot download file: %v", err)
		return
	}

	if err := os.RemoveAll(outPutFolder); err != nil {
		log.Println("Cannot remove file: " + outPutFolder)
	}

	if err := os.RemoveAll(w.UploadFolder + "/" + jobId); err != nil {
		log.Println("Cannot remove file: " + outPutFolder)
	}

}
