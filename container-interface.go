package container

import "time"

// ParametersBag is a bag that will provide parameters values to the container
type ParametersBag interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetDurationSlide(name string) []time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
}

// Container is a basic definition for a container
type Container interface {
	Get(name string) interface{}
	Set(name string, service interface{}) Container

	SetParametersBag(bag ParametersBag) Container
	GetParametersBag() ParametersBag
}
