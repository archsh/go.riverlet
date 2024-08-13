package riverlet

type Job[T Argument] struct {
	Args        T
	Priority    int
	Identity    string
	Attempts    int
	MaxAttempts int
	Status      STATUS
}
