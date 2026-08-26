package worker

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Mondal-Prasun/w2w/job"
)

func (w *Worker) generateThumbnail(jobDet *job.Job) {

	outPutFolder := w.ProcessJobFolder + "/" + jobDet.JobUniqueId

	if err := os.MkdirAll(outPutFolder, os.ModeDir); err != nil {
		log.Println("Cannot make " + outPutFolder)
		w.addFiledJob(jobDet)
		return
	}

	outPutFile := fmt.Sprintf("%v/%v_thumbnail.jpg", outPutFolder, jobDet.JobUniqueId)

	intutFilePath := jobDet.JobFileDestination

	cmd := exec.Command("ffmpeg", "-ss", "00:00:02", "-i", intutFilePath, "-frames:v", "1", "-y", outPutFile)

	if err := cmd.Run(); err != nil {
		log.Printf("Cannot generate thumbnail: %v \n", err)

		w.addFiledJob(jobDet)
		return
	}

	jobDet.JobProcessFolderDestination = outPutFolder

	w.changeJobStatus(jobDet, job.Done)

}
