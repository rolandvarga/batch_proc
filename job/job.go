package job

import (
	"fmt"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

var (
	batchSize       = 100
	numberOfWorkers = 500
	numberOfJobs    = 50000
	topic           = "batch_proc"
)

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
	return Job{Data: data, Workers: make([]Worker, 500)}
}

// Object represents a single element in our input data.
type Object struct {
	ID   string `json:"id"`
	Seq  int64  `json:"seq"`
	Data string `json:"data"`
}

// GenerateGroupID returns the groupID that each object is grouped under
// when assigned to a worker.
func (o *Object) GenerateGroupID() int64 {
	return o.Seq / int64(batchSize)
}

// GetIndex returns the modulo of each sequence based on the current batch size.
// This makes up the index used to store the object in Worker.Objects
func (o *Object) GetIndex() int64 {
	return o.Seq % int64(batchSize)
}

// Message returns the concatenated string representation of fields id, seq & data.
func (o *Object) Message() string {
	return fmt.Sprintf("%s %d %s", o.ID, o.Seq, o.Data)
}

// Process will take an object and produce it as an entry for kafka topics.
func (o *Object) Process() {
	// partition := 0
	// fmt.Println(o.Message())
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	conn, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println(err)
	}

	msg := o.Message()

	conn.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})
}

// Worker is responsible for holding and processing a single batch of objects.
type Worker struct {
	GroupID    int64
	Objects    []Object
	Queue      chan Object
	Done       chan bool
	Killsignal chan bool
}

// NewWorker initializes a new Worker instance.
func NewWorker(queue chan Object, groupID int64, done, ks chan bool) *Worker {
	return &Worker{
		Queue:      queue,
		GroupID:    groupID,
		Done:       done,
		Killsignal: ks,
		Objects:    make([]Object, batchSize),
	}
}

// Run allows the worker to listen on its Queue channel and process a batch when ready.
func (w *Worker) Run() {
	count := 0
	for true {
		if count == 100 {
			for _, obj := range w.Objects {
				obj.Process()
				w.Done <- true
			}
		}
		select {
		case o := <-w.Queue:
			// received new object from queue
			idx := o.GetIndex()
			w.Objects[idx] = o
			count++
		case <-w.Killsignal:
			fmt.Printf("stopping worker with ID %d\n", w.GroupID)
			return
		}
	}
}

// Run starts processing the dataset found in Job
func (j *Job) Run() {
	//channel for terminating the workers
	killsig := make(chan bool)
	done := make(chan bool)

	for i := 0; i < numberOfWorkers; i++ {
		q := make(chan Object)

		worker := NewWorker(q, int64(i), done, killsig)
		j.Workers[i] = *worker
		go worker.Run()
	}

	for _, obj := range j.Data.Objects {
		go func(obj Object) {
			gid := obj.GenerateGroupID()
			j.Workers[gid].Queue <- obj
		}(obj)
	}

	// a deadlock occurs if c >= numberOfJobs
	for c := 0; c < numberOfJobs; c++ {
		<-done
	}

	fmt.Println("done!")

	close(killsig)
	time.Sleep(2 * time.Second)
}
