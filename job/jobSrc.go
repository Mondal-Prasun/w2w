package job

type JobTypes string

var (
	Thumbnail   JobTypes = "Thumbnail"
	Segmentaion JobTypes = "Segmentation"
	Rescale     JobTypes = "rescale"
)

type JobStatus string

var (
	Done    JobStatus = "Done"
	Pending JobStatus = "Pending"
	Failed  JobStatus = "Failed"
)

type Job struct {
	JobUniqueId        string    `json:"JobUniqueId"`
	JobFileDestination string    `json:"JobFileDestination"`
	JobType            JobTypes  `json:"JobType"`
	Status             JobStatus `json:"Status"`
}
