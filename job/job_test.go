package job_test

import (
	"reflect"
	"testing"

	"github.com/rolandvarga/batch_proc/job"
)

func TestGenerateGroupID(t *testing.T) {
	var cases = []struct {
		input int64
		want  int64
	}{
		{0, 0},
		{1, 0},
		{99, 0},
		{100, 1},
		{110, 1},
		{200, 2},
		{201, 2},

		{301, 3},
		{309, 3},
		{310, 3},
		{350, 3},
		{399, 3},

		{1051, 10},
		{1151, 11},

		{19365, 193},
		{19465, 194},

		{49999, 499},
		{50000, 500},
	}

	for _, c := range cases {
		obj := job.Object{Seq: c.input}

		got := obj.GenerateGroupID()
		if got != c.want {
			t.Errorf("got '%d' want '%d'", got, c.want)
		}
	}
}

func TestAddObject(t *testing.T) {
	var cases = []struct {
		objects []job.Object
		groupID int64
	}{
		{
			[]job.Object{
				job.Object{ID: "object_1", Seq: 1, Data: ""},
				job.Object{ID: "object_3", Seq: 3, Data: ""},
				job.Object{ID: "object_2", Seq: 2, Data: ""},
			}, 1,
		},
	}

	for _, c := range cases {
		w := job.NewWorker(c.groupID)

		for _, obj := range c.objects {
			w.Add(obj)
		}

		if !reflect.DeepEqual(c.objects, w.Objects) {
			t.Errorf("got '%v' want '%v'", w.Objects, c.objects)
		}
	}
}

func TestSize(t *testing.T) {
	var cases = []struct {
		objects []job.Object
		groupID int64
		size    int
	}{
		{
			[]job.Object{
				job.Object{ID: "object_1", Seq: 1, Data: ""},
				job.Object{ID: "object_3", Seq: 3, Data: ""},
				job.Object{ID: "object_2", Seq: 2, Data: ""},
			}, 0, 3,
		},
		{
			[]job.Object{}, 0, 0,
		},
	}

	for _, c := range cases {
		w := job.NewWorker(c.groupID)

		for _, obj := range c.objects {
			w.Add(obj)
		}

		got := w.Size()
		if got != c.size {
			t.Errorf("got '%d' want '%d'", got, c.size)
		}
	}
}
func TestSort(t *testing.T) {
	var cases = []struct {
		objects []job.Object
		groupID int64
		want    []job.Object
	}{
		{
			[]job.Object{
				job.Object{ID: "object_199", Seq: 199, Data: ""},
				job.Object{ID: "object_103", Seq: 103, Data: ""},
				job.Object{ID: "object_102", Seq: 102, Data: ""},
				job.Object{ID: "object_154", Seq: 154, Data: ""},
			},
			1,
			[]job.Object{
				job.Object{ID: "object_102", Seq: 102, Data: ""},
				job.Object{ID: "object_103", Seq: 103, Data: ""},
				job.Object{ID: "object_154", Seq: 154, Data: ""},
				job.Object{ID: "object_199", Seq: 199, Data: ""},
			},
		},
	}

	for _, c := range cases {
		w := job.NewWorker(c.groupID)

		for _, obj := range c.objects {
			w.Add(obj)
		}

		w.SortObjects()
		got := w.Objects
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("got '%v' want '%v'", got, c.want)
		}
	}
}

func FindWorkerWithID(t *testing.T) {
	var cases = []struct {
		job     job.Job
		groupID int64
		want    int
	}{
		{job.Job{
			Workers: []job.Worker{
				job.Worker{GroupID: 1},
				job.Worker{GroupID: 2},
				job.Worker{GroupID: 3},
			},
		}, 4, -1},

		{job.Job{
			Workers: []job.Worker{
				job.Worker{GroupID: 1},
				job.Worker{GroupID: 2},
				job.Worker{GroupID: 3},
			},
		}, 2, 1},
	}

	for _, c := range cases {
		got := c.job.FindWorkerWithID(c.groupID)
		if got != c.want {
			t.Errorf("got '%d' want '%d' for groupID %d", got, c.want, c.groupID)
		}
	}
}

func TestAddWorker(t *testing.T) {
	var cases = []struct {
		workers []job.Worker
	}{
		{
			[]job.Worker{
				job.Worker{GroupID: 1},
				job.Worker{GroupID: 2},
				job.Worker{GroupID: 3},
			},
		},
	}

	for _, c := range cases {
		j := job.NewJob(job.Data{})
		for _, w := range c.workers {
			j.Add(w)
		}
		if !reflect.DeepEqual(c.workers, j.Workers) {
			t.Errorf("got '%v' want '%v'", j.Workers, c.workers)
		}
	}
}
