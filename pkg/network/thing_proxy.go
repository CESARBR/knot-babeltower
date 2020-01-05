package network

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	Create(id, name, authorization string) (idGenerated string, err error)
}

type proxy struct {
	url    string
	logger logging.Logger
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) ThingProxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return proxy{url, logger}
}

// Create proxy the http request to thing service
func (p proxy) Create(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	return idGenerated, nil
}
