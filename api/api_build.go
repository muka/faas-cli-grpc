//go:generate sh api.sh

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/openfaas/faas-cli/builder"
	"github.com/openfaas/faas-cli/commands"
	"github.com/openfaas/faas-cli/stack"
)

func (f *faas) Build(ctx context.Context, msg *BuildRequest) (*Response, error) {

	var yamlFile string
	var regex string
	var filter string

	if len(msg.Archive) > 0 {

		target, err := ioutil.TempDir("", "faas_"+strconv.Itoa(int(time.Now().Unix())))
		if err != nil {
			return nil, err
		}

		Unzip(msg.Archive, target)

		// defer os.RemoveAll(target)
	}

	return nil, nil

	var services stack.Services
	if len(yamlFile) > 0 {
		parsedServices, err := stack.ParseYAMLFile(yamlFile, regex, filter)
		if err != nil {
			return nil, err
		}

		if parsedServices != nil {
			services = *parsedServices
		}
	}

	if pullErr := commands.PullTemplates(""); pullErr != nil {
		return nil, fmt.Errorf("could not pull templates for OpenFaaS: %v", pullErr)
	}

	if len(services.Functions) > 0 {
		err := commands.BuildStack(&services, int(msg.Parallel), msg.Shrinkwrap)
		if err != nil {
			return nil, err
		}
		return okResponse(), nil
	}

	if len(msg.Image) == 0 {
		return nil, fmt.Errorf("please provide a valid --image name for your Docker image")
	}
	if len(msg.Handler) == 0 {
		return nil, fmt.Errorf("please provide the full path to your function's handler")
	}
	if len(msg.Name) == 0 {
		return nil, fmt.Errorf("please provide the deployed --name of your function")
	}
	if err := builder.BuildImage(
		msg.Image,
		msg.Handler,
		msg.Name,
		msg.Lang,
		msg.NoCache,
		msg.Squash,
		msg.Shrinkwrap,
	); err != nil {
		return nil, err
	}

	return okResponse(), nil
}
