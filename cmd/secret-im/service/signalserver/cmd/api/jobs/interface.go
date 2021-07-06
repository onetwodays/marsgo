package jobs

type Job interface {
	Start() error
	Stop()
}

