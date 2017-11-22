//go:generate sh api.sh

package api

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/openfaas/faas-cli/commands"
	"github.com/openfaas/faas-cli/stack"
)

func (f *faas) Build(ctx context.Context, msg *BuildRequest) (*Response, error) {

	var yamlFile string
	var parsedServices *stack.Services

	if len(msg.Archive) > 0 {

		name := msg.Name
		if len(msg.Name) == 0 {
			name = "faas_" + strconv.Itoa(int(time.Now().Unix()))
		}

		target := filepath.Join(
			os.Getenv("workdir"),
			"./store/",
			name,
		)

		err := Unzip(msg.Archive, target)
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

		parsedServices, err = stack.ParseYAMLFile(yamlFile, msg.Regex, msg.Filter)
		if err != nil {
			return nil, err
		}

		// consider all paths as local
		basePath := filepath.Dir(yamlFile)
		for fname, fx := range parsedServices.Functions {
			for i := range fx.EnvironmentFile {
				fx.EnvironmentFile[i] = filepath.Join(basePath, fx.EnvironmentFile[i])
			}
			fx.Handler = filepath.Join(basePath, fx.Handler)
			parsedServices.Functions[fname] = fx
		}

		log.Printf("yml read %s", yamlFile)

	} else {
		return nil, errors.New("Missing archive field")
	}

	arg := commands.BuildArguments{
		YamlFile:     yamlFile,
		Services:     parsedServices,
		Filter:       msg.Filter,
		Regex:        msg.Regex,
		Image:        msg.Image,
		Handler:      msg.Handler,
		FunctionName: msg.Name,
		Language:     msg.Language,
		Nocache:      msg.NoCache,
		Squash:       msg.Squash,
		Parallel:     int(msg.Parallel),
		Shrinkwrap:   msg.Shrinkwrap,
	}

	err := commands.Build(arg)
	if err != nil {
		return nil, err
	}
	return okResponse(), nil
}
