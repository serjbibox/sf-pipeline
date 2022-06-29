package pipeline

type Pipeline struct {
	pipeUnits []PipeUnit
	done      <-chan struct{}
}

func PipelineInt(units ...PipeUnit) *Pipeline {
	return &Pipeline{pipeUnits: units}
}

func (pl *Pipeline) Setup(done <-chan struct{}) *Pipeline {
	pl.done = done
	for idx := range pl.pipeUnits {
		switch pl.pipeUnits[idx].(type) {
		case *PipeUnitInt:
			pl.pipeUnits[idx].(*PipeUnitInt).done = done
		}
	}
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
	return source
}
