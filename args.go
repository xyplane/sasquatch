package main

import "os"
import "fmt"
import "flag"
import "errors"
import "strings"



var (
	ErrKeyRequired = errors.New("key required")
	ErrValueRequired = errors.New("value required")
	ErrStdinAlreadyUsed = errors.New("stdin already used")
)


var (
	vFlag bool
	htmlFlag bool
	versionFlag bool
	dataFlags []DataSource
	//executeFlags []Executers
	templateFlags []TemplateSource

)


var KeyValueFlagSeparators = []string{ "=", ":" }


func init() {
	flag.BoolVar(&vFlag, "v", false, "Verbose")
	flag.BoolVar(&htmlFlag, "html", false, "Use HTML safe template")
	flag.Var(&dataFlagValue{}, "d", "Define a new value")
	flag.Var(&templateFlagValue{}, "t", "Define a new template")
}


func parseArguments() error {
	flag.Parse()

	var templateFileArgs []TemplateSource

	if flag.NArg() == 0 {
		if !stdinUsed {
			var t, err = newTemplateFileArg(stdinAlias)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			templateFileArgs = append(templateFileArgs, t)
		}
	} else {
		for _, arg := range flag.Args() {
			var t, err = newTemplateFileArg(arg)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			templateFileArgs = append(templateFileArgs, t)
		}
	}

	templateFlags = append(templateFileArgs, templateFlags...)
	return nil
}

func splitKeyValueFlag(value string, keyReq bool, valReq bool) (string, string, error) {
	var idx = -1
	var sep = ""
	for _, s := range KeyValueFlagSeparators {
		var i = strings.Index(value, s)
		if (i > -1) && ((idx == -1) || (idx > i)) {
			idx = i
			sep = s
		}
	}
	var kv []string
	if idx == -1 {
		kv = []string{ "", value }
	} else {
		kv = strings.SplitN(value, sep, 2)
	}
	if keyReq && (kv[0] == "") {
		return kv[0], kv[1], ErrKeyRequired
	}
	if valReq && (kv[1] == "") {
		return kv[0], kv[1], ErrValueRequired
	}
	return kv[0], kv[1], nil
}


type dataFlagValue struct {}

func (*dataFlagValue) String() string {
	return "key=value"
}

func (*dataFlagValue) Set(value string) error {
	var k, v, err = splitKeyValueFlag(value, true, false)
	if err != nil {
		return err
	}
	var dFlag = &dataFlag{ k, v }
	dataFlags = append(dataFlags, dFlag)
	return nil
}


type templateFlagValue struct {}

func (*templateFlagValue) String() string {
	return "key=value"
}

func (*templateFlagValue) Set(value string) error {
	var k, v, err = splitKeyValueFlag(value, true, true)
	if err != nil {
		return err
	}
	var tFlag = &templateFlag{ k, v, false, "" }
	templateFlags = append(templateFlags, tFlag)
	return nil
}

