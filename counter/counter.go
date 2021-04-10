package counter

import (
	"errors"
	"github.com/gliderlabs/logspout/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Counters map[string]Counter
}

type Counter struct {
	Match string
	Level string
}

type Adapter struct {
	Config Config
	Count  func (msg string, level string)
}

func readConfig(filename string) Config {

	var config Config

	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	err = yaml.Unmarshal(yamlFile, &config)
	check(err)
	return config
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var logCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "log_count",
		Help: "Number of times a certain log message was seen",
	},
	[]string{"msg", "level"},
)


func init() {
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(logCounter)

	router.AdapterFactories.Register(NewCounterAdapter, "counter")
}


func NewCounterAdapter(_ *router.Route) (router.LogAdapter, error) {

	var fileName = os.Getenv("LOG_COUNTER_FILE")
	if fileName == "" {
		return nil, errors.New("LOG_COUNTER_FILE must be set and contain the path to the configuration file.")
	}

	return &Adapter{
		Config: readConfig(fileName),
		Count: func(msg string, level string) { logCounter.WithLabelValues(msg, level).Inc() },
	}, nil
}


// Simple adapter implementation.
func (a *Adapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		for label, counter := range a.Config.Counters {
			if strings.Contains(message.Data, counter.Match) {
				a.Count(label, counter.Level)
			}
		}
	}
}

