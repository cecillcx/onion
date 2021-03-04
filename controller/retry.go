package controller

import (
	"onion"
)

type RetryCtrl struct {
	retryTimes int
	errorCheck func([]interface{}) error
}

func NewRetryCtrl(maxRetryTimes int, errChecker func([]interface{}) error) onion.Controller {
	c := RetryCtrl{}
	c.retryTimes = maxRetryTimes
	c.errorCheck = errChecker
	return &c
}

func (r *RetryCtrl) Run(executor *onion.ExecuteContext, latestRet []interface{}) []interface{} {
	for i := 0; i < r.retryTimes; i++ {
		if r.errorCheck(latestRet) == nil {
			break
		}
		latestRet = executor.ReExecute()
	}
	return latestRet
}

func (r *RetryCtrl) GetControllerPos() onion.CtrlPos {
	return onion.After
}
