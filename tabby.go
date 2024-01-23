package tabby

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type VarType int

const (
	STRING VarType = iota
	INT
	BOOL
)

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

func (t *Tabby) Run(rawArgs []string) error {
	if rawArgs == nil {
		rawArgs = os.Args[1:]
	}
	//Apps
	var apps []string
	i := 0
	for _, argv := range rawArgs {
		if argv[0] == '-' {
			break
		}
		apps = append(apps, argv)
		i += 1
	}
	rawArgs = rawArgs[i:]

	//App
	app := t.mainApp
	var subApp Application
	ok := false
	if err := app.Init(nil); err != nil {
		return errors.New("error: '" + t.mainApp.Name() + "' cause:" + err.Error())
	}
	appPath := []string{t.mainApp.Name()}
	for _, appName := range apps {
		appPath = append(appPath, appName)
		subApp, ok = app.SubApplication(appName)
		if !ok {
			if t.unknownApp == nil {
				return errors.New(fmt.Sprintf("App '%s' not exists", strings.Join(appPath, "/")))
			} else {
				app = t.unknownApp
				break
			}
		}
		if err := subApp.Init(app); err != nil {
			return errors.New("error: '" + strings.Join(appPath, "/") + "' cause:" + err.Error())
		}
		app = subApp
	}

	finalAppPath := strings.Join(appPath, "/")

	//Args
	strArgs, err := t.parser(rawArgs)
	if err != nil {
		return err
	}

	params := app.ParamTypes()
	args := map[string]any{}
	for k, v := range strArgs {
		if param, ok := params[k]; ok {
			switch param.paramType {
			case STRING:
				args[param.identify] = v
			case INT:
				args[param.identify], err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
			case BOOL:
				if len(v) == 0 {
					args[param.identify] = true
				} else {
					return errors.New(fmt.Sprintf("App '%s': argument '%s' is boolean.", finalAppPath, param.identify))
				}
			default:
				return errors.New(fmt.Sprintf("unknown argument type id '%d'", param.paramType))
			}
		} else {
			return errors.New(fmt.Sprintf("App '%s': argument '%s' not support", finalAppPath, k))
		}
	}
	//------Default Args
	for _, param := range params {
		if _, ok := args[param.identify]; !ok {
			if param.defaultValue != nil {
				args[param.identify] = param.defaultValue
			} else {
				var alias []string
				for alia, param1 := range params {
					if param1.identify == param.identify {
						alias = append(alias, "-"+alia)
					}
				}
				return errors.New(fmt.Sprintf("App '%s': required parameter '%s' not provided(%s)", finalAppPath, param.identify, strings.Join(alias, ",")))
			}
		}
	}
	//Run
	err = app.Main(NewArguments(appPath, args))
	if err != nil {
		return errors.New("App '" + finalAppPath + "' error:" + err.Error())
	}

	return nil
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
