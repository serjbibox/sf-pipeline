package pipeline

type PipeUnit interface {
	Run() <-chan int
}
