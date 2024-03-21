package tabby

import (
	"fmt"
	"strings"
)

type Parameter struct {
	identify, help string
	alias          []string
	defaultValue   DefaultValue
}

type Arguments struct {
	empty   bool
	args    map[string]any
	appPath []string
}

func (a *Arguments) Get(key string) any {
	if v, ok := a.args[key]; ok {
		return v
	} else if !a.empty {
		panic("argument '" + key + "' not registered")
	} else {
		panic("argument '" + key + "' is required,please handle empty arguments")
	}
}

func (a *Arguments) AppPath() []string {
	return a.appPath
}

func (a *Arguments) IsEmpty() bool {
	return a.empty
}

func NewArguments(empty bool, appPath []string, args map[string]any) Arguments {
	return Arguments{
		args:    args,
		appPath: appPath,
		empty:   empty,
	}
}

type Application interface {
	Init(Application) error
	Detail() (string, string)
	Help(...string)
	Main(Arguments) (*TabbyContainer, error)
	SetParam(string, string, DefaultValue, ...string)
	Params() []Parameter
	SubApplication(string) (Application, bool)
}

type BaseApplication struct {
	apps          map[string]Application
	params        []Parameter
	width, height uint
}

func (ba *BaseApplication) Init(Application) error {
	return nil
}

func (ba *BaseApplication) Help(parts ...string) {
	for _, part := range parts {
		fmt.Println(part)
	}

	for _, param := range ba.params {
		alias := AddPrefix(param.alias, "-")
		fmt.Printf("   -%s | %s(%s)\n", param.identify, param.help, strings.Join(alias, ","))
	}

	if len(ba.apps) > 0 {
		fmt.Println("Subcommands:")
		for _, app := range ba.apps {
			name, desc := app.Detail()
			fmt.Printf("   %s | %s\n", name, desc)
		}
	}

}

// SetParam default设置为nil则必须提供
func (ba *BaseApplication) SetParam(identify, help string, defaultValue DefaultValue, alias ...string) {
	ba.params = append(ba.params, Parameter{
		identify:     identify,
		help:         help,
		alias:        append(alias, identify),
		defaultValue: defaultValue,
	})
}

func (ba *BaseApplication) Params() []Parameter {
	return ba.params
}

func (ba *BaseApplication) SubApplication(name string) (Application, bool) {
	if app, ok := ba.apps[name]; ok {
		return app, true
	} else {
		return nil, false
	}
}

func (ba *BaseApplication) SubApplications() map[string]Application {
	return ba.apps
}

func NewBaseApplication(width, height uint, apps []Application) *BaseApplication {
	ba := &BaseApplication{
		apps:   make(map[string]Application),
		width:  width,
		height: height,
	}
	for _, app := range apps {
		appName, _ := app.Detail()
		if _, ok := ba.apps[appName]; ok {
			panic("'" + appName + "' exists")
		} else {
			ba.apps[appName] = app
		}
	}
	return ba
}
