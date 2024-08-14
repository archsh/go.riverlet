package riverlet

var (
	defaultRuntime *Runtime
)

// RegisterWorker register a worker to builtin runtime
func RegisterWorker(worker Worker[Argument]) error {
	return defaultRuntime.RegisterWorker(worker)
}

// AddJob add a new job to builtin runtime
func AddJob(job Job[Argument]) error {
	return defaultRuntime.AddJob(job)
}

// Start builtin runtime
func Start() error {
	return defaultRuntime.Start()
}

// Stop builtin runtime
func Stop() error {
	return defaultRuntime.Stop()
}

func init() {
	defaultRuntime, _ = NewRuntime()
}
