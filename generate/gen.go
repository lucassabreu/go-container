package generate

import (
	"io"

	"github.com/lucassabreu/go-container/config"
)

// OutputTo Generates the Go code and send to the Writer
func (cf ContainerFormatter) OutputTo(w io.Writer) error {
	_, err := w.Write([]byte(""))
	return err
}

// Generate a ContainerFormatter
func Generate(cfg config.Config) (ContainerFormatter, error) {
	c := ContainerFormatter{}

	return c, nil
}
