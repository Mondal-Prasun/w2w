package job

type JobTypes interface {
	ToString() string
}

type thumbnail struct{}
type segmentation struct{}
type rescale struct{}

func (thumbnail) ToString() string    { return "Thumbnail" }
func (segmentation) ToString() string { return "Segmentation" }
func (rescale) ToString() string      { return "Rescale" }

var (
	Thumbnail   JobTypes = thumbnail{}
	Segmentaion JobTypes = segmentation{}
	Rescale     JobTypes = rescale{}
)

type JobStatus interface {
	ToString() string
}

type done struct{}
type pending struct{}
type failed struct{}

func (done) ToString() string    { return "Done" }
func (pending) ToString() string { return "Pending" }
func (failed) ToString() string  { return "Failed" }

var (
	Done    JobStatus = done{}
	Pending JobStatus = pending{}
	Failed  JobStatus = failed{}
)

type Job struct {
	JobUniqueId        string
	JobFileDestination string
	JobType            JobTypes
	Status             JobStatus
}
