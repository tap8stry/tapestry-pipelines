//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package tkn

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	"github.ibm.scs.com/tapestry/pkg/utils"
)

const (
	task           string = "Task"
	taskrun        string = "TaskRun"
	pipeline       string = "Pipeline"
	triggerbinding string = "TriggerBinding"
)

type parsedTknObjects struct {
	ManifestFilepath string
	ManifestFilehash string
	TknObjects       []tknObject
}

type tknObject struct {
	GroupKind     string
	RuntimeObject []byte
}

//getTknResources :
func getTknResources(file string) []parsedTknObjects {
	parsedObjs := []parsedTknObjects{}
	if utils.IsYAMLFile(file) {
		if filebuf, err := ioutil.ReadFile(file); err == nil {
			p := parsedTknObjects{}
			p.ManifestFilepath = file
			p.ManifestFilehash = fmt.Sprintf("%x", md5.Sum(filebuf))
			p.TknObjects = parseK8sYaml(filebuf)
			parsedObjs = append(parsedObjs, p)
		}
	}
	return parsedObjs
}

func parseK8sYaml(fileR []byte) []tknObject {
	dObjs := []tknObject{}
	acceptedK8sTypes := regexp.MustCompile(fmt.Sprintf("(%s|%s|%s|%s)",
		task, pipeline, taskrun, triggerbinding))
	fileAsString := string(fileR[:])
	sepYamlfiles := strings.Split(fileAsString, "---")
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		_, groupVersionKind, err := decode([]byte(f), nil, nil)
		if err != nil {
			continue
		}
		if acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			d := tknObject{}
			d.GroupKind = groupVersionKind.Kind
			d.RuntimeObject = []byte(f)
			dObjs = append(dObjs, d)
		}
	}
	return dObjs
}
