package container

import "time"

// Container is a basic definition for a container
type Container interface {
	Get(name string) interface{}
	Set(name string, service interface{}) Container

	SetParameters(map[string]interface{}) Container

	GetParameter(name string) interface{}

	GetBoolP(name string) bool
	GetBoolSlideP(name string) bool
	GetDurationP(name string) time.Duration
	GetDurationSlideP(name string) time.Duration
	GetFloat32P(name string) float32
	GetFloat32SlideP(name string) float32
	GetFloat64P(name string) float64
	GetFloat64SlideP(name string) float64
	GetIntP(name string) int
	GetIntSlideP(name string) int
	GetInt8P(name string) int8
	GetInt8SlideP(name string) int8
	GetInt16P(name string) int16
	GetInt16SlideP(name string) int16
	GetInt32P(name string) int32
	GetInt32SlideP(name string) int32
}
