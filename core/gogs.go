package core

import (
	"bytes"
	"code-spider/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

func Gogs(target string, cookie string) []string {
	download_list := []string{}
	if len(cookie) == 0 {
		checkGogs(target)
	}
	url_api := target + "/explore/repos"
	resp, err := util.DoGet(url_api, cookie)
	if err != nil {
		return nil
	}
	projectDoc, err := html.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil
	}
	projectNode, err := htmlquery.QueryAll(projectDoc, "//*[@class=\"ui header\"]/a")
	if err != nil {
		return nil
	}
	for _, urlNode := range projectNode {
		res := htmlquery.SelectAttr(urlNode, "href")
		url := fmt.Sprintf("%s%s/archive/master.zip", target, res)
		download_list = append(download_list, url)
	}

	return download_list

}

func checkGogs(target string) {
	url := target + "/user/sign_up"
	resp, err := util.DoGet(url, "")
	if err != nil {
		return
	}
	if strings.Contains(string(resp.Body()), "<label for=\"user_name\">") {
		fmt.Println(fmt.Sprintf("[!] %s 注册接口存在", url))
	}
}
