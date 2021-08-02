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
