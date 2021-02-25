package testing

import (
	"sync"

	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

type TestStatus int

const (
	TestStatusInProgress TestStatus = iota
	TestStatusSuccess
	TestStatusError
	TestStatusFail
)

func (ts TestStatus) String() string {
	switch ts {
	case TestStatusError:
		return "error"
	case TestStatusFail:
		return "fail"
	case TestStatusSuccess:
		return "success"
	case TestStatusInProgress:
		return "in progress"
	default:
		panic("unsupported test status")
	}
}

type Tester interface {
	Run([]string, <-chan int) error
}

type Test struct {
	mu sync.Mutex

	cmd    []string
	repo   *repository.Repository
	result TestResult
}

func (Test) Run(stop <-chan struct{}) {

}

func (t Test) Wait() {

}

type TestResult struct {
	Success bool
	Output  []byte
}
