package etp

import (
	"bytes"
	"etp/model"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

type Printer struct {
	Model            model.Printer
	Template         string
	SuppressCommands bool

	encoder *encoding.Encoder
}

func New(modelName, template string) (Printer, error) {
	m, err := model.New(modelName)
	if err != nil {
		return Printer{}, err
	}

	if !utf8.Valid([]byte(template)) {
		return Printer{}, fmt.Errorf("not valid utf8 template")
	}

	return Printer{
		m,
		template,
		false,
		charmap.CodePage852.NewEncoder(),
	}, nil
}

func (p Printer) WriteTo(w io.Writer, data interface{}) error {
	t, err := p.NewPrintTemplate("WriteTo")
	if err != nil {
		return err
	}
	data, err = p.encode(data)
	if err != nil {
		return fmt.Errorf("data encoding: %s", err)
	}

	buf := &bytes.Buffer{}
	if _, err := buf.Write(p.initCommands()); err != nil {
		return err
	}
	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("template execute: %s", err)
	}

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func (p Printer) initCommands() []byte {
	if p.SuppressCommands {
		return nil
	}
	//TODO code page
	// this is latin-2
	return append(
		[]byte(model.Init()),
		model.CP852()...,
	)
}

func (p Printer) encode(data interface{}) (interface{}, error) {
	//TODO better encoding ... use reflect
	switch d := data.(type) {
	case string:
		s, err := p.encoder.String(d)
		if err != nil {
			return nil, err
		}
		return s, nil
	case []interface{}:
		s := make([]interface{}, len(d))
		for i := range d {
			ret, err := p.encode(d[i])
			if err != nil {
				return nil, err
			}
			s[i] = ret
		}
		return s, nil
	case map[string]interface{}:
		m := make(map[string]interface{}, len(d))
		for k, v := range d {
			ret, err := p.encode(v)
			if err != nil {
				return nil, err
			}
			m[k] = ret
		}
		return m, nil
	}
	return data, nil
}

func (p Printer) NewTemplate(name string) (*template.Template, error) {
	t := template.New(name)
	p.AddFuncs(t)
	return t.Parse(p.Template)
}

func (p Printer) NewPrintTemplate(name string) (*template.Template, error) {
	content, err := p.encoder.String(p.Template)
	if err != nil {
		return nil, fmt.Errorf("encoding: %s", err)
	}

	t := template.New(name)

	p.AddFuncs(t)

	content = p.ConvertNewLines(content)

	return t.Parse(content)
}

var helpers map[string]model.Command = map[string]model.Command{
	"raw": {
		"(...byte) Write bytes without encoding",
		func(b ...byte) string {
			return fmt.Sprint(string(b))
		},
	},
	"esc": {
		"(...byte) Write escape command followed by bytes of your choice (equals \"raw 0x1b ...byte\")",
		func(b ...byte) string {
			return string(byte(0x1b)) + string(b)
		},
	},
	"rawFile": {
		"(string) Write file content without encoding. If file read fails, return empty string.",
		func(file string) string {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return ""
			}
			return string(data)
		},
	},
}

func (p Printer) Commands() map[string]model.Command {
	commands := map[string]model.Command{}

	for k, command := range helpers {
		commands[k] = command
	}

	mc := model.CommonCommands

	if p.Model != nil {
		mc = p.Model.Commands()
	}
	for k, command := range mc {
		commands[k] = command
	}
	return commands
}

func (p Printer) AddFuncs(t *template.Template) {
	tf := map[string]interface{}{}

	for k, command := range p.Commands() {
		if p.SuppressCommands {
			tf[k] = func() string { return "" }
			continue
		}
		tf[k] = command.Function
	}

	t.Funcs(tf)
}

func (p Printer) ConvertNewLines(s string) string {
	if mc, ok := p.Model.(interface {
		ConvertNewLines(string) string
	}); ok {
		return mc.ConvertNewLines(s)
	}

	if strings.IndexByte(s, '\n') == -1 {
		return s
	}

	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\n", "\r\n", -1)
	return s
}
