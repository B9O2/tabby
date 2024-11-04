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
	IgnoreUnsupportedArgs() bool
}

type BaseApplication struct {
	apps                  map[string]Application
	params                []Parameter
	width, height         uint
	ignoreUnsupportedArgs bool
}

func (ba *BaseApplication) Init(Application) error {
	return nil
}

func (ba *BaseApplication) Help(parts ...string) {
	for _, part := range parts {
		fmt.Println(part)
	}
	params := map[string]string{}

	max := 0
	for _, param := range ba.params {
		alias := AddPrefix(param.alias, "-")
		key := fmt.Sprintf("   -%s (%s)", param.identify, param.defaultValue.String())
		value := fmt.Sprintf("%s(%s)", param.help, strings.Join(alias, ","))
		l := len(key)
		if l > max {
			max = l
		}
		params[key] = value
	}

	for k, v := range params {
		fmt.Println(k + ":" + strings.Repeat(" ", max-len(k)+2) + v)
	}

	if len(ba.apps) > 0 {
		fmt.Println("Subcommands:")
		for _, app := range ba.apps {
			name, desc := app.Detail()
			fmt.Printf("   %s : %s\n", name, desc)
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

func (ba *BaseApplication) IgnoreUnsupportedArgs() bool {
	return ba.ignoreUnsupportedArgs
}

func NewBaseApplication(ignoreUnsupportedArgs bool, width, height uint, apps []Application) *BaseApplication {
	ba := &BaseApplication{
		apps:                  make(map[string]Application),
		width:                 width,
		height:                height,
		ignoreUnsupportedArgs: ignoreUnsupportedArgs,
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
