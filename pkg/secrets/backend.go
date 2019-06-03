package secrets

import (
	"fmt"
	"sync"
)

type Backend interface {
	Init(map[string]string) error
	Get(string) (string, error)
}

// BackendInstances are instantiated backends
var BackendInstances map[string]Backend

// BackendFunctions is a map of functions that return Backends
var BackendFunctions map[string]func() Backend

var initLock sync.Mutex

// BackendInstantiate instantiates a Backend of type `backendType`
func BackendInstantiate(name string, backendType string) error {
	if BackendInstances == nil {
		BackendInstances = make(map[string]Backend)
	}

	function, found := BackendFunctions[backendType]
	if !found {
		return fmt.Errorf("Unkown backend type: %v", backendType)
	}

	BackendInstances[name] = function()

	return nil
}

// BackendRegister registers a new backend type with name `name`
// function is a function that returns a backend of that type
func BackendRegister(name string, function func() Backend) {
	if BackendFunctions == nil {
		BackendFunctions = make(map[string]func() Backend)
	}

	BackendFunctions[name] = function
}

// BackendInitFromEnv initializes a backend looking into Env for config data
func BackendInitFromEnv() error {
	initLock.Lock()
	defer initLock.Unlock()

	config, err := BackendConfigFromEnv()
	if err != nil {
		return err
	}

	err = BackendInstantiate(config.Name, config.Type)
	if err != nil {
		return err
	}

	err = BackendInstances[config.Name].Init(config.Parameters)

	return err
}
