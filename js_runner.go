package main

import (
	"fmt"
	"runtime/debug"

	"github.com/dop251/goja"
	"golang.org/x/sync/errgroup"
)

type Ruleset struct {
	tag      string
	url      string
	behavior string
	content  string
}

// downloadRulesets 并发下载规则集，保持调用顺序
// 从 JS 的 rulesets() 函数中提取规则集 URL，并发下载但按原始顺序返回
func downloadRulesets(vm *goja.Runtime) (resultLines []*Ruleset, err error) {
	jsRulesetsFunc := vm.Get("rulesets")

	if jsRulesetsFunc == nil || goja.IsUndefined(jsRulesetsFunc) || goja.IsNull(jsRulesetsFunc) {
		return
	}

	rulesetsFunc, ok := goja.AssertFunction(jsRulesetsFunc)
	if !ok {
		return nil, fmt.Errorf("rulesets must be a function")
	}

	registrations := make([]*Ruleset, 0, 64)
	register := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		behavior := "classical"
		if len(call.Arguments) >= 3 && !goja.IsUndefined(call.Argument(2)) {
			behavior = call.Argument(2).String()
		}
		registrations = append(registrations, &Ruleset{
			tag:      call.Argument(0).String(),
			url:      call.Argument(1).String(),
			behavior: behavior,
		})
		return goja.Undefined()
	})

	_, err = rulesetsFunc(goja.Undefined(), register)
	if err != nil {
		return nil, err
	}

	errGroup := new(errgroup.Group)
	limiter := make(chan bool, 8)
	resultLines = make([]*Ruleset, len(registrations))

	for i, registration := range registrations {
		i := i
		registration := registration
		errGroup.Go(func() error {
			limiter <- true
			defer func() {
				<-limiter
			}()

			content, e := GetOrPut(registration.url, FetchString)
			if e != nil {
				return e
			}

			registration.content = content
			resultLines[i] = registration
			return nil
		})
	}

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}

	return
}

// ExecJs 执行 JS 脚本，支持 rulesets 和 buildConfig 函数
func ExecJs(
	script string, template string, proxies SubscriptionData, legacyRelay bool,
) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[panic] %v\n%s", r, string(debug.Stack()))
		}
	}()

	vm := goja.New()
	err = vm.Set("log", func(v any) {
		L().Info(fmt.Sprintf("[JS] %v", v))
	})
	if err != nil {
		return
	}

	_, err = vm.RunString(script)
	if err != nil {
		return
	}

	ruleLines, err := downloadRulesets(vm)
	if err != nil {
		return
	}

	conf, err := BuildTemplate(template, proxies, ruleLines)
	if err != nil {
		return
	}

	buildConfigFunc := func(map[string]any, bool) {}
	jsBuildConfigFunc := vm.Get("buildConfig")
	if jsBuildConfigFunc != nil {
		err = vm.ExportTo(jsBuildConfigFunc, &buildConfigFunc)
		if err != nil {
			return
		}
	}

	buildConfigFunc(conf, legacyRelay)

	result, err = Marshal(conf)

	return
}
