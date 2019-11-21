package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
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

type objKnot struct {
	ID string `json:"id"`
}

type objMetadata struct {
	Knot objKnot `json:"knot"`
}

func (p proxy) getJSONBody(id, name string) ([]byte, error) {
	body := struct {
		Name     string      `json:"name"`
		Metadata objMetadata `json:"metadata"`
	}{
		Name: name,
		Metadata: objMetadata{
			Knot: objKnot{
				ID: id,
			},
		},
	}
	return json.Marshal(body)
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) ThingProxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return proxy{url, logger}
}

// Create registers a thing on service and return the id generated
func (p proxy) Create(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	/**
	 * Add Timeout in http.Client to avoid blocking the request.
	 */
	client := &http.Client{Timeout: 10 * time.Second}
	jsonBody, err := p.getJSONBody(id, name)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	req, err := http.NewRequest("POST", p.url+"/things", bytes.NewBuffer(jsonBody))
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		p.logger.Errorf("Status not created: %d", resp.StatusCode)
		switch resp.StatusCode {
		case http.StatusConflict:
			// TODO: Verify uniqueness from KNoT ID
			err = entities.ErrEntityExists{Msg: "Thing exists"}
		case http.StatusForbidden:
			err = entities.ErrNoPerm{Msg: "The authorization token has no permission to create a thing"}
		}
		return "", err
	}

	locationHeader := resp.Header.Get("Location")
	idGenerated = locationHeader[len("/things/"):] // get substring after "/things/"

	return idGenerated, err
}
