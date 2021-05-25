package model

import "fmt"

type Printer interface {
	Model() string
	Commands() map[string]Command
}

var availableModels map[string]func() Printer = map[string]func() Printer{
	"TM-T88IV": func() Printer { return Tmt88iv{} },
}

type errModelNotSupported struct {
	model string
}

func (e errModelNotSupported) Error() string {
	return fmt.Sprintf("Model \"%s\" not supported", e.model)
}

func New(model string) (Printer, error) {
	mc, ok := availableModels[model]
	if !ok {
		return nil, errModelNotSupported{model}
	}
	m := mc()

	return m, nil
}
