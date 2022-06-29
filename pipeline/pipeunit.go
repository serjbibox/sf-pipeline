package pipeline

type UnitFunc func(<-chan struct{}, <-chan int) <-chan int

type PipeUnitInt struct {
	unitFunc UnitFunc
	input    <-chan int
	done     <-chan struct{}
}

func UnitInt(f UnitFunc) *PipeUnitInt {
	return &PipeUnitInt{unitFunc: f}
}

func (pu *PipeUnitInt) Run() <-chan int {
	return pu.unitFunc(pu.done, pu.input)
}
