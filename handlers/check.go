package handlers

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// CheckGroup represents a struct that checks on whether our service is ready to receive requests or not
type CheckGroup struct {
	build string
	log   *zerolog.Logger
}

// NewCheckGroup creates a new instance of CheckGroup
func NewCheckGroup(log *zerolog.Logger, build string) CheckGroup {
	return CheckGroup{
		build: build,
		log:   log,
	}
}

// Liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (cg CheckGroup) Liveness(w http.ResponseWriter, r *http.Request) {

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     cg.build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	Respond(w, info, http.StatusOK, true, nil)
}
