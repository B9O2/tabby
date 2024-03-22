package tabby

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Int = NewTransfer("int", func(s string) (any, error) {
		if i, err := strconv.ParseInt(s, 10, 0); err != nil {
			return 0, err
		} else {
			return int(i), nil
		}
	})

	String = NewTransfer("string", func(s string) (any, error) {
		return s, nil
	})

	Bool = NewTransfer("int", func(s string) (any, error) {
		if len(s) == 0 {
			return true, nil
		} else {
			return false, errors.New("boolean requires no value")
		}
	})
)

type DefaultValue struct {
	name     string
	value    any
	transfer func(string) (any, error)
}

func (dv DefaultValue) String() string {
	switch v := dv.value.(type) {
	case nil:
		return "*"
	case string:
		return "\"" + v + "\""
	default:
		return fmt.Sprint(dv.value)
	}
}

func DefaultParser(rawArgs []string) (map[string]string, error) {
	//RawArgs
	currentKey := ""
	strArgs := map[string]string{}
	for _, argv := range rawArgs {
		if len(argv) <= 0 {
			continue
		}
		if argv[0] == '-' {
			key := argv[1:]
			if currentKey != "" { //未分配值
				strArgs[currentKey] = ""
			}
			currentKey = key
		} else {
			strArgs[currentKey] = argv
			currentKey = ""
		}
	}
	if currentKey != "" {
		strArgs[currentKey] = ""
	}
	return strArgs, nil
}

type Tabby struct {
	name                string
	mainApp, unknownApp Application
	parser              func([]string) (map[string]string, error)
}

func (t *Tabby) SetParser(parser func([]string) (map[string]string, error)) error {
	if parser != nil {
		t.parser = parser
		return nil
	} else {
		return errors.New("parser is nil")
	}
}

func (t *Tabby) Run(rawArgs []string) (*TabbyContainer, error) {
	if rawArgs == nil {
		rawArgs = os.Args[1:]
	}

	//Apps
	var apps []string
	i := 0
	for _, argv := range rawArgs {

		if len(argv) > 0 && argv[0] == '-' {
			break
		} else {
			apps = append(apps, argv)
		}
		i += 1
	}
	rawArgs = rawArgs[i:]

	//App
	app := t.mainApp
	var subApp Application
	ok := false
	name, _ := t.mainApp.Detail()
	if err := app.Init(nil); err != nil {
		return nil, errors.New("error: '" + name + "' cause:" + err.Error())
	}
	appPath := []string{name}
	for _, appName := range apps {
		appPath = append(appPath, appName)
		subApp, ok = app.SubApplication(appName)
		if !ok {
			if t.unknownApp == nil {
				return nil, fmt.Errorf("App '%s' not exists", strings.Join(appPath, "/"))
			} else {
				app = t.unknownApp
				if err := app.Init(t.mainApp); err != nil {
					appName, _ := app.Detail()
					return nil, errors.New("error: '" + name + "/" + appName + "' cause:" + err.Error())
				}
				break
			}
		}
		if err := subApp.Init(app); err != nil {
			return nil, errors.New("error: '" + strings.Join(appPath, "/") + "' cause:" + err.Error())
		}
		app = subApp
	}

	finalAppPath := strings.Join(appPath, "/")

	//Args
	empty := true
	strArgs, err := t.parser(rawArgs)
	if err != nil {
		return nil, err
	}

	params := app.Params()
	args := map[string]any{}

	for _, param := range params {
		for _, alia := range param.alias {
			if v, ok := strArgs[alia]; ok {
				if value, err1 := param.defaultValue.transfer(v); err1 != nil {
					return nil, fmt.Errorf("App '%s': argument '%s(%s)' :error: %s", finalAppPath, param.identify, param.defaultValue.name, err1.Error())
				} else {
					args[param.identify] = value
					empty = false
				}
				delete(strArgs, alia)
			}
		}
	}

	if len(strArgs) > 0 {
		return nil, fmt.Errorf(
			"App '%s': unsupported parameters '%s'",
			finalAppPath,
			strings.Join(AddPrefix(MapKeys[string, string](strArgs), "-"), ","))
	}

	//DefaultArgs
	for _, param := range params {
		if _, ok := args[param.identify]; !ok {
			if param.defaultValue.value != nil {
				args[param.identify] = param.defaultValue.value
			} else if !empty {
				return nil, fmt.Errorf(
					"App '%s': required parameter '%s' not provided(%s)",
					finalAppPath, param.identify,
					strings.Join(AddPrefix(param.alias, "-"), ","))
			}
		}
	}

	//Run
	tc, err := app.Main(NewArguments(empty, appPath, args))
	if err != nil {
		return nil, errors.New("App '" + finalAppPath + "' error:" + err.Error())
	}
	return tc, nil
}

func (t *Tabby) SetUnknownApp(app Application) {
	t.unknownApp = app
}

func NewTabby(name string, mainApp Application) *Tabby {
	t := &Tabby{
		name:    name,
		mainApp: mainApp,
		parser:  DefaultParser,
	}

	return t
}

func NewTransfer(name string, transfer func(string) (any, error)) func(any) DefaultValue {
	if transfer != nil {
		return func(a any) DefaultValue {
			return DefaultValue{
				name:     name,
				value:    a,
				transfer: transfer,
			}
		}
	} else {
		panic("transfer '" + name + "' is nil")
	}
}
