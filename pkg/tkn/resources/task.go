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
