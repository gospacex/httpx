package httpx

import "fmt"

type AdapterFactory func() Server

var adapterRegistry = make(map[string]AdapterFactory)

func RegisterAdapter(name string, factory AdapterFactory) {
	adapterRegistry[name] = factory
}

func getAdapter(name string) (AdapterFactory, error) {
	factory, ok := adapterRegistry[name]
	if !ok {
		return nil, fmt.Errorf("adapter %q not found, available adapters: gin, hertz, nethttp", name)
	}
	return factory, nil
}