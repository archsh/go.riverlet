package riverlet

type Worker[T Argument] interface {
	Run(arg T) error
	DefArgs() T
}
