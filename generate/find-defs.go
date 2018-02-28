package generate

type defLooker struct {
	funcs map[string]map[string]interface{}
}

// DefLooker is a helper to find definitions
var DefLooker = defLooker{
	funcs: make(map[string]map[string]interface{}),
}

func (dl defLooker) getFunc(pkg, funcName string) error {
	return nil
}
