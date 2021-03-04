package onion

import (
	"reflect"
)

type ExecuteContext struct {
	orinfunc  interface{}
	args      []reflect.Value
	ctrlsDone []Controller
}

func (e *ExecuteContext) SetArgs(args []reflect.Value) {
	e.args = args
}

func (e *ExecuteContext) SetFunc(orifunc interface{}) {
	e.orinfunc = orifunc
}

func (e *ExecuteContext) execute() []interface{} {
	funcValue := reflect.ValueOf(e.orinfunc)
	var res []interface{}
	for _, r := range funcValue.Call(e.args) {
		res = append(res, r.Interface())
	}
	return res
}

func (e *ExecuteContext) ReExecute() []interface{} {
	for _, c := range e.ctrlsDone {
		c.Run(e, nil)
	}

	return e.execute()
}

func (e *ExecuteContext) AddCtrlDone(c Controller) {
	e.ctrlsDone = append(e.ctrlsDone, c)
}

type Onion struct {
	executor    *ExecuteContext
	orinfunc    interface{}
	beforeCtrls []Controller
	afterCtrls  []Controller
}

func NewOnion() Onion {
	return Onion{}
}

const (
	Before = -1
	After  = 1
)

type CtrlPos int

// 用户可自行实现相关的控制器逻辑
type Controller interface {
	Run(*ExecuteContext, []interface{}) []interface{}

	// 提供两个位置的插入，即执行前&执行后
	// 执行前返回 -1，执行后返回 1
	GetControllerPos() CtrlPos
}

func (o *Onion) ApplyFunc(ori interface{}) *Onion {
	o.orinfunc = ori
	return o
}

// func (o *Onion) invoke(args ...interface{}) []reflect.Value {
func (o *Onion) Invoke(args ...interface{}) []interface{} {
	// reset context and env
	o.afterCtrls = nil
	o.beforeCtrls = nil
	o.executor = &ExecuteContext{}

	params := []reflect.Value{}
	for _, a := range args {
		params = append(params, reflect.ValueOf(a))
	}
	o.executor.SetFunc(o.orinfunc)
	o.executor.SetArgs(params)

	for _, c := range o.beforeCtrls {
		c.Run(o.executor, nil)
		o.executor.AddCtrlDone(c)
	}

	ret := o.executor.execute()

	for _, c := range o.afterCtrls {
		ret = c.Run(o.executor, ret)
		o.executor.AddCtrlDone(c)
	}

	return ret
}

func (o *Onion) AddCtrl(ctrl Controller, args ...interface{}) *Onion {
	if ctrl.GetControllerPos() == Before {
		o.beforeCtrls = append(o.beforeCtrls, ctrl)
	} else {
		o.afterCtrls = append(o.afterCtrls, ctrl)
	}

	return o
}
