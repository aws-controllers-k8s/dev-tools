package testing

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Result struct {
	Type     string
	Service  string
	Logs     [][]byte
	Error    error
	ExitCode int
	Exited   bool
}

func (r *Result) SaveLogs(directory string) error {
	current := time.Now()
	filename := fmt.Sprintf("%s-%s-%s.logs", r.Service, r.Type, current.Format("2021-01-01 00:00:00"))
	fullPath := filepath.Join(directory, filename)

	logs := []byte{}
	for _, logLine := range r.Logs {
		logs = append(logs, logLine...)
	}
	return ioutil.WriteFile(fullPath, logs, 0644)
}

func (r *Result) TableLine() string {
	return fmt.Sprintf("%s\t%v\t\n", r.Service, r.Status())
}

func (r *Result) Status() TestStatus {
	if !r.Exited {
		return TestStatusInProgress
	} else if r.Error != nil {
		return TestStatusError
	} else if r.ExitCode != 0 {
		return TestStatusFail
	} else {
		return TestStatusSuccess
	}
}
