package util

import (
	"crypto/tls"
	"github.com/FFFFFFStatic/resty"
	"time"
)

func DoGet(fullUrl string, cookie string) (*resty.Response, error) {
	client := resty.New()
	client.
		SetRetryCount(0).
		SetRetryWaitTime(time.Duration(100) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(2000) * time.Millisecond).
		SetTimeout(time.Duration(3) * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept", "*/*")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("Connection", "close")
	client.Header.Set("X-Forwarded-For", "127.0.0.1")
	client.Header.Set("Accept-Encoding", "gzip, deflate")
	client.Header.Set("Upgrade-Insecure-Requests", "1")
	client.Header.Set("Cookie", cookie)
	resp, err := client.R().Get(fullUrl)
	if err != nil {
		//fmt.Println(fmt.Sprintf("%s 访问异常 %s", fullUrl, err.Error()))
		return resp, err
	}
	return resp, err
}

func DoGetNoRedirect(fullUrl string, cookie string) (*resty.Response, error) {
	client := resty.New()
	client.
		SetRetryCount(0).
		SetRetryWaitTime(time.Duration(100) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(2000) * time.Millisecond).
		SetTimeout(time.Duration(3) * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept", "*/*")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("Connection", "close")
	client.Header.Set("X-Forwarded-For", "127.0.0.1")
	client.Header.Set("Accept-Encoding", "gzip, deflate")
	client.Header.Set("Upgrade-Insecure-Requests", "1")
	client.Header.Set("Cookie", cookie)
	client.SetRedirectPolicy(resty.DomainCheckRedirectPolicy(""))
	resp, err := client.R().Get(fullUrl)
	if err != nil {
		//fmt.Println(fmt.Sprintf("%s 访问异常", fullUrl))
		return resp, err
	}
	return resp, err
}

func DoPOST(fullUrl string, data map[string]string, cookie string) (*resty.Response, error) {
	client := resty.New()
	client.
		SetRetryCount(0).
		SetRetryWaitTime(time.Duration(100) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(2000) * time.Millisecond).
		SetTimeout(time.Duration(3) * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept", "*/*")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("Connection", "close")
	client.Header.Set("X-Forwarded-For", "127.0.0.1")
	client.Header.Set("Accept-Encoding", "gzip, deflate")
	client.Header.Set("Upgrade-Insecure-Requests", "1")
	client.Header.Set("Cookie", cookie)
	client.SetRedirectPolicy(resty.DomainCheckRedirectPolicy(""))
	resp, err := client.R().SetFormData(data).Post(fullUrl)
	if err != nil {
		//fmt.Println(fmt.Sprintf("%s 访问异常", fullUrl))
		return resp, err
	}
	return resp, err
}

func DoPOSTFile(fullUrl string, path string) (*resty.Response, error) {
	client := resty.New()
	client.
		SetRetryCount(0).
		SetRetryWaitTime(time.Duration(100) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(2000) * time.Millisecond).
		SetTimeout(time.Duration(3) * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept", "*/*")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("Connection", "close")
	client.Header.Set("X-Forwarded-For", "127.0.0.1")
	client.Header.Set("Accept-Encoding", "gzip, deflate")
	client.Header.Set("Upgrade-Insecure-Requests", "1")
	resp, err := client.R().SetFile("file", path).Post(fullUrl)
	if err != nil {
		//fmt.Println(fmt.Sprintf("%s 访问异常", fullUrl))
		return resp, err
	}
	return resp, err
}
