package riverlet

var (
	defaultRuntime *Runtime
)

func RegisterWorker(worker Worker[Argument]) error {
	return defaultRuntime.RegisterWorker(worker)
}

func AddJob(job Job[Argument]) error {
	return defaultRuntime.AddJob(job)
}

func NewJob[T Argument](arg T) Job[T] {
	var job = Job[T]{
		Args: arg,
	}
	return job
}

func init() {
	defaultRuntime, _ = NewRuntime()
}
