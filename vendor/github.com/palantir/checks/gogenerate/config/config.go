// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"io/ioutil"
	"sort"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GoGenerate struct {
	// Generators is a map from the name of a generator to its configuration.
	Generators Generators `yaml:"generators" json:"generators"`
}

type Generators map[string]GeneratorConfig

func (g Generators) SortedKeys() []string {
	var sorted []string
	for k := range g {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

type GeneratorConfig struct {
	// GoGenDir is the relative path to the directory in which "go generate" should be run.
	GoGenDir string `yaml:"go-generate-dir" json:"go-generate-dir"`
	// GenPaths is the configuration that specifies the criteria for matching the output files and directories
	// generated by the "go generate" command. Any file or directory that is matched by the matchers are used to
	// determine whether or not the "go generate" command caused any changes.
	GenPaths matcher.NamesPathsCfg `yaml:"gen-paths" json:"gen-paths"`
	// Environment specifies values for the environment variables that should be set for the generator. For example, the
	// following would set GOOS to "darwin" and GOARCH to "amd64":
	//
	//   environment:
	//     GOOS: darwin
	//     GOARCH: amd64
	Environment map[string]string `yaml:"environment" json:"environment"`
}

func Load(configPath, jsonContent string) (GoGenerate, error) {
	var yml []byte
	if configPath != "" {
		var err error
		yml, err = ioutil.ReadFile(configPath)
		if err != nil {
			return GoGenerate{}, errors.Wrapf(err, "failed to read file %s", configPath)
		}
	}
	cfg, err := LoadFromStrings(string(yml), jsonContent)
	if err != nil {
		return GoGenerate{}, err
	}
	return cfg, nil
}

func LoadFromStrings(ymlContent, _ string) (GoGenerate, error) {
	cfg := GoGenerate{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return GoGenerate{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	return cfg, nil
}
