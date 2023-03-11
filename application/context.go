package application

import (
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

/*
+-------------+
|   Context   |
+-------------+
*/

type ApplicationContext struct {
	*configurableContext
	*applicationArguments
	*configureEnvironment
}

func newApplicationContext(ctx *configurableContext, args *applicationArguments, env *configureEnvironment) *ApplicationContext {
	return &ApplicationContext{ctx, args, env}
}

type configurableContext struct {
	Pid             int
	ApplicationName string
	DisplayName     string
	GoVersion       string
	StartupDate     *time.Time
	instances       map[string]interface{}
	logger          *logrus.Logger
}

func newConfigurableContext(logger *logrus.Logger) *configurableContext {
	_startupDate := time.Now()
	_pid := os.Getegid()
	_name := GetProjectName()
	_goVersion := runtime.Version()
	_instances := make(map[string]interface{})
	return &configurableContext{
		Pid:             _pid,
		ApplicationName: _name,
		DisplayName:     _name,
		GoVersion:       _goVersion,
		StartupDate:     &_startupDate,
		instances:       _instances,
		logger:          logger,
	}
}
func (ctx *configurableContext) RegisterInstance(i interface{}) {
	iType := reflect.TypeOf(i).String()
	iType = strings.TrimLeft(iType, "*")
	ctx.instances[iType] = i
}
func (ctx *configurableContext) GetInstance(i interface{}) interface{} {
	iType := reflect.TypeOf(i).String()
	iType = strings.TrimLeft(iType, "*")
	value, _ := ctx.instances[iType]
	return value
}
func (ctx *configurableContext) GetOrDefaultInstance(i interface{}) interface{} {
	iType := reflect.TypeOf(i).String()
	iType = strings.TrimLeft(iType, "*")
	value, ok := ctx.instances[iType]
	if ok {
		return value
	}
	return ctx
}
func (ctx *configurableContext) GetInstanceAll() interface{} {
	var instances []interface{}
	for _, instance := range ctx.instances {
		instances = append(instances, instance)
	}
	return instances
}
func (ctx *configurableContext) Logger() *logrus.Logger {
	if ctx.logger != nil {
		ctx.logger = newLogger()
	}
	return ctx.logger
}
