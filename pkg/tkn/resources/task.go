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
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	"github.ibm.scs.com/tapestry/pkg/sign"
	"go.uber.org/zap"
)

// ParseTknTask parses Task
func ParseTknTask(r []byte) *v1beta1.Task {
	if r == nil {
		return nil
	}
	var task v1beta1.Task
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(r, nil, &task)
	if err != nil {
		zap.S().Debugf("error parsing `task' object: %v", err)
	}
	return &task
}

//GetTaskSteps :
func GetTaskSteps(task *v1beta1.Task) []sign.StepSC {
	var stepList []sign.StepSC
	for _, s := range task.TaskSpec().Steps {
		var step sign.StepSC
		step.Name = s.Name
		step.ImageRef = s.Image
		stepList = append(stepList, step)
	}
	return stepList
}
