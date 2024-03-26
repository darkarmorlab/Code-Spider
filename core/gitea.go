package core

import (
	"bytes"
	"code-spider/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

func Gitea(target string, cookie string, gb bool) []string {
	download_list := []string{}
	//var username string
	//var password string
	//if len(cookie) == 0 {
	//	username, password, cookie = checkGitea(target)
	//}
	url_api := target + "/explore/repos"
	resp, _ := util.DoGet(url_api, cookie)
	projectDoc, _ := html.Parse(bytes.NewReader(resp.Body()))
	pageRedirect, _ := htmlquery.Query(projectDoc, "//a[@class=' item navigation']")
	pagehref := htmlquery.SelectAttr(pageRedirect, "href")
	pagenum, _ := strconv.Atoi(strings.Split(strings.Split(pagehref, "&")[0], "=")[1])
	fmt.Println(fmt.Sprintf("[!] 获取Gitea页码, 共%s页", strconv.Itoa(pagenum)))
	for i := 1; i < pagenum+1; i++ {
		pageurl := target + fmt.Sprintf("/explore/repos?page=%s&sort=recentupdate&q=&topic=false&language=&only_show_relevant=false", strconv.Itoa(i))
		resp, _ = util.DoGet(pageurl, cookie)
		projectDoc, _ = html.Parse(bytes.NewReader(resp.Body()))
		projectNode, _ := htmlquery.QueryAll(projectDoc, "//a[@class='name']")
		for _, urlNode := range projectNode {
			res := htmlquery.SelectAttr(urlNode, "href")
			if !gb {
				branch_api := fmt.Sprintf("%s%s/branches", target, res)
				resp, _ = util.DoGet(branch_api, cookie)
				projectDoc, _ = html.Parse(bytes.NewReader(resp.Body()))
				projectNode, _ = htmlquery.QueryAll(projectDoc, "//a[@class='item archive-link']")
				for _, branch := range projectNode {
					zip := htmlquery.SelectAttr(branch, "href")
					download_url := fmt.Sprintf("%s%s", target, zip)
					if strings.Contains(download_url, ".zip") {
						fmt.Println(fmt.Sprintf("[!] 获取zip: %s", download_url))
						download_list = append(download_list, download_url)
					}
				}
			} else {
				download_url1 := fmt.Sprintf("%s%s/archive/main.zip", target, res)
				download_list = append(download_list, download_url1)
				download_url2 := fmt.Sprintf("%s%s/archive/master.zip", target, res)
				download_list = append(download_list, download_url2)
			}
		}
	}

	//if len(password) != 0 {
	//	GiteaAccountDelete(target, username, password, cookie)
	//}

	return download_list

}

func checkGitea(target string) (username string, password string, cookie string) {
	url := target + "/user/sign_up"
	resp, err := util.DoGet(url, "")
	if err != nil {
		return
	}
	if strings.Contains(string(resp.Body()), "<label for=\"user_name\">") {
		fmt.Println(fmt.Sprintf("[!] %s 注册接口存在", url))
		username = util.RandSeq(6)
		email := fmt.Sprintf("%s@email.com", username)
		password = util.RandSeq(6)
		data := map[string]string{
			"user_name": username,
			"email":     email,
			"password":  password,
			"retype":    password,
		}
		resp, err = util.DoPOST(url, data, "")
		if resp.StatusCode() == 303 {
			fmt.Println(fmt.Sprintf("[+] %s %s %s %s 自动注册成功", url, username, email, password))
			cookie = strings.Split(resp.Header().Get("Set-Cookie"), ";")[0]
			return username, password, cookie
		} else {
			fmt.Println(fmt.Sprintf("[-] %s 需要邮箱激活，自动注册失败", url))
			return "", "", ""
		}
	}
	return "", "", ""
}

func GiteaAccountDelete(target string, username string, password string, cookie string) {
	url := target + "/user/settings/account/delete"
	data := map[string]string{
		"_autofill_dummy_username": "",
		"_autofill_dummy_password": "",
		"password":                 password,
	}
	resp, _ := util.DoPOST(url, data, cookie)
	if resp.StatusCode() == 303 {
		fmt.Println(fmt.Sprintf("[+] %s %s %s 自动删除成功", url, username, password))
	}
}
