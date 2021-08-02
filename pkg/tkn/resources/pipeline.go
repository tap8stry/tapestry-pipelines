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

package resources

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	"go.uber.org/zap"
)

//PipelineCntr :
type PipelineCntr struct {
	Pipeline v1beta1.Pipeline
}

// ParseTknPipeline parses Task
func ParseTknPipeline(r []byte) (*v1beta1.Pipeline, error) {
	var pipeline v1beta1.Pipeline
	if r == nil {
		return nil, nil
	}

	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(r, nil, &pipeline)
	if err != nil {
		zap.S().Debugf("error parsing `pipeline' object: %v", err)
		return nil, err
	}

	return &pipeline, nil
}

//GetTaskNames :
func GetTaskNames(pipeline *v1beta1.Pipeline) ([]string, error) {
	taskNameList := []string{}
	tasks := pipeline.Spec.Tasks
	for _, t := range tasks {
		taskNameList = append(taskNameList, t.TaskRef.Name)
	}
	return taskNameList, nil
}
