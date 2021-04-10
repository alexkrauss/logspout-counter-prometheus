package counter

import (
	"gotest.tools/assert"
	"testing"
)


var exampleConfig = Config{Counters: map[string]Counter{
	"hello": {"Hello", "INFO"},
	"world": {"World", "ERROR"},
}}

func TestFileParsing(t *testing.T) {

	cfg := readConfig("example_counters.yaml")
	assert.DeepEqual(t, cfg, exampleConfig)
}

//func TestCounterAdapter_Stream(t *testing.T) {
//
//	var result = []string{}
//
//	var a = Adapter{
//		Config: exampleConfig,
//		Count: func(msg string, level string) {
//			result = append(result, msg + "_" + level)
//		},
//	}
//}