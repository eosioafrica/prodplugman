package ppman

import (
	"net/http"
	"context"
	"github.com/sirupsen/logrus"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

var (

	pause  =  "v1/producer/pause"  // pause production
	resume =  "v1/producer/resume" // resume from paused production
	state  =  "v1/producer/paused" // return true/false depending on whether the node is currently paused
)

type Status struct {

	State 	string
	Err		error
}

type NodeosEndpoint struct {

	URL 			string
	Port 			int
	Resource		string

	Err  			error
}

func (pp *NodeosEndpoint) PpHandler() {


}

func (pp *NodeosEndpoint) PauseProduction(ctx *context.Context) {


	url := fmt.Sprintf("%s:%s/%s", pp.URL, pp.Port, pause)

	result, _, err := handleRequest(ctx, "POST", url)
	if err != nil { pp.Err = err }
	if result >= 200 && result <= 299 {

		return
	} else {

		pp.Err= errors.New("Error attempting to pause production.")
		logrus.Error("Error attempting to pause production.")
	}
}

func (pp *NodeosEndpoint) ResumeProduction(ctx *context.Context)  {

	url := fmt.Sprintf("%s:%s/%s", pp.URL, pp.Port, resume)

	result, _, err := handleRequest(ctx, "POST", url)
	if err != nil { pp.Err = err }
	if result >= 200 && result <= 299 {

		return
	} else {

		pp.Err= errors.New("Error attempting to resume production.")
		logrus.Error("Error attempting to resume production.")
	}
}

func (pp *NodeosEndpoint) ProductionState(ctx *context.Context) (Status) {

	url := fmt.Sprintf("%s:%s/%s", pp.URL, pp.Port, state)

	result, body,  err := handleRequest(ctx, "POST", url)
	if err != nil { pp.Err = err }
	if result >= 200 && result <= 299 {

		htmlData, err := ioutil.ReadAll(body)

		if err != nil {
			fmt.Println(err)
			pp.Err = err
			return Status{ "unknown", err }
		}

		switch string(htmlData) {

			case "true", "TRUE", "True":
				return Status{ "on", err }
			case "false", "FALSE", "False":
				return Status{ "off", err }
		}

		pp.Err = errors.New("Syntax err parsing bool at ppman.ProductionState")
		return Status{ "unknown", pp.Err }

	} else {

		pp.Err= errors.New("Error attempting to retriece current production state.")
		logrus.Error("Error attempting to retriece current production state.")
		return Status{ "unknown", pp.Err }
	}
}

func (pp *NodeosEndpoint) CheckIfApiIsReachable(ctx *context.Context) (bool) {

	url := fmt.Sprintf("%s:%s/%s", pp.URL, pp.Port, pp.Resource)

	result, _, err := handleRequest(ctx, "GET", url)
	if err != nil { pp.Err = err }
	if result >= 200 && result <= 299 {

		return true
	} else {

		pp.Err= errors.New("Failed to quiry nodeos through API.")
		logrus.Error("Failed to quiry nodeos through API.")

		return false
	}
}

func handleRequest(ctx *context.Context, method string, url string) (int, io.Reader, error) {

	req, _ := http.NewRequest(method, url, nil)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)

	if err != nil { return -1, nil, err }

	defer resp.Body.Close()

	return resp.StatusCode, resp.Body, nil
}
