package fljobs

import (
	"fmt"
	"net/url"
	"strings"

	workers "github.com/jrallison/go-workers"
	"github.com/nerde/fuji-lane-back/flconfig"
)

const defaultQueue = "default"

// JobFunc is a callable job
type JobFunc func(*Context) error

// Adapter provides an interface to a background jobs implementation
type Adapter interface {
	Enqueue(string, ...interface{}) (string, error)
	Add(string, JobFunc)
	Work() error
}

// WorkersAdapter uses goworker to enqueue jobs
type WorkersAdapter struct {
	Queue     string
	StatsPort int
	jobs      map[string]JobFunc
}

// Enqueue a job to be processed in the background
func (a *WorkersAdapter) Enqueue(class string, args ...interface{}) (string, error) {
	return workers.EnqueueWithOptions(a.Queue, class, args, workers.EnqueueOptions{Retry: true})
}

// Add a new job definition
func (a *WorkersAdapter) Add(class string, f JobFunc) {
	a.jobs[class] = f
}

func (a *WorkersAdapter) handle(message *workers.Msg) {
	class := message.Json.Get("class").MustString()
	if err := a.jobs[class](NewContext(a.Queue, message.Args().MustArray())); err != nil {
		panic(err)
	}
}

// Work to start processing enqueued jobs
func (a *WorkersAdapter) Work() error {
	workers.Process(a.Queue, a.handle, 25)

	statsPort := a.StatsPort
	if statsPort == 0 {
		statsPort = 80
	}
	go workers.StatsServer(statsPort)

	workers.Run()

	return nil
}

func configureWorkers(uri string) {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	workers.Configure(getConfig(u))
}

func getConfig(u *url.URL) map[string]string {
	config := map[string]string{
		"server":   u.Host,
		"database": strings.Trim(u.Path, "/"),
		"pool":     "10",
		"process":  "1",
	}

	for k, v := range u.Query() {
		config[k] = v[0]
	}

	if pass, exists := u.User.Password(); exists {
		config["password"] = pass
	}

	return config
}

// NewWorkersAdapter creates a new GoWorkerAdapter
func NewWorkersAdapter() *WorkersAdapter {
	uri := flconfig.Config.RedisURL
	if uri == "" {
		panic(fmt.Errorf("REDIS_URL is not defined"))
	}

	configureWorkers(uri)

	return &WorkersAdapter{Queue: defaultQueue, StatsPort: flconfig.Config.JobsStatsPort, jobs: map[string]JobFunc{}}
}
