package processor

import "github.com/Alonza0314/nf-example/pkg/app"

type ProcessorNf interface {
	app.App

	Processor() *Processor
}

type Processor struct {
	ProcessorNf
}

func NewProcessor(nf ProcessorNf) (*Processor, error) {
	p := &Processor{
		ProcessorNf: nf,
	}
	return p, nil
}
