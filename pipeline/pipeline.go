package pipeline

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
)

type Pipeline struct {
	pipeUnits []PipeUnit
	elog      *log.Logger
	ilog      *log.Logger
	done      <-chan struct{}
}

func PipelineInt(units ...PipeUnit) *Pipeline {
	l := log.New(os.Stdout, "PipeLine INFO\t", log.Ldate|log.Ltime)
	l.Println("Создан пайплайн с юнитами:")
	for _, val := range units {
		fmt.Printf("\t\t\t\t\t%v, \n", runtime.FuncForPC(reflect.ValueOf(val.(*PipeUnitInt).unitFunc).Pointer()).Name())
	}

	return &Pipeline{
		pipeUnits: units,
		elog:      log.New(os.Stdout, "PipeLine ERROR\t", log.Ldate|log.Ltime),
		ilog:      l,
	}

}

func (pl *Pipeline) Setup(done <-chan struct{}) *Pipeline {
	pl.done = done
	for idx := range pl.pipeUnits {
		switch pl.pipeUnits[idx].(type) {
		case *PipeUnitInt:
			pl.pipeUnits[idx].(*PipeUnitInt).done = done
		}
	}
	pl.ilog.Println("Условия завершения работы пайплайна настроены")
	return pl
}

func (pl *Pipeline) Run(source <-chan int) <-chan int {
	for idx := range pl.pipeUnits {
		switch pl.pipeUnits[idx].(type) {
		case *PipeUnitInt:
			pl.pipeUnits[idx].(*PipeUnitInt).input = source
			source = pl.pipeUnits[idx].(*PipeUnitInt).Run()
		}

	}
	pl.ilog.Println("Пайплайн запущен")
	return source
}
