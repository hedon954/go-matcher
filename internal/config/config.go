package config

// Configer provides a way to get the config.
type Configer[T any] interface {
	Get() *T
}
