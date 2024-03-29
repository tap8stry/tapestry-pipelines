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
	"context"
	"strings"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/tapestry-pipelines/pkg/sign"
	"github.com/tapestry-pipelines/pkg/tkn/resources"
	"github.com/tapestry-pipelines/pkg/utils"
)

type tknResources struct {
	RawPipelines []tknPipeline
	RawTasks     []tknTask
}

type tknTask struct {
	Task     *v1beta1.Task
	Filepath string
}

type tknPipeline struct {
	Pipeline *v1beta1.Pipeline
	Filepath string
}

//GenSignCandidates :
func GenSignCandidates(ctx context.Context, filepath, pipeline string) (*sign.TknSignCandidates, error) {
	var signCandidates sign.TknSignCandidates

	files, _ := utils.FilePathWalkDir(filepath)

	var rawResources tknResources

	for _, f := range files {
		pObj := getTknResources(f)
		if len(pObj) == 0 {
			// zap.S().Debugf("no tkn pipeline definitions found")
			continue
		}
		for _, p := range pObj {
			for _, o := range p.TknObjects {
				// zap.S().Debugf("file %s parsed resource object: %s", f, o.GroupKind)
				scanTknObject(f, o.GroupKind, o.RuntimeObject, &rawResources)
			}
		}
	}
	parseSignCandidates(&rawResources, &signCandidates, &pipeline)
	return &signCandidates, nil
}

func scanTknObject(filepath, kind string, objDataBuf []byte, raw *tknResources) {
	switch kind {
	case "Task":
		t := resources.ParseTknTask(objDataBuf)
		// zap.S().Debugf("task name %v", t.Name)
		// if _, exists := taskTmp[t.Name]; exists {
		// 	zap.S().Debugf("task `%s` defined in multiple files:", t.Name)
		// }
		rawT := tknTask{}
		rawT.Task = t
		rawT.Filepath = filepath
		raw.RawTasks = append(raw.RawTasks, rawT)
	case "Pipeline":
		// p := resources.PipelineCntr{}
		pipeline, _ := resources.ParseTknPipeline(objDataBuf)
		// taskNames, _ := resources.GetTaskNames(pipeline)
		// zap.S().Debugf("task refs %v ", taskNames)
		rawP := tknPipeline{}
		rawP.Filepath = filepath
		rawP.Pipeline = pipeline
		raw.RawPipelines = append(raw.RawPipelines, rawP)
	}
}

func parseSignCandidates(rawResources *tknResources, signCandidates *sign.TknSignCandidates, optPipeline *string) {
	for _, p := range rawResources.RawPipelines {
		if *optPipeline != "" && !strings.EqualFold(*optPipeline, p.Pipeline.Name) {
			continue
		}
		var pipelineCandidate sign.PipelineSC
		pipelineCandidate.Filepath = p.Filepath
		pipelineCandidate.Name = p.Pipeline.Name
		pipelineCandidate.PipelineObj = p.Pipeline
		taskNames, _ := resources.GetTaskNames(p.Pipeline)
		for _, taskName := range taskNames {
			for _, t := range rawResources.RawTasks {
				if strings.EqualFold(taskName, t.Task.Name) {
					var taskCandidate sign.TaskSC
					taskCandidate.Filepath = t.Filepath
					taskCandidate.Name = t.Task.Name
					taskCandidate.TaskObj = t.Task
					taskCandidate.Steps = resources.GetTaskSteps((t.Task))
					pipelineCandidate.TaskRefs = append(pipelineCandidate.TaskRefs, taskCandidate)
				}
			}
		}
		signCandidates.PipelinesSC = append(signCandidates.PipelinesSC, pipelineCandidate)
	}
}
