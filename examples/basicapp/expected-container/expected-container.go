package container

import (
	"fmt"
	"strings"

	cBase "github.com/lucassabreu/go-container"
	"github.com/lucassabreu/go-container/examples/basicapp"
)

// ReportContainer is a container for the application
type ReportContainer struct {
	parametersBag cBase.ParametersBag

	composedReport                         *basicapp.Reporter
	simpleReportA                          *basicapp.Reporter
	simpleReportB                          *basicapp.Reporter
	reportThatFails                        *basicapp.Reporter
	reportThatNotFails                     *basicapp.Reporter
	simplerReport                          *basicapp.SimplerReport
	composedSimplerReport                  *basicapp.ComposedSimplerReport
	composedSimplerReportWithSimplerReport *basicapp.ComposedSimplerReport
	proxyReporter                          *basicapp.ProxyReporter
}

func (c *ReportContainer) onInitError(serviceName string, err error) {
	panic(fmt.Errorf("Error on initializing service \"%s\": %v", serviceName, err))
}

// SetParametersBag sets the ParametersBag for the container
func (c *ReportContainer) SetParametersBag(bag cBase.ParametersBag) *ReportContainer {
	c.parametersBag = bag
	return c
}

// GetParametersBag get the container's ParametersBag
func (c *ReportContainer) GetParametersBag() cBase.ParametersBag {
	return c.parametersBag
}

// GetComposedReport returns the ComposedReport service
func (c *ReportContainer) GetComposedReport() *basicapp.Reporter {
	if c.composedReport == nil {
		t := basicapp.NewComposedReport(
			*c.GetSimpleReportA(),
			*c.GetSimpleReportB(),
			*c.GetSimplerReport(),
		)
		c.composedReport = &t
	}
	return c.composedReport
}

// GetSimpleReportA returns the SimpleReportA service
func (c *ReportContainer) GetSimpleReportA() *basicapp.Reporter {
	if c.simpleReportA == nil {
		t := basicapp.NewSimpleReporter(
			"A",
		)
		c.simpleReportA = &t
	}
	return c.simpleReportA
}

// GetSimpleReportB returns the SimpleReportB service
func (c *ReportContainer) GetSimpleReportB() *basicapp.Reporter {
	if c.simpleReportB == nil {
		t := basicapp.NewSimpleReporter(
			"B",
		)
		c.simpleReportB = &t
	}
	return c.simpleReportB
}

// GetReportThatFails returns the ReportThatFails service
func (c *ReportContainer) GetReportThatFails() *basicapp.Reporter {
	if c.reportThatFails == nil {
		t, err := basicapp.NewReportThatFails()
		c.onInitError("ReportThatFails", err)
		c.reportThatFails = &t
	}
	return c.reportThatFails
}

// GetReportThatNotFails returns the ReportThatNotFails service
func (c *ReportContainer) GetReportThatNotFails() *basicapp.Reporter {
	if c.reportThatNotFails == nil {
		t, err := basicapp.NewReportThatNotFails()
		c.onInitError("ReportThatNotFails", err)
		c.reportThatNotFails = &t
	}
	return c.reportThatNotFails
}

// GetSimplerReport returns the SimplerReport service
func (c *ReportContainer) GetSimplerReport() *basicapp.SimplerReport {
	if c.simplerReport == nil {
		c.simplerReport = &basicapp.SimplerReport{
			ReportName: "something simple",
		}
	}
	return c.simplerReport
}

// GetComposedSimplerReport returns the ComposedSimplerReport service
func (c *ReportContainer) GetComposedSimplerReport() *basicapp.ComposedSimplerReport {
	if c.composedSimplerReport == nil {
		c.composedSimplerReport = &basicapp.ComposedSimplerReport{
			SimplerReport: basicapp.SimplerReport{
				ReportName: "something simple",
			},
			Extra: "spice",
		}
	}
	return c.composedSimplerReport
}

// GetComposedSimplerReportWithSimplerReport returns the ComposedSimplerReportWithSimplerReport service
func (c *ReportContainer) GetComposedSimplerReportWithSimplerReport() *basicapp.ComposedSimplerReport {
	if c.composedSimplerReport == nil {
		c.composedSimplerReport = &basicapp.ComposedSimplerReport{
			SimplerReport: *c.GetSimplerReport(),
			Extra:         "spice",
		}
	}
	return c.composedSimplerReportWithSimplerReport
}

// GetProxyReporter makes things complicated
func (c *ReportContainer) GetProxyReporter() *basicapp.ProxyReporter {
	if c.proxyReporter == nil {
		c.proxyReporter = &basicapp.ProxyReporter{
			CSR: c.GetComposedSimplerReportWithSimplerReport(),
		}
	}
	return c.proxyReporter
}

// Get returns a service for its name
func (c *ReportContainer) Get(name string) interface{} {
	switch strings.ToLower(name) {
	case "parametersbag":
		return c.GetParametersBag()
	case "composedreport":
		return c.GetComposedReport()
	case "simplereporta":
		return c.GetSimpleReportA()
	case "simplereportb":
		return c.GetSimpleReportB()
	case "reportthatfails":
		return c.GetReportThatFails()
	case "reportthatnotfails":
		return c.GetReportThatNotFails()
	case "simplerreport":
		return c.GetSimplerReport()
	case "composedsimplerreport":
		return c.GetComposedSimplerReport()
	case "composedsimplerreportwithsimplerreport":
		return c.GetComposedSimplerReportWithSimplerReport()
	case "proxyreporter":
		return c.GetProxyReporter()
	default:
		return nil
	}
}
