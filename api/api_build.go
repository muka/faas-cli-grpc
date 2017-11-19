//go:generate sh api.sh

package api

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/openfaas/faas-cli/builder"
	"github.com/openfaas/faas-cli/commands"
	"github.com/openfaas/faas-cli/stack"
)

func (f *faas) Build(ctx context.Context, msg *BuildRequest) (*Response, error) {

	var yamlFile string

	if len(msg.Archive) > 0 {

		target, err := ioutil.TempDir("", "faas_"+strconv.Itoa(int(time.Now().Unix())))
		if err != nil {
			return nil, err
		}

		err = Unzip(msg.Archive, target)
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(target)

		if len(msg.Yaml) > 0 {
			ymlpath := filepath.Join(target, msg.Yaml)
			if _, err := os.Stat(ymlpath); !os.IsNotExist(err) {
				yamlFile = ymlpath
			}
		}

		if len(yamlFile) == 0 {
			ymlpath := filepath.Join(target, "stack.yml")
			if _, err := os.Stat(ymlpath); !os.IsNotExist(err) {
				yamlFile = ymlpath
			}
		}

		if len(yamlFile) == 0 {
			return nil, errors.New("Cannot find YAML file in archive")
		}

		log.Printf("yml read %s", yamlFile)

	} else {
		return nil, errors.New("Missing archive field")
	}

	// consider all paths as local
	basePath := filepath.Dir(yamlFile)

	var services stack.Services
	if len(yamlFile) > 0 {
		parsedServices, err := stack.ParseYAMLFile(yamlFile, msg.Regex, msg.Filter)
		if err != nil {
			return nil, err
		}

		if parsedServices != nil {
			services = *parsedServices
			for fname, fx := range services.Functions {
				for i := range fx.EnvironmentFile {
					fx.EnvironmentFile[i] = filepath.Join(basePath, fx.EnvironmentFile[i])
				}
				fx.Handler = filepath.Join(basePath, fx.Handler)
				services.Functions[fname] = fx
			}
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
		filepath.Join(basePath, msg.Handler),
		msg.Name,
		msg.Language,
		msg.NoCache,
		msg.Squash,
		msg.Shrinkwrap,
	); err != nil {
		return nil, err
	}

	return okResponse(), nil
}
