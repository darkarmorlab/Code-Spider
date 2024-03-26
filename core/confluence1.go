package core

import (
	"bytes"
	"code-spider/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/FFFFFFStatic/resty"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func ConfluenceLogin(client resty.Client, target string, user string, pwd string) (*resty.Response, error) {
	url_api := target + "/login.action"
	// 准备登录请求的表单数据
	loginData := map[string]string{
		"os_username": user,
		"os_password": pwd,
	}
	// 发送登录请求
	//client.SetRedirectPolicy(resty.DomainCheckRedirectPolicy(""))
	resp, loginErr := client.R().SetFormData(loginData).Post(url_api)

	if loginErr != nil {
		fmt.Println("登录请求失败:", loginErr)
		return resp, loginErr
	}
	//fmt.Println(string(resp.Body()))
	return resp, loginErr
}

func Confluence(target string, cookie string, user string, pwd string) []string {
	// 创建 Resty 客户端
	client := resty.New()

	download_list := []string{}
	if len(user) != 0 && len(pwd) != 0 {
		_, err := ConfluenceLogin(*client, target, user, pwd)
		if err != nil {
			return download_list
		}

		resp, err := client.R().Get(target + "/dashboard.aciton")
		if err != nil {
			fmt.Println("请求失败:", err)
			return download_list
		}
		if strings.Contains(string(resp.Body()), "/logout.action") {
			fmt.Println(fmt.Sprintf("[+] %s 登陆成功 %s/%s", target, user, pwd))
			download_list = ConfluenceAuthorized(*client, target)
		} else {
			fmt.Println(fmt.Sprintf("[-] %s 登陆失败", target, user, pwd))
			return download_list
		}
	} else {
		download_list = ConfluenceUnauthorized(target)
	}
	return download_list
}

func ConfluenceAuthorized(client resty.Client, target string) []string {
	url_api1 := target + "/rest/experimental/search?cql=type%20=%20space&expand=space.icon&_=1695003716277"
	resp1, err := client.R().Get(url_api1)
	if err != nil {
		return nil
	}
	download_list := []string{}

	var confluencemain ConfluenceMain
	if resp1.StatusCode() == 200 {
		//fmt.Println(fmt.Sprintf("[+] %s 存在未授权", target))
		err = json.Unmarshal(resp1.Body(), &confluencemain)
		confluencedata := confluencemain.Results
		if err != nil {
			return nil
		}
		for i := 0; i < len(confluencedata); i++ {

			var NodeUrllist []string
			var NodeTitlelist []string

			confluencespace := confluencedata[i].Space
			confluenceexpandable := confluencespace.Expandable
			homepageid := strings.Split(confluenceexpandable.Homepage, "/")[len(strings.Split(confluenceexpandable.Homepage, "/"))-1]

			fmt.Println(fmt.Sprintf("[!] 获取空间名称: %s url: %s id: %s", confluencedata[i].Title, confluencedata[i].Url, homepageid))

			NodeUrllist = append(NodeUrllist, "/pages/viewpage.action?pageId="+homepageid)
			NodeTitlelist = append(NodeTitlelist, confluencedata[i].Title)
			url_api2 := target + fmt.Sprintf("/plugins/pagetree/naturalchildren.action?decorator=none&excerpt=false&sort=position&reverse=false&disableLinks=false&expandCurrent=true&placement=sidebar&hasRoot=true&pageId=%s&treeId=0&startDepth=0&mobile=false&ancestors=%s&treePageId=%s&_=1695026583391", homepageid, homepageid, homepageid)
			resp2, err := client.R().Get(url_api2)
			if err != nil {
				return nil
			}
			ConfluenceDoc, err := html.Parse(bytes.NewReader(resp2.Body()))
			if err != nil {
				return nil
			}
			ConfluenceNode, err := htmlquery.QueryAll(ConfluenceDoc, "//span[@class='plugin_pagetree_children_span']/a") // <span class="plugin_pagetree_children_span"
			if err != nil {
				return nil
			}
			for _, urlNode := range ConfluenceNode {
				if strings.Contains(htmlquery.SelectAttr(urlNode, "href"), ".action?pageId=") {
					fmt.Println(fmt.Sprintf("[!] 获取页面: %s 标题: %s", htmlquery.SelectAttr(urlNode, "href"), htmlquery.OutputHTML(urlNode, false)))
					NodeUrllist = append(NodeUrllist, htmlquery.SelectAttr(urlNode, "href"))
					NodeTitlelist = append(NodeTitlelist, htmlquery.OutputHTML(urlNode, false))
					NodeUrllist_, NodeTitlelist_ := CheckurlNodeIndex1(client, target, homepageid, htmlquery.SelectAttr(urlNode, "href"))
					NodeUrllist = append(NodeUrllist, NodeUrllist_...)
					NodeTitlelist = append(NodeTitlelist, NodeTitlelist_...)
				}
			}

			NodeUrllist = util.RemoveDuplicate(NodeUrllist)
			NodeTitlelist = util.RemoveDuplicate(NodeTitlelist)

			for j := 0; j < len(NodeUrllist); j++ {
				//fmt.Println(NodeUrllist[j], NodeTitlelist[j])
				//https://47.112.169.165/spaces/flyingpdf/pdfpageexport.action?pageId=68198340
				//url := fmt.Sprintf("%s/spaces/flyingpdf/pdfpageexport.action?%s&title=%s", target, strings.Split(NodeUrllist[j], "?")[1], NodeTitlelist[j])
				//https://47.112.169.165/exportword?pageId=68198340
				NodeTitlelist[j] = strings.Replace(NodeTitlelist[j], "/", "-", -1)
				url := fmt.Sprintf("%s/exportword?%s&title=%s", target, strings.Split(NodeUrllist[j], "?")[1], NodeTitlelist[j])
				download_list = append(download_list, url)
			}
		}
	}
	return download_list
}

func CheckurlNodeIndex1(client resty.Client, target string, treePageId string, pageId string) ([]string, []string) {
	pageId = strings.Split(strings.Split(pageId, "?")[1], "=")[1]
	url_api := target + fmt.Sprintf("/plugins/pagetree/naturalchildren.action?decorator=none&excerpt=false&sort=position&reverse=false&disableLinks=false&expandCurrent=true&placement=sidebar&hasRoot=true&pageId=%s&treeId=0&startDepth=0&mobile=false&treePageId=%s&_=1695028277882", pageId, treePageId)
	//fmt.Println(url_api)
	resp, err := client.R().Get(url_api)
	if err != nil {
		return []string{}, []string{}
	}
	ConfluenceDoc, err := html.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		return []string{}, []string{}
	}
	ConfluenceNode, err := htmlquery.QueryAll(ConfluenceDoc, "//div[@class='plugin_pagetree_children_content']/span") // <div class="plugin_pagetree_children_content"> <span class="plugin_pagetree_children_span"
	if err != nil {
		return []string{}, []string{}
	}

	var NodeUrllist []string
	var NodeTitlelist []string
	for _, urlNode := range ConfluenceNode {
		//if strings.Contains(htmlquery.SelectAttr(urlNode, "href"), ".action?pageId=") {
		if strings.Contains(htmlquery.SelectAttr(urlNode, "id"), "childrenspan") {
			id := strings.Replace(htmlquery.SelectAttr(urlNode, "id"), "childrenspan", "", -1)
			if strings.Contains(id, "-") {
				id = strings.Split(id, "-")[0]
			}
			id = "/pages/viewpage.action?pageId=" + id
			title, _ := htmlquery.Query(urlNode, "//a")
			fmt.Println(fmt.Sprintf("[!] 获取页面: %s 标题: %s", id, htmlquery.OutputHTML(title, false)))
			NodeUrllist = append(NodeUrllist, id)
			NodeTitlelist = append(NodeTitlelist, htmlquery.OutputHTML(title, false))
			NodeUrllist_, NodeTitlelist_ := CheckurlNodeIndex1(client, target, treePageId, id)
			NodeUrllist = append(NodeUrllist, NodeUrllist_...)
			NodeTitlelist = append(NodeTitlelist, NodeTitlelist_...)
		}
	}
	return NodeUrllist, NodeTitlelist
}

func ConfluenceUnauthorized(target string) []string {
	url_api1 := target + "/rest/experimental/search?cql=type%20=%20space&expand=space.icon&_=1695003716277"
	resp1, err := util.DoGet(url_api1, "")
	if err != nil {
		return nil
	}
	download_list := []string{}

	var confluencemain ConfluenceMain
	if resp1.StatusCode() == 200 {
		fmt.Println(fmt.Sprintf("[+] %s 存在未授权", target))
		err = json.Unmarshal(resp1.Body(), &confluencemain)
		confluencedata := confluencemain.Results
		if err != nil {
			return nil
		}
		for i := 0; i < len(confluencedata); i++ {

			var NodeUrllist []string
			var NodeTitlelist []string

			confluencespace := confluencedata[i].Space
			confluenceexpandable := confluencespace.Expandable
			homepageid := strings.Split(confluenceexpandable.Homepage, "/")[len(strings.Split(confluenceexpandable.Homepage, "/"))-1]

			fmt.Println(fmt.Sprintf("[!] 获取空间名称: %s url: %s id: %s", confluencedata[i].Title, confluencedata[i].Url, homepageid))

			NodeUrllist = append(NodeUrllist, "/pages/viewpage.action?pageId="+homepageid)
			NodeTitlelist = append(NodeTitlelist, confluencedata[i].Title)
			url_api2 := target + fmt.Sprintf("/plugins/pagetree/naturalchildren.action?decorator=none&excerpt=false&sort=position&reverse=false&disableLinks=false&expandCurrent=true&placement=sidebar&hasRoot=true&pageId=%s&treeId=0&startDepth=0&mobile=false&ancestors=%s&treePageId=%s&_=1695026583391", homepageid, homepageid, homepageid)
			resp2, err := util.DoGet(url_api2, "")
			if err != nil {
				return nil
			}
			ConfluenceDoc, err := html.Parse(bytes.NewReader(resp2.Body()))
			if err != nil {
				return nil
			}
			ConfluenceNode, err := htmlquery.QueryAll(ConfluenceDoc, "//span[@class='plugin_pagetree_children_span']/a") // <span class="plugin_pagetree_children_span"
			if err != nil {
				return nil
			}
			for _, urlNode := range ConfluenceNode {
				if strings.Contains(htmlquery.SelectAttr(urlNode, "href"), ".action?pageId=") {
					fmt.Println(fmt.Sprintf("[!] 获取页面: %s 标题: %s", htmlquery.SelectAttr(urlNode, "href"), htmlquery.OutputHTML(urlNode, false)))
					NodeUrllist = append(NodeUrllist, htmlquery.SelectAttr(urlNode, "href"))
					NodeTitlelist = append(NodeTitlelist, htmlquery.OutputHTML(urlNode, false))
					NodeUrllist_, NodeTitlelist_ := CheckurlNodeIndex(target, homepageid, htmlquery.SelectAttr(urlNode, "href"), "")
					NodeUrllist = append(NodeUrllist, NodeUrllist_...)
					NodeTitlelist = append(NodeTitlelist, NodeTitlelist_...)
				}
			}

			NodeUrllist = util.RemoveDuplicate(NodeUrllist)
			NodeTitlelist = util.RemoveDuplicate(NodeTitlelist)

			for j := 0; j < len(NodeUrllist); j++ {
				//fmt.Println(NodeUrllist[j], NodeTitlelist[j])
				//https://47.112.169.165/spaces/flyingpdf/pdfpageexport.action?pageId=68198340
				//url := fmt.Sprintf("%s/spaces/flyingpdf/pdfpageexport.action?%s&title=%s", target, strings.Split(NodeUrllist[j], "?")[1], NodeTitlelist[j])
				//https://47.112.169.165/exportword?pageId=68198340
				url := fmt.Sprintf("%s/exportword?%s&title=%s", target, strings.Split(NodeUrllist[j], "?")[1], NodeTitlelist[j])
				download_list = append(download_list, url)
			}
		}
	} else {
		fmt.Println(fmt.Sprintf("[-] %s 不存在未授权", target))
	}
	return download_list
}

func ConfluenceDownload(target string, url string, folder string, user string, pwd string) string {
	// 创建 Resty 客户端
	client := resty.New()
	var filename string
	nt := time.Now().Format("2006-01-02 15:01:05")
	if len(user) != 0 && len(pwd) != 0 {
		_, err := ConfluenceLogin(*client, target, user, pwd)
		if err != nil {
			fmt.Println("登陆失败")
		}

		if strings.Contains(url, "&") && !strings.Contains(url, ".action?pageId=") && !strings.Contains(url, "exportword?pageId=") {
			filename = strings.Split(url, "&")[0]
			filename = strings.Replace(filename, "?r=", "", -1)
			filename = strings.Replace(filename, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
		} else if strings.Contains(url, "&") && strings.Contains(url, ".action?pageId=") {
			filename = strings.Split(strings.Split(url, "&")[1], "=")[1] + ".pdf"
		} else if strings.Contains(url, "&") && strings.Contains(url, "exportword?pageId=") {
			filename = strings.Replace(url, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
			filename = filename + "-" + strings.Split(strings.Split(url, "&")[1], "=")[1] + ".doc"
		} else if url != "" {
			filename = strings.Replace(url, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
		} else {
			return ""
		}
		var fpath string
		if strings.Contains(url, ".action?pageId=") || strings.Contains(url, "exportword?pageId=") {
			fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, filename)
			fpath = fmt.Sprintf("%s/%s", folder, filename)
		} else {
			fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, url)
			fpath = fmt.Sprintf("%s/%s", folder, filename)
		}
		url = strings.Split(url, "&")[0]
		resp, err := client.R().Get(url)
		if err != nil {
			fmt.Println("文件下载失败:", err)
		}

		// 检查响应状态码
		if resp.StatusCode() != 200 {
			fmt.Println("文件下载失败，状态码:", resp.StatusCode())
		}

		newFile, err := os.Create(fpath)
		if err != nil {
			fmt.Println(err.Error())
			return "process failed for" + filename
		}
		defer newFile.Close()

		_, err = newFile.Write(resp.Body())

		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		if strings.Contains(url, "&") && !strings.Contains(url, ".action?pageId=") && !strings.Contains(url, "exportword?pageId=") {
			filename = strings.Split(url, "&")[0]
			filename = strings.Replace(filename, "?r=", "", -1)
			filename = strings.Replace(filename, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
		} else if strings.Contains(url, "&") && strings.Contains(url, ".action?pageId=") {
			filename = strings.Split(strings.Split(url, "&")[1], "=")[1] + ".pdf"
		} else if strings.Contains(url, "&") && strings.Contains(url, "exportword?pageId=") {
			filename = strings.Replace(url, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
			filename = filename + "-" + strings.Split(strings.Split(url, "&")[1], "=")[1] + ".doc"
		} else if url != "" {
			filename = strings.Replace(url, "/", "-", -1)
			filename = filename[strings.Index(filename, ":")+3:]
			filename = strings.Replace(filename, ":", "-", -1)
			filename = strings.Replace(filename, "?", "-", -1)
		} else {
			return ""
		}
		var fpath string
		if strings.Contains(url, ".action?pageId=") || strings.Contains(url, "exportword?pageId=") {
			fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, filename)
			fpath = fmt.Sprintf("%s/%s", folder, filename)
		} else {
			fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, url)
			fpath = fmt.Sprintf("%s/%s", folder, filename)
		}
		newFile, err := os.Create(fpath)
		if err != nil {
			fmt.Println(err.Error())
			return "process failed for" + filename
		}
		defer newFile.Close()

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := http.Client{
			Timeout:   9000 * time.Second,
			Transport: tr,
		}
		url = strings.Split(url, "&")[0]
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			return "download failed for" + filename
		}
		defer resp.Body.Close()

		_, err = io.Copy(newFile, resp.Body)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return filename
}
