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
	} else {
		panic("argument '" + key + "' not registered")
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

type BaseApplication struct {
	apps   map[string]Application
	params []Parameter
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
		fmt.Printf("   -%s %s(%s)\n", param.identify, param.help, strings.Join(alias, ","))
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

func NewBaseApplication(apps []Application) *BaseApplication {
	ba := &BaseApplication{
		apps: make(map[string]Application),
	}
	for _, app := range apps {
		appName := app.Name()
		if _, ok := ba.apps[appName]; ok {
			panic("'" + appName + "' exists")
		} else {
			ba.apps[appName] = app
		}
	}
	return ba
}

type Application interface {
	Init(Application) error
	Name() string
	Help(...string)
	Main(Arguments) error
	SetParam(string, string, DefaultValue, ...string)
	Params() []Parameter
	SubApplication(string) (Application, bool)
}
