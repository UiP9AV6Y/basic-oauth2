package web

import (
	"encoding/json"
	"net/http"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/version"
)

const (
	HealthEndpoint = "/health"
)

type HealthRouter struct {
	logger *log.Controller
}

func NewHealthRouter(logger *log.Controller) *HealthRouter {
	router := &HealthRouter{
		logger: logger,
	}

	return router
}

func (r *HealthRouter) Handler() http.Handler {
	mux := http.NewServeMux()

	r.DecorateHandler(mux)

	return mux
}

func (r *HealthRouter) DecorateHandler(mux *http.ServeMux) {
	mux.HandleFunc(HealthEndpoint, r.HandleHealth)
}

// HandleHealth returns the system health information.
func (r *HealthRouter) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/health+json")
	w.WriteHeader(http.StatusOK)

	body := map[string]interface{}{
		"version":   "1",
		"releaseId": version.Version(),
	}
	body["checks"], err = r.collectChecksStatus()
	if err == nil {
		body["status"] = "pass"
	} else {
		body["status"] = "fail"
		body["output"] = err.Error()
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(body)
	if err != nil {
		r.logger.Error().Printf("%s: %v", HealthEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *HealthRouter) collectChecksStatus() (map[string]interface{}, error) {
	checks := map[string]interface{}{}

	return checks, nil
}
