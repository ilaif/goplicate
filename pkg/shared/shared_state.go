package shared

// State a shared state struct to pass state during goplicate sync
// between different project runs.
type State struct {
	Message string // Message the message for the change request
}
