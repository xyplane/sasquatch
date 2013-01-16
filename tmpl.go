package main

import "os"
import "io"
import "io/ioutil"
import "errors"
import "strings"
import htmltemplate "html/template" 
import texttemplate "text/template" 


var tmpl Template


var (
	ErrNotHTMLTmpl = errors.New("expecting HTML template")
	ErrNotTextTmpl = errors.New("expecting Text template")
)


var stdinUsed = false
var stdinAlias = "-"


type Template interface {
	New(string) Template
	Parse(string) error
	Lookup(string) Template
	Execute(io.Writer, map[string]interface{}) error
}

func newTemplate(name string, html bool) (Template, error) {
	if html {
		return newHTMLTemplate(name)
	}
	return newTextTemplate(name)
}


type HTMLTemplate struct {
	tmpl *htmltemplate.Template
}

func (ht *HTMLTemplate) New(name string) Template {
	return &HTMLTemplate{ ht.tmpl.New(name) }
}
 
func (ht *HTMLTemplate) Parse(source string) error {
	var _, err = ht.tmpl.Parse(source)
	return err
}

func (ht *HTMLTemplate) Lookup(name string) Template {
	var t = ht.tmpl.Lookup(name)
	if t == nil {
		return nil
	}
	return &HTMLTemplate{ t }
}

func (ht *HTMLTemplate) Execute(wr io.Writer, data map[string]interface{}) error {
	return ht.tmpl.Execute(wr, data)
}

func newHTMLTemplate(name string) (Template, error) {
	if tmpl == nil {
		var t *htmltemplate.Template
		t = htmltemplate.New(name)
		tmpl = &HTMLTemplate{t}
		return tmpl, nil
	}
	var _, check = tmpl.(*HTMLTemplate)
	if !check {
		return nil, ErrNotHTMLTmpl
	}
	var t = tmpl.Lookup(name)
	if t != nil {
		return t, nil
	}
	t = tmpl.New(name)
	return t, nil
}


type TextTemplate struct {
	tmpl *texttemplate.Template
}

func (tt *TextTemplate) New(name string) Template {
	return &TextTemplate{ tt.tmpl.New(name) }
}
 
func (tt *TextTemplate) Parse(source string) error {
	var _, err = tt.tmpl.Parse(source)
	return err
}

func (tt *TextTemplate) Lookup(name string) Template {
	var t = tt.tmpl.Lookup(name)
	if t == nil {
		return nil
	}
	return &TextTemplate{ t }
}

func (tt *TextTemplate) Execute(wr io.Writer, data map[string]interface{}) error {
	return tt.tmpl.Execute(wr, data)
}

func newTextTemplate(name string) (Template, error) {
	if tmpl == nil {
		var t *texttemplate.Template
		t = texttemplate.New(name)
		tmpl = &TextTemplate{t}
		return tmpl, nil
	}
	var _, check = tmpl.(*TextTemplate)
	if !check {
		return nil, ErrNotTextTmpl
	}
	var t = tmpl.Lookup(name)
	if t != nil {
		return t, nil
	}
	t = tmpl.New(name)
	return t, nil
}


//type AsyncTemplateParser func(chan<-error)

type TemplateSource interface {
	ReadTemplate(chan error) 
	AddTemplate() error
}


type templateFlag struct {
	Name string
	text string
	read bool
	src string
}

 
func (tf *templateFlag) ReadTemplate(rtn chan error) {
	rtn<- tf.ReadText()
}

func (tf *templateFlag) ReadText() error {
	if tf.read {
		return nil
	}
	tf.src = tf.text
	tf.read = true
	return nil
}

func (tf *templateFlag) AddTemplate() error {
	var readerr = tf.ReadText()
	if readerr != nil {
		return readerr
	}
	var t, err = newTemplate(tf.Name, htmlFlag)
	if err != nil {
		return err
	}
	return t.Parse(tf.src)
}


type templateFileArg struct {
	name string
	path string
	read bool
	src string
}

func newTemplateFileArg(arg string) (*templateFileArg, error) {
	var name = "stdin"
	var path = strings.TrimSpace(arg)
	if path == stdinAlias {
		if stdinUsed {
			return nil, ErrStdinAlreadyUsed
		}
		stdinUsed = true
	} else {
		var info, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
		name = info.Name()
	}
	return &templateFileArg{ name, path, false, "" }, nil
}

func (tf *templateFileArg) ReadTemplate(rtn chan error) {
	rtn<- tf.ReadPath()
}

func (tf *templateFileArg) ReadPath() error {
	if tf.read {
		return nil
	}
	var buf []byte
	var err error
	if tf.path == stdinAlias {
		buf, err = ioutil.ReadAll(os.Stdin)
	} else {
		buf, err = ioutil.ReadFile(tf.path)
	}
	if err != nil {
		return err
	}
	tf.src = string(buf)
	tf.read = true
	return nil
}

func (tf *templateFileArg) AddTemplate() error {
	var readerr = tf.ReadPath()
	if readerr != nil {
		return readerr
	}
	var t, err = newTemplate(tf.name, htmlFlag)
	if err != nil {
		return err
	}
	return t.Parse(tf.src)
}



//type templateFileFlag struct {
//	Name string
//	Path string
//}


