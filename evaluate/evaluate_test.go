package evaluate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jthomperoo/custom-pod-autoscaler/evaluate"
	"github.com/jthomperoo/custom-pod-autoscaler/models"
	"github.com/jthomperoo/custom-pod-autoscaler/test"
)

const (
	testMetricPod                   = "test pod"
	testMetricValue                 = "test value"
	testEvaluationInvalidEvaluation = "{ \"invalid\": \"invalid\"}"
	testEvaluationInvalidJSON       = "invalid}"
	testEvaluationTargetReplicas    = int32(3)
	invalidYAML                     = "- in: -: valid - yaml"
	testEvaluate                    = "test evaluate"
	testMetric                      = "test metric"
	testInterval                    = 1234
	testHost                        = "1.2.3.4"
	testPort                        = 1234
	testMetricTimeout               = 4321
	testEvaluateTimeout             = 8765
	testNamespace                   = "test namespace"
	testScaleTargetRefKind          = "test kind"
	testScaleTargetRefName          = "test name"
	testScaleTargetRefAPIVersion    = "test api version"
	testDeploymentName              = "test deployment"
	testExecuteError                = "test error"
	testExecuteSuccess              = "test success"
)

type successExecuteValidEvaluation struct{}

func (e *successExecuteValidEvaluation) ExecuteWithPipe(command string, value string, timeout int) (*bytes.Buffer, error) {
	// Convert into JSON
	jsonEvaluation, err := json.Marshal(getTestEvaluation())
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString(string(jsonEvaluation))
	return &buffer, nil
}

type successExecuteInvalidEvaluation struct{}

func (e *successExecuteInvalidEvaluation) ExecuteWithPipe(command string, value string, timeout int) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString(testEvaluationInvalidEvaluation)
	return &buffer, nil
}

type successExecuteInvalidJSON struct{}

func (e *successExecuteInvalidJSON) ExecuteWithPipe(command string, value string, timeout int) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString(testEvaluationInvalidJSON)
	return &buffer, nil
}

func TestGetEvaluation_ExecuteFail(t *testing.T) {
	resourceMetrics := getTestResourceMetrics()
	evaluator := &evaluate.Evaluator{
		Config:   test.GetTestConfig(),
		Executer: &test.FailExecute{},
	}

	_, err := evaluator.GetEvaluation(resourceMetrics)
	if err == nil {
		t.Errorf("Expected error due to executer failing and returning an error")
		return
	}

	if err.Error() != testExecuteError {
		t.Errorf("error mismatch (-want +got):\n%s", cmp.Diff(testExecuteError, err.Error()))
	}
}

func TestGetEvaluation_ExecuteSuccessValidJSON(t *testing.T) {
	resourceMetrics := getTestResourceMetrics()
	testEvaluation := getTestEvaluation()
	evaluator := &evaluate.Evaluator{
		Config:   test.GetTestConfig(),
		Executer: &successExecuteValidEvaluation{},
	}

	evaluation, err := evaluator.GetEvaluation(resourceMetrics)
	if err != nil {
		t.Error(err)
		return
	}

	if !cmp.Equal(testEvaluation, evaluation) {
		t.Errorf("Evaluation mismatch (-want +got):\n%s", cmp.Diff(testEvaluation, evaluation))
	}
}

func TestGetEvaluation_ExecuteSuccessInvalidEvaluation(t *testing.T) {
	resourceMetrics := getTestResourceMetrics()
	evaluator := &evaluate.Evaluator{
		Config:   test.GetTestConfig(),
		Executer: &successExecuteInvalidEvaluation{},
	}
	_, err := evaluator.GetEvaluation(resourceMetrics)

	if err == nil {
		t.Errorf("Expected error due to executer returning an invalid evaluation")
		return
	}

	if _, ok := err.(*evaluate.ErrInvalidEvaluation); !ok {
		t.Errorf("Expected invalid evaluation, instead got: %v", err)
	}

	if err.Error() != fmt.Sprintf("Invalid evaluation returned by evaluator: %s", testEvaluationInvalidEvaluation) {
		t.Errorf("Error mismatch (-want +got):\n%s", cmp.Diff(testEvaluationInvalidEvaluation, err.Error()))
	}
}

func TestGetEvaluation_ExecuteSuccessInvalidJSONSyntax(t *testing.T) {
	resourceMetrics := getTestResourceMetrics()
	evaluator := &evaluate.Evaluator{
		Config:   test.GetTestConfig(),
		Executer: &successExecuteInvalidJSON{},
	}

	_, err := evaluator.GetEvaluation(resourceMetrics)

	if err == nil {
		t.Errorf("Expected error due to executer returning invalid JSON syntax")
		return
	}

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected invalid JSON syntax, instead got: %v", err)
	}
}

func getTestResourceMetrics() *models.ResourceMetrics {
	return &models.ResourceMetrics{
		DeploymentName: testDeploymentName,
		Deployment:     test.GetTestDeployment(),
		Metrics: []*models.Metric{
			&models.Metric{
				Pod:   testMetricPod,
				Value: testMetricValue,
			},
		},
	}
}

func getTestEvaluation() *models.Evaluation {
	targetReplicas := testEvaluationTargetReplicas
	return &models.Evaluation{
		TargetReplicas: &targetReplicas,
	}
}