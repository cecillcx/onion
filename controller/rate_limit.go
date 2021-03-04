package controller

import (
	"fmt"
	"onion"
	"time"

	limiter "github.com/juju/ratelimit"
)

type RateLimiterCtrl struct {
	QPS         int64
	rateLimiter *limiter.Bucket
}

func NewRateLimiterCtrl(qps int64) onion.Controller {
	r := RateLimiterCtrl{}
	r.QPS = qps
	r.rateLimiter = limiter.NewBucketWithQuantum(1*time.Second, r.QPS, r.QPS)
	return &r
}

func (r *RateLimiterCtrl) Run(*onion.ExecuteContext, []interface{}) []interface{} {
	r.rateLimiter.Wait(1)
	return nil
}

func (r *RateLimiterCtrl) GetControllerPos() onion.CtrlPos {
	return onion.Before
}

func main() {
	fmt.Println("vim-go")
}
