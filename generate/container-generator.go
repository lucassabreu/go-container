package generate

// ContainerGenerator generate golang code for the a Contaienr based on commands
type ContainerGenerator struct {
	ContainerPackage string
	ContainerName    string

	Packages map[string]string
	Services map[string]service
}
