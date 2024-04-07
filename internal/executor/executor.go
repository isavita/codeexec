package executor

type Executor interface {
	Execute(code, language string) (string, error)
}
