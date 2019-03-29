package fujilane

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-simplejson"

	"github.com/nerde/fuji-lane-back/fljobs"
)

type jobRoutes map[string]fljobs.JobFunc

type inlineJobsAdapter struct {
	jobs jobRoutes
}

func (a *inlineJobsAdapter) Enqueue(class string, args ...interface{}) (string, error) {
	bytes, _ := json.Marshal(args)
	js, _ := simplejson.NewJson(bytes) // Forcing numbers to become json.Number so they are properly casted inside the job
	a.jobs[class](fljobs.NewContext("test", js.MustArray()))
	return "fake-job-id", nil
}

func (a *inlineJobsAdapter) Add(class string, f fljobs.JobFunc) {
	a.jobs[class] = f
}

func (a *inlineJobsAdapter) Work() error {
	return errors.New("Inline adapter cannot work")
}
