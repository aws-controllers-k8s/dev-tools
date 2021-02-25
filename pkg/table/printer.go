package table

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
)

var (
	ErrorElementsSize           = errors.New("wrong elements size")
	ErrorUnsupportedElementType = errors.New("unsupported element type")
)

type Printer struct {
	// protects tabewriter.Writer from concurrent write operations
	mu sync.Mutex

	// table number of columns
	columns int
	tw      *tabwriter.Writer
}

// NewPrinter returns a pointer to a new Printer
func NewPrinter(columns int) *Printer {
	return &Printer{
		columns: columns,
		tw:      tabwriter.NewWriter(log.Writer(), 0, 0, 1, ' ', 0),
	}
}

// AddRaw adds a new elements to the table raws
func (p *Printer) AddRaw(elements ...interface{}) error {
	if len(elements) != p.columns {
		return ErrorElementsSize
	}

	elementsStr := []interface{}{}
	for _, element := range elements {
		newElement, err := formatElement(element)
		if err != nil {
			return err
		}
		elementsStr = append(elementsStr, newElement)
	}

	rawLine := strings.Repeat("%s\t", p.columns) + "\n"
	line := fmt.Sprintf(
		rawLine,
		elementsStr...,
	)

	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.tw.Write([]byte(line))
	return err
}

// Print flushes tabwriter
func (p *Printer) Print() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.tw.Flush()
	return err
}

func formatElement(element interface{}) (string, error) {
	var newElement interface{} = element
	switch e := element.(type) {
	case string:
	case []string:
		newElement = strings.Join(e, ",")
	case []int:
		newElement = strings.Trim(
			strings.Join(
				strings.Fields(
					fmt.Sprint(e),
				),
				","),
			"[]",
		)
	case bool:
		newElement = strconv.FormatBool(e)
	case int:
		newElement = strconv.Itoa(e)
	default:
		return "", ErrorUnsupportedElementType
	}
	return newElement.(string), nil
}
