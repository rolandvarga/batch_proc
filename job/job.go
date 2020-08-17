package job

import (
	"fmt"
	"sort"
)

var (
	batchSize = 100
)

// Object represents a single element in our input data.
type Object struct {
	ID   string `json:"id"`
	Seq  int64  `json:"seq"`
	Data string `json:"data"`
}

// GenerateGroupID returns the groupID that each object is grouped under
// when assigned to a worker.
func (o *Object) GenerateGroupID() int64 {
	return o.Seq / 100
}

// Data represents a collection of objects from our input data.
type Data struct {
	Objects []Object `json:"objects"`
}

// Job defines the operation responsible for processing our data set.
type Job struct {
	Data    Data
	Workers []Worker
}

// NewJob initializes a new Job instance
func NewJob(data Data) Job {
	return Job{Data: data}
}

// FindWorkerWithID checks to see if a worker has been already spawned for
// the current job. Returns its index if so.
func (j *Job) FindWorkerWithID(groupID int64) int {
	for i, w := range j.Workers {
		if w.GroupID == groupID {
			return i
		}
	}
	return -1
}

// Add new Worker to Job
func (j *Job) Add(worker Worker) {
	j.Workers = append(j.Workers, worker)
}

// Run starts processing the dataset found in Job
func (j *Job) Run() {
	for _, obj := range j.Data.Objects {
		groupID := obj.GenerateGroupID()

		idx := j.FindWorkerWithID(groupID)
		if idx == -1 {
			worker := NewWorker(groupID)
			j.Add(*worker)
			idx = len(j.Workers) - 1
		}

		j.Workers[idx].Add(obj)

		if j.Workers[idx].Size() >= batchSize {
			j.Workers[idx].SortObjects()
			j.Workers[idx].Process()
		}
	}
}

// Worker is responsible for holding and processing a single batch of objects.
type Worker struct {
	GroupID int64
	Objects []Object
}

// NewWorker initializes a new Worker instance.
func NewWorker(groupID int64) *Worker {
	return &Worker{GroupID: groupID}
}

// Process all Objects in worker.
func (w *Worker) Process() {
	for _, obj := range w.Objects {
		fmt.Println(obj)
	}
}

// Add new Object to Worker.
func (w *Worker) Add(object Object) {
	w.Objects = append(w.Objects, object)
}

// Size returns number of Objects in Worker.
func (w *Worker) Size() int {
	return len(w.Objects)
}

// SortObjects sorts elements of Worker.Objects based on their Object.Seq value.
// See section of BySeq for definition.
func (w *Worker) SortObjects() {
	sort.Sort(BySeq(w.Objects))
}

// BySeq is a custom implementation of sort. In which we sort elements of
// Job.Objects in ascending order based on field Object.Seq.
type BySeq []Object

func (s BySeq) Len() int           { return len(s) }
func (s BySeq) Less(i, j int) bool { return s[i].Seq < s[j].Seq }
func (s BySeq) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
