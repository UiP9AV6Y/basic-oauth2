package web

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidData = errors.New("malformed health information received")
)

type Consumer interface {
	UnmarshalSample() interface{}
	Consume(interface{}) error
}

// API response structure as defined by the IETF Draft
// for Health Check Response Format for HTTP APIs
// (https://inadarei.github.io/rfc-healthcheck/)
type HealthResponse struct {
	Status      string                   `json:"status,omitempty"`
	Output      string                   `json:"output,omitempty"`
	Version     string                   `json:"version,omitempty"`
	ReleaseId   string                   `json:"releaseId,omitempty"`
	ServiceId   string                   `json:"serviceId,omitempty"`
	Description string                   `json:"description,omitempty"`
	Notes       []string                 `json:"notes,omitempty"`
	Checks      map[string][]interface{} `json:"checks,omitempty"`
	Links       map[string]string        `json:"links,omitempty"`
}

// HealthConsumer is a response consumer implementation expecting
// API health information (HealthResponse) as input.
type HealthConsumer []string

// NewAPIHealthConsumer creates a health consumer instance
// based on the IETF Draft for Health Check Response Format for HTTP APIs.
// (https://inadarei.github.io/rfc-healthcheck/)
func NewAPIHealthConsumer() HealthConsumer {
	return NewHealthConsumer("pass", "ok", "up")
}

func NewHealthConsumer(okStatus ...string) HealthConsumer {
	return HealthConsumer(okStatus)
}

func (c HealthConsumer) UnmarshalSample() interface{} {
	return new(HealthResponse)
}

func (c HealthConsumer) Consume(data interface{}) error {
	health, ok := data.(HealthResponse)
	if !ok {
		return ErrInvalidData
	}

	for _, healthy := range c {
		if health.Status == healthy {
			return nil
		}
	}

	return fmt.Errorf("health status is %q", health.Status)
}
