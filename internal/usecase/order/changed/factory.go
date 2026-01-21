package changed

import "courier-service/internal/model"

type Factory struct {
	processors map[model.OrderStatus]Processor
}

func NewFactory(processors map[model.OrderStatus]Processor) *Factory {
	return &Factory{processors: processors}
}

func (f *Factory) Get(status model.OrderStatus) (Processor, bool) {
	p, ok := f.processors[status]
	return p, ok
}
