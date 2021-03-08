package table

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatElement(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name    string
		arg     interface{}
		output  string
		wantErr bool
	}{
		{
			name:    "format string",
			arg:     "a typical string",
			output:  "a typical string",
			wantErr: false,
		},
		{
			name:    "format int",
			arg:     888,
			output:  "888",
			wantErr: false,
		},
		{
			name:    "format boolean",
			arg:     true,
			output:  "true",
			wantErr: false,
		},
		{
			name:    "format string slice",
			arg:     []string{"a", "b", "c"},
			output:  "a,b,c",
			wantErr: false,
		},
		{
			name:    "format int slice",
			arg:     []int{1, 2, 3},
			output:  "1,2,3",
			wantErr: false,
		},
		{
			name:    "unsupported type",
			arg:     time.Duration(0),
			output:  "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		o, err := formatElement(tc.arg)
		if (err != nil) != tc.wantErr {
			t.Errorf("hello() error = %v, wantErr %v", err, tc.wantErr)
			continue
		}
		assert.Equal(o, tc.output)
	}
}

func TestPrinterAddRaw(t *testing.T) {

}

func TestPrinterFlush(t *testing.T) {

}
