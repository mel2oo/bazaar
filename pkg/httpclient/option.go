package httpclient

import (
	"sync"
	"time"
)

var (
	optionPool = &sync.Pool{
		New: func() interface{} {
			return &option{
				header: make(map[string][]string),
			}
		},
	}
)

// Option 自定义设置http请求
type Option func(*option)

type option struct {
	ttl         time.Duration
	header      map[string][]string
	retryTimes  int
	retryDelay  time.Duration
	retryVerify RetryVerify
}

func (o *option) reset() {
	o.ttl = 0
	o.header = make(map[string][]string)
	o.retryTimes = 0
	o.retryDelay = 0
	o.retryVerify = nil
}

func getOption() *option {
	opt := optionPool.Get().(*option)
	opt.reset()
	return opt
}

func releaseOption(opt *option) {
	if opt != nil {
		opt.reset()
		optionPool.Put(opt)
	}
}

// WithTTL 本次http请求最长执行时间
func WithTTL(ttl time.Duration) Option {
	return func(opt *option) {
		opt.ttl = ttl
	}
}

// WithHeader 设置http header，可以调用多次设置多对key-value
func WithHeader(key, value string) Option {
	return func(opt *option) {
		opt.header[key] = []string{value}
	}
}

// WithOnFailedRetry 设置失败重试
func WithOnFailedRetry(retryTimes int, retryDelay time.Duration, retryVerify RetryVerify) Option {
	return func(opt *option) {
		opt.retryTimes = retryTimes
		opt.retryDelay = retryDelay
		opt.retryVerify = retryVerify
	}
}
