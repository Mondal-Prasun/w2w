package job

type JobTypes string

var (
	Thumbnail   JobTypes = "Thumbnail"
	Segmentaion JobTypes = "Segmentation"
	Rescale     JobTypes = "Rescale"
)

type JobStatus string

var (
	Done    JobStatus = "Done"
	Pending JobStatus = "Pending"
	Failed  JobStatus = "Failed"
)

type Job struct {
	JobUniqueId                 string   `json:"JobUniqueId"`
	JobFileDestination          string   `json:"JobFileDestination"`
	JobProcessFolderDestination string   `json:"JobProcessFolderDestination"`
	JobType                     JobTypes `json:"JobType"`
	JobArgs                     struct {
		InPixel string `json:"InPixel"`
		InRatio string `json:"InRatio"`
	}
	Status JobStatus `json:"Status"`
}
