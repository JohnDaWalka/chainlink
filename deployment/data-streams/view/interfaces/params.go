package interfaces

import "context"

// ContractView defines the base interface for any contract view
type ContractView interface {
	// SerializeView converts the view to string format like JSON
	SerializeView() (string, error)
}

// ContractViewGenerator is an interface type used mostly to standardize the view generator implementations
type ContractViewGenerator[C, P any, V ContractView] interface {
	Generate(ctx context.Context, contract C, params P) (V, error)
}
