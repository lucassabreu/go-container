package def

// Container structed
type Container struct {
	// Packages that will be used and their aliases (when needed)
	Packages map[string]*string

	// Services to be managed by the container
	Services map[string]Service

	// Services aliases
	Aliases map[string]string
}

// Service represents a entry in the services config
type Service struct {
	Factory   *string
	Arguments *[]interface{}
	Struct    *string
	InitMap   *map[string]interface{}
}

// IsByFactory when a factory function factory should be called
func (s Service) IsByFactory() bool {
	return s.Factory != nil
}

// IsByInitialization when the struct should be initialized
func (s Service) IsByInitialization() bool {
	return s.Struct != nil
}
