package interfaces

// ISource defines an interface for navigator data source
type ISource interface {
	Get(string) (any, bool)
	Set(string, any) bool
}

// Interface represents qualifier to access fields of navigator
type IQual interface {
	// Function to access next part of qualifier
	Next() (string, bool)
	// Function to reset qualifier back to zero offset start state.
	Reset()
	// Function to get an independent copy of current qual
	Copy() IQual
}
