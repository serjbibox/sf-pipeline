package pipeline

import (
	"log"
	"os"
	"reflect"
	"runtime"
)

type UnitFunc func(<-chan struct{}, <-chan int) <-chan int

type PipeUnitInt struct {
	unitFunc UnitFunc
	elog     *log.Logger
	ilog     *log.Logger
	input    <-chan int
	done     <-chan struct{}
}

func NewUnitInt(f UnitFunc) *PipeUnitInt {
	l := log.New(os.Stdout, "Unit INFO\t", log.Ldate|log.Ltime)
	l.Printf("Юнит с функцией %v создан", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	return &PipeUnitInt{
		unitFunc: f,
		elog:     log.New(os.Stdout, "Unit ERROR\t", log.Ldate|log.Ltime),
		ilog:     l,
	}
}

func (pu *PipeUnitInt) Run() <-chan int {
	pu.ilog.Printf("Юнит с функцией %v запущен", runtime.FuncForPC(reflect.ValueOf(pu.unitFunc).Pointer()).Name())
	return pu.unitFunc(pu.done, pu.input)
}

func (pu *PipeUnitInt) SetDone(done <-chan struct{}) {
	pu.done = done
}
func (pu *PipeUnitInt) SetInput(input <-chan int) {
	pu.input = input
}
