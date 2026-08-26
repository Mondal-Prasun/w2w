package worker

import (
	"bytes"
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

func (w *Worker) generateHLSSegmentation(jobDet *job.Job) {

	outPutFolder := w.ProcessJobFolder + "/" + jobDet.JobUniqueId

	if err := os.MkdirAll(outPutFolder, os.ModeDir); err != nil {
		log.Println("Cannot make " + outPutFolder)
		w.addFiledJob(jobDet)
		return
	}

	segmentFilePath := fmt.Sprintf("%v/segment_%%03d.ts", outPutFolder)

	outPutM3U8FilePath := fmt.Sprintf("%v/output.m3u8", outPutFolder)

	//

	cmd := exec.Command("ffmpeg",
		"-i", jobDet.JobFileDestination,
		"-map", "0:v:0",
		"-map", "0:a?", // Safe for videos without audio
		"-c:v", "libx264", "-crf", "21", "-preset", "fast",
		"-c:a", "aac", "-b:a", "128k",
		"-g", "48", "-keyint_min", "48", "-sc_threshold", "0",
		"-hls_time", "2",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentFilePath,
		outPutM3U8FilePath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Cannot generate HLS Segmentation: %v | Stderr: %s\n", err, stderr.String())
		w.addFiledJob(jobDet)
		return
	}
	jobDet.JobProcessFolderDestination = outPutFolder

	w.changeJobStatus(jobDet, job.Done)

}

func (w *Worker) generateRescale(jobDet *job.Job) {

	outPutFolder := w.ProcessJobFolder + "/" + jobDet.JobUniqueId

	if err := os.MkdirAll(outPutFolder, os.ModeDir); err != nil {
		log.Println("Cannot make " + outPutFolder)
		w.addFiledJob(jobDet)
		return
	}

	inputFilePath := jobDet.JobFileDestination

	outPutFilePath := fmt.Sprintf("%v/%v.mp4", outPutFolder, jobDet.JobUniqueId)

	cmd := exec.Command("ffmpeg", "-i", inputFilePath,
		"-vf", "scale=1080:1440:force_original_aspect_ratio=decrease,pad=1080:1440:(ow-iw)/2:(oh-ih)/2:black",
		"-c:a", "copy", outPutFilePath)

	if err := cmd.Run(); err != nil {
		log.Printf("Cannot generate Rescale: %v \n", err)

		w.addFiledJob(jobDet)
		return
	}

	jobDet.JobProcessFolderDestination = outPutFolder

	w.changeJobStatus(jobDet, job.Done)

}
