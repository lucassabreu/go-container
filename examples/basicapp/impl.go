package basicapp

import (
	"fmt"
	"strings"
)

// Reporter represents a report like implementation
type Reporter interface {
	Name() string
	Report() string
}

// NewSimpleReporter creates a report that just shows its name
func NewSimpleReporter(name string) Reporter {
	return simple{name}
}

type simple struct {
	name string
}

func (s simple) Name() string {
	return s.name
}

func (s simple) Report() string {
	return fmt.Sprintf("Simple Report named %s\n", s.name)
}

// NewReportThatFails just fails at creating something
func NewReportThatFails() (Reporter, error) {
	return nil, fmt.Errorf("Don't like you")
}

// NewReportThatNotFails that do not fail
func NewReportThatNotFails() (Reporter, error) {
	return NewSimpleReporter("same, same, but different"), nil
}

// NewComposedReport creates a Reporter composed of others
func NewComposedReport(rs ...Reporter) Reporter {
	return composed{reports: rs}
}

type composed struct {
	reports []Reporter
}

func (c composed) Name() string {
	b := strings.Builder{}
	for _, r := range c.reports {
		b.WriteString(r.Name())
		b.WriteRune(',')
	}
	return b.String()
}

func (c composed) Report() string {
	b := strings.Builder{}
	for _, r := range c.reports {
		b.WriteString(r.Report())
		b.WriteRune('\n')
	}
	return b.String()
}

// SimplerReport is a struct without "builder"
type SimplerReport struct {
	ReportName string
}

// Name is the name of it
func (s SimplerReport) Name() string {
	return s.ReportName
}

// Report reports something
func (s SimplerReport) Report() string {
	return "Losing creativity\n"
}

// ComposedSimplerReport will be used to tests Go composition
type ComposedSimplerReport struct {
	SimplerReport
	Extra string
}

// Report reports something (composed)
func (s ComposedSimplerReport) Report() string {
	return "Composed of " + s.SimplerReport.Report()
}

// NewComposedSimplerReport hello
func NewComposedSimplerReport(name string) ComposedSimplerReport {
	return ComposedSimplerReport{
		SimplerReport: SimplerReport{
			ReportName: name,
		},
		Extra: "extra " + name,
	}
}

// ProxyReporter proxy calls to other
type ProxyReporter struct {
	CSR *ComposedSimplerReport
}

// Name shows name
func (pr ProxyReporter) Name() string {
	return pr.CSR.Name()
}

// Report shows report
func (pr ProxyReporter) Report() string {
	return pr.CSR.Report()
}
