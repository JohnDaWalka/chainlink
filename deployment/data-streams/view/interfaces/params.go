package interfaces

// ContractView defines the base interface for any contract view
type ContractView interface {
	// SerializeView converts the view to a JSON string
	SerializeView() (string, error)
}
