package generate

import (
	"io"

	"github.com/lucassabreu/go-container/config"
)

// OutputTo Generates the Go code and send to the Writer
func (cf ContainerGenerator) OutputTo(w io.Writer) error {
	_, err := w.Write([]byte(""))
	return err
}

// Generate a ContainerGenerator
func Generate(cfg config.Config) (ContainerGenerator, error) {
	c := ContainerGenerator{}

	return c, nil
}
