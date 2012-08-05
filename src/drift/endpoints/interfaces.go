package endpoints

type Endpoint interface {
	Start() bool
	Stop() bool
}
