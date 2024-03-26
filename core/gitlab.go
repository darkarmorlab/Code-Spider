package core

import (
	"bytes"
	"code-spider/util"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

func Gitlab(target string, cookie string, gb bool) []string {
	download_list := []string{}
	url_api := target + "/explore/projects"
	projectNode := GitlabProjectAPI1(url_api, cookie)
	// fmt.Println(projectNode, len(projectNode))
	var projecturl []string
	if len(projectNode) == 0 {
		projectNode = GitlabProjectAPI2(url_api, cookie)
		if len(projectNode) == 0 {
			projectNode = GitlabProjectAPI3(url_api, cookie)
		}
		for _, urlNode := range projectNode {
			result := htmlquery.SelectAttr(urlNode, "href")
			projecturl = append(projecturl, result)
		}
	} else {
		for _, urlNode := range projectNode {
			result := htmlquery.SelectAttr(urlNode, "href")
			projecturl = append(projecturl, result)
		}
	}
	projecturl = util.RemoveDuplicate(projecturl)
	// 拉取代码分支代码
	if !gb {
		for _, url := range projecturl {
			getdescription_url := fmt.Sprintf("%s%s", target, url)
			resp, err := util.DoGet(getdescription_url, cookie)
			if err != nil {
				fmt.Println(fmt.Sprintf("[-] %s 请求失败", getdescription_url))
			}
			projectDoc, err := html.Parse(bytes.NewReader((resp.Body())))
			if err != nil {
				return nil
			}
			projectNode, err = htmlquery.QueryAll(projectDoc, "//meta[@name='description']")
			if err != nil {
				return nil
			}
			var description string
			for _, node := range projectNode {
				description = htmlquery.SelectAttr(node, "content")
			}
			branch_api := fmt.Sprintf("%s%s/refs?ref=master&search=", target, url)
			resp, _ = util.DoGetNoRedirect(branch_api, cookie)
			// fmt.Println(branch_api, resp.StatusCode())
			if resp.StatusCode() == 200 {
				// 使用 map[string]interface{} 解析 JSON
				var respjson map[string]interface{}
				err = json.Unmarshal(resp.Body(), &respjson)
				if err != nil {
					fmt.Println("解析 JSON 时出错:", err)
				}
				Branches := respjson["Branches"].([]interface{})
				// http://x.x.x.x/jojo/dy/-/archive/master/dy-master.zip
				// fmt.Println(Branches[0])
				if len(Branches) != 0 {
					fileName := url[strings.LastIndex(url, "/")+1:] + fmt.Sprintf("-%s.zip", Branches[0])
					test_download_url := fmt.Sprintf("%s%s/-/archive/%s/%s", target, url, Branches[0], fileName)
					//fmt.Println(test_download_url)
					resp, _ = util.DoGetNoRedirect(test_download_url, cookie)
					// fmt.Println(resp.StatusCode())
					for _, branch := range Branches {
						fmt.Println(fmt.Sprintf("[+] %s%s 获取到代码分支 %s", target, url, branch))
						var download_url string
						if resp.StatusCode() == 200 {
							fileName = url[strings.LastIndex(url, "/")+1:] + fmt.Sprintf("-%s.zip", branch)
							download_url = fmt.Sprintf("%s%s/-/archive/%s/%s&description=%s", target, url, branch, fileName, description)
							//fmt.Println(download_url)
						} else {
							// http://x.x.x.x/chenPeng/TenementApp_Dabai/repository/hyc_dev/archive.zip
							download_url = fmt.Sprintf("%s%s/repository/%s/archive.zip&description=%s", target, url, branch, description)
							// fmt.Println(download_url)
						}
						if description == "GitLab Community Edition" || description == "GitLab Enterprise Edition" {
							download_url = strings.Split(download_url, "&")[0]
						}
						download_list = append(download_list, download_url)
					}
				}
			} else {
				getdescription_url = fmt.Sprintf("%s%s", target, url)
				resp, err = util.DoGet(getdescription_url, cookie)
				if err != nil {
					fmt.Println(fmt.Sprintf("[-] %s 请求失败", getdescription_url))
				}
				projectDoc, err = html.Parse(bytes.NewReader((resp.Body())))
				if err != nil {
					return nil
				}
				projectNode, err = htmlquery.QueryAll(projectDoc, "//meta[@name='description']")
				if err != nil {
					return nil
				}
				for _, node := range projectNode {
					description = htmlquery.SelectAttr(node, "content")
				}

				branch_api1 := fmt.Sprintf("%s%s/branches", target, url)
				resp, err = util.DoGet(branch_api1, cookie)
				if err != nil {
					return nil
				}
				projectDoc, err = html.Parse(bytes.NewReader((resp.Body())))
				if err != nil {
					return nil
				}
				projectNode, err = htmlquery.QueryAll(projectDoc, "//ul[@class='content-list all-branches']/li/div/a")
				if len(projectNode) == 0 {
					projectNode, err = htmlquery.QueryAll(projectDoc, "//ul[@class='bordered-list top-list all-branches']/li/h4/a")
				}
				if err != nil {
					return nil
				}
				var branch_urls []string
				for _, urlNode := range projectNode {
					result := htmlquery.SelectAttr(urlNode, "href")
					if strings.Contains(result, "tree") {
						branch_urls = append(branch_urls, result)
					}
				}
				branch_urls = util.RemoveDuplicate(branch_urls)
				for _, branch := range branch_urls {
					// http://x.x.x.x/project/node_take_out/repository/archive.zip?ref=master
					branch_ := strings.Split(branch, "/")
					branch_ = branch_[len(branch_)-1:]
					fmt.Println(fmt.Sprintf("[+] %s%s 获取到代码分支 %s", target, url, branch_[0]))
					download_url := fmt.Sprintf("%s%s/repository/archive.zip?ref=%s&description=%s", target, url, branch_[0], description)
					if description == "GitLab Community Edition" || description == "GitLab Enterprise Edition" {
						download_url = strings.Split(download_url, "&")[0]
					}
					download_list = append(download_list, download_url)
				}
			}
		}
	} else {
		// 不拉取代码分支代码
		for _, url := range projecturl {
			getdescription_url := fmt.Sprintf("%s%s", target, url)
			resp, err := util.DoGet(getdescription_url, cookie)
			if err != nil {
				fmt.Println(fmt.Sprintf("[-] %s 请求失败", getdescription_url))
			}
			projectDoc, err := html.Parse(bytes.NewReader((resp.Body())))
			if err != nil {
				return nil
			}
			projectNode, err = htmlquery.QueryAll(projectDoc, "//meta[@name='description']")
			if err != nil {
				return nil
			}
			var description string
			for _, node := range projectNode {
				description = htmlquery.SelectAttr(node, "content")
			}
			branch_api := fmt.Sprintf("%s%s/refs?ref=master&search=", target, url)
			resp, _ = util.DoGetNoRedirect(branch_api, cookie)
			// fmt.Println(branch_api, resp.StatusCode())
			if resp.StatusCode() == 200 {
				// http://x.x.x.x/jojo/dy/-/archive/master/dy-master.zip
				// fmt.Println(Branches[0])
				fileName := url[strings.LastIndex(url, "/")+1:] + "-master.zip"
				test_download_url := fmt.Sprintf("%s%s/-/archive/master/%s", target, url, fileName)
				//fmt.Println(test_download_url)
				resp, _ = util.DoGetNoRedirect(test_download_url, cookie)
				// fmt.Println(resp.StatusCode())
				var download_url string
				if resp.StatusCode() == 200 {
					fileName = url[strings.LastIndex(url, "/")+1:] + "-master.zip"
					download_url = fmt.Sprintf("%s%s/-/archive/master/%s&description=%s", target, url, fileName, description)
					//fmt.Println(download_url)
				} else {
					// http://x.x.x.x/chenPeng/TenementApp_Dabai/repository/hyc_dev/archive.zip
					download_url = fmt.Sprintf("%s%s/repository/master/archive.zip&description=%s", target, url, description)
					// fmt.Println(download_url)
				}
				if description == "GitLab Community Edition" || description == "GitLab Enterprise Edition" {
					download_url = strings.Split(download_url, "&")[0]
				}
				download_list = append(download_list, download_url)
			} else {
				download_url := fmt.Sprintf("%s%s/repository/archive.zip?ref=master&description=%s", target, url, description)
				if description == "GitLab Community Edition" || description == "GitLab Enterprise Edition" {
					download_url = strings.Split(download_url, "&")[0]
				}
				download_list = append(download_list, download_url)
			}
		}
	}
	/*
		for _, urlNode := range projectNode {
			result := htmlquery.SelectAttr(urlNode, "href")
			fileName := result[strings.LastIndex(result, "/")+1:] + "-master.zip"
			url := fmt.Sprintf("%s%s/-/archive/master/%s", target, result, fileName)
			download_list = append(download_list, url)
		}
	*/
	download_list = util.RemoveDuplicate(download_list)

	return download_list
}

func GitlabProjectAPI1(target string, cookie string) []*html.Node {
	resp, err := util.DoGet(target, cookie)
	if err != nil {
		return nil
	}
	projectDoc, err := html.Parse(bytes.NewReader((resp.Body())))
	if err != nil {
		return nil
	}
	projectNode, err := htmlquery.QueryAll(projectDoc, "//a[@class='text-plain']")
	if err != nil {
		return nil
	}
	return projectNode
}

func GitlabProjectAPI2(target string, cookie string) []*html.Node {
	resp, err := util.DoGet(target, cookie)
	if err != nil {
		return nil
	}
	projectDoc, err := html.Parse(bytes.NewReader((resp.Body())))
	if err != nil {
		return nil
	}
	projectNode, err := htmlquery.QueryAll(projectDoc, "//a[@class='project']")
	if err != nil {
		return nil
	}
	return projectNode
}

func GitlabProjectAPI3(target string, cookie string) []*html.Node {
	resp, err := util.DoGet(target, cookie)
	if err != nil {
		return nil
	}
	projectDoc, err := html.Parse(bytes.NewReader((resp.Body())))
	if err != nil {
		return nil
	}
	projectNode, err := htmlquery.QueryAll(projectDoc, "//ul[@class='bordered-list top-list']/li/h4/a")
	if err != nil {
		return nil
	}
	return projectNode
}
