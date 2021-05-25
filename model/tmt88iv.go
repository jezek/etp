package model

const Model = "TM-T88IV"

type Tmt88iv struct{}

func (m Tmt88iv) Model() string {
	return "TMT88IV"
}

func (m Tmt88iv) Commands() map[string]Command {
	res := map[string]Command{}
	for k, v := range CommonCommands {
		res[k] = v
	}

	delete(res, "ds")
	delete(res, "nods")
	return res
}
