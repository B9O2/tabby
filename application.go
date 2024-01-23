package tabby

import (
	"errors"
)

type Parameter struct {
	identify     string
	paramType    VarType
	defaultValue any
}

type Arguments struct {
	args    map[string]any
	appPath []string
}

func (a *Arguments) Get(key string) (any, error) {
	if v, ok := a.args[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("argument '" + key + "' not exists")
	}
}

func (a *Arguments) AppPath() []string {
	return a.appPath
}

func NewArguments(appPath []string, args map[string]any) Arguments {
	return Arguments{
		args:    args,
		appPath: appPath,
	}
}

type BaseApplication struct {
	apps  map[string]Application
	types map[string]Parameter
}

func (ba *BaseApplication) Init(Application) error {
	return nil
}

// SetParam default设置为nil则必须提供
func (ba *BaseApplication) SetParam(identify string, defaultValue any, argType VarType, alias ...string) {
	for _, alia := range append(alias, identify) {
		ba.types[alia] = Parameter{
			identify:     identify,
			paramType:    argType,
			defaultValue: defaultValue,
		}
	}
}

func (ba *BaseApplication) ParamTypes() map[string]Parameter {
	return ba.types
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
		apps:  make(map[string]Application),
		types: map[string]Parameter{},
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
	Main(Arguments) error
	SetParam(string, any, VarType, ...string)
	ParamTypes() map[string]Parameter
	SubApplication(string) (Application, bool)
}
