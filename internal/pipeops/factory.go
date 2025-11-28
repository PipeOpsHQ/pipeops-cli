package pipeops

// Factory functions for creating clients, can be swapped for testing
var (
	NewClientFunc           = NewClient
	NewClientWithConfigFunc = NewClientWithConfig
)
