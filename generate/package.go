package generate

import "github.com/lucassabreu/go-container/scan"

// Package represents a package to be imported by the container
type Package struct {
	name     string
	fullName string
	alias    *string
	scanned  scan.Package
}

// UniqueName package name or alias (if informed), should be unique
func (p Package) UniqueName() string {
	if p.alias != nil {
		return *p.alias
	}

	return p.name
}

// Path is the import path of the package
func (p Package) Path() string {
	return p.fullName
}

// ScannedPackage returns the scan.Package which it is based
func (p Package) ScannedPackage() scan.Package {
	return p.scanned
}
