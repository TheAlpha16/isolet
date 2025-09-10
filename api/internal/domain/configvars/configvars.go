package configvars

type ConfigKey[T any] struct {
	Name    string
	Default T
}

var (
	TestKey = ConfigKey[string]{"test_key", "hellobruh"}
)
