package blademaster

import (
	"net/url"
	"regexp"
	"strings"

	"marsgo/pkg/log"
)


//中间件.根据Host和url来判断是否跨域

func matchHostSuffix(suffix string) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return strings.HasSuffix(strings.ToLower(uri.Host), suffix)
	}
}

func matchPattern(pattern *regexp.Regexp) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return pattern.MatchString(strings.ToLower(uri.String()))
	}
}

// CSRF returns the csrf middleware to prevent invalid cross site request.
// Only referer is checked currently.
func CSRF(allowHosts []string, allowPattern []string) HandlerFunc {
	validations := []func(*url.URL) bool{}

	addHostSuffix := func(suffix string) {
		validations = append(validations, matchHostSuffix(suffix))
	}
	addPattern := func(pattern string) {
		validations = append(validations, matchPattern(regexp.MustCompile(pattern)))
	}

	for _, r := range allowHosts {
		addHostSuffix(r)
	}
	for _, p := range allowPattern {
		addPattern(p)
	}

	return func(c *Context) {
		referer := c.Request.Header.Get("Referer") //当浏览器向web服务器发送请求的时候，一般会带上Referer，告诉服务器该网页是从哪个页面链接过来的，服务器因此可以获得一些信息用于处理。
		if referer == "" {
			log.V(5).Info("The request's Referer or Origin header is empty.")
			c.AbortWithStatus(403)
			return
		}
		illegal := true
		if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
			for _, validate := range validations {
				if validate(uri) {
					illegal = false
					break
				}
			}
		}
		if illegal {
			log.V(5).Info("The request's Referer header `%s` does not match any of allowed referers.", referer)
			c.AbortWithStatus(403)
			return
		}
	}
}
