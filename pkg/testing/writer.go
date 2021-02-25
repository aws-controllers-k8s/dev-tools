package testing

import (
	"io"
	"sync"
)

func NewMultiWriter()

type MultiWriter struct {
	mu sync.Mutex

	writers []io.Writer
}
