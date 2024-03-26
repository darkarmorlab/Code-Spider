package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/FFFFFFStatic/resty"
	"strings"
	"time"
)

func Shimo(cookie string) []string {
	var download_list []string
	client := resty.New()
	client.
		SetRetryCount(0).
		SetRetryWaitTime(time.Duration(100) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(2000) * time.Millisecond).
		SetTimeout(time.Duration(3) * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept", "application/nd.shimo.v2+json")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("Cookie", cookie)
	client.Header.Set("Referer", "https://shimo.im/")
	ShimoEnterpriceSearch(client)
	download_list = append(download_list, ShimoTeamSpaceSearch(client)...)
	download_list = append(download_list, ShimoPersonSpaceSearch(client)...)
	return download_list
}

// 企业信息
func ShimoEnterpriceSearch(client *resty.Client) {
	url := "https://shimo.im/lizard-api/org/departments/1/users?perPage=100&page=1"
	resp, err := client.R().Get(url)
	if err != nil {
		fmt.Println("请求失败", err.Error())
	}
	// 使用 map[string]interface{} 解析 JSON
	var respjson map[string]interface{}
	err = json.Unmarshal(resp.Body(), &respjson)
	if err != nil {
		fmt.Println("解析 JSON 时出错:", err)
	}
	total := respjson["total"].(float64)
	fmt.Println(fmt.Sprintf("共获取到企业人员总数（默认跑企业id为1): %.f", total))
	users := respjson["users"].([]interface{})
	for i := 0; i < len(users); i++ {
		user := users[i].(map[string]interface{})
		name := user["name"].(string)
		email := user["email"].(string)
		fmt.Println(fmt.Sprintf("获取到成员 %s  email: %s", name, email))
	}
}

// 团队空间
func ShimoTeamSpaceSearch(client *resty.Client) []string {
	var download_list []string

	url1 := "https://shimo.im/panda-api/file/spaces?orderBy=updatedAt"
	resp1, err := client.R().Get(url1)
	if err != nil {
		fmt.Println("请求失败", err.Error())
	}
	// 使用 map[string]interface{} 解析 JSON
	var respjson1 map[string]interface{}
	err = json.Unmarshal(resp1.Body(), &respjson1)
	if err != nil {
		fmt.Println("解析 JSON 时出错:", err)
	}
	spaces := respjson1["spaces"].([]interface{})
	for i := 0; i < len(spaces); i++ {
		space := spaces[i].(map[string]interface{})
		spacename := space["name"].(string)
		guid := space["guid"].(string)
		//fmt.Println(guid)
		url2 := fmt.Sprintf("https://shimo.im/lizard-api/files?folder=%s", guid)
		resp2, err := client.R().Get(url2)
		if err != nil {
			fmt.Println("请求失败", err.Error())
		}
		// 使用 map[string]interface{} 解析 JSON
		var files []interface{}
		err = json.Unmarshal(resp2.Body(), &files)
		if err != nil {
			fmt.Println("解析 JSON 时出错:", err)
		}
		for j := 0; j < len(files); j++ {
			file := files[j].(map[string]interface{})
			filename := file["name"].(string)
			fileurl := file["url"].(string)
			fileurl = strings.Join(strings.Split(fileurl, "/")[len(strings.Split(fileurl, "/"))-1:], "")
			url3 := fmt.Sprintf("https://shimo.im/lizard-api/office-gw/files/export?fileGuid=%s&type=md", fileurl)
			resp3, err := client.R().Get(url3)
			if err != nil {
				fmt.Println("请求失败", err.Error())
			}

			var respjson3 map[string]interface{}
			err = json.Unmarshal(resp3.Body(), &respjson3)
			if err != nil {
				fmt.Println("解析 JSON 时出错:", err)
			}
			taskId := respjson3["taskId"].(string)
			url4 := fmt.Sprintf("https://shimo.im/lizard-api/office-gw/files/export/progress?taskId=%s", taskId)
			resp4, err := client.R().Get(url4)
			if err != nil {
				fmt.Println("请求失败", err.Error())
			}
			var downloadjson map[string]interface{}
			err = json.Unmarshal(resp4.Body(), &downloadjson)
			if err != nil {
				fmt.Println("解析 JSON 时出错:", err)
			}
			data := downloadjson["data"].(map[string]interface{})
			downloadUrl := data["downloadUrl"].(string)
			downloadname := strings.Replace(spacename, "/", "-", -1) + "-" + strings.Replace(filename, "/", "-", -1)
			downloadUrl = fmt.Sprintf("%s&&downloadname=%s", downloadUrl, downloadname)
			download_list = append(download_list, downloadUrl)
		}
	}
	return download_list
}

func ShimoPersonSpaceSearch(client *resty.Client) []string {
	var download_list []string
	url1 := "https://shimo.im/lizard-api/files"
	resp1, err := client.R().Get(url1)
	if err != nil {
		fmt.Println("请求失败", err.Error())
	}
	var files []interface{}
	err = json.Unmarshal(resp1.Body(), &files)
	for i := 0; i < len(files); i++ {
		file := files[i].(map[string]interface{})
		filename := file["name"].(string)
		fileurl := file["url"].(string)
		if !strings.Contains(fileurl, "/shortcut/") {
			fileurl = strings.Join(strings.Split(fileurl, "/")[len(strings.Split(fileurl, "/"))-1:], "")
			url2 := fmt.Sprintf("https://shimo.im/lizard-api/office-gw/files/export?fileGuid=%s&type=md", fileurl)
			resp2, err := client.R().Get(url2)
			if err != nil {
				fmt.Println("请求失败", err.Error())
			}

			var respjson2 map[string]interface{}
			err = json.Unmarshal(resp2.Body(), &respjson2)
			if err != nil {
				fmt.Println("解析 JSON 时出错:", err)
			}
			taskId := respjson2["taskId"].(string)
			url3 := fmt.Sprintf("https://shimo.im/lizard-api/office-gw/files/export/progress?taskId=%s", taskId)
			resp3, err := client.R().Get(url3)
			if err != nil {
				fmt.Println("请求失败", err.Error())
			}
			var downloadjson map[string]interface{}
			err = json.Unmarshal(resp3.Body(), &downloadjson)
			if err != nil {
				fmt.Println("解析 JSON 时出错:", err)
			}
			data := downloadjson["data"].(map[string]interface{})
			downloadUrl := data["downloadUrl"].(string)
			downloadname := "个人空间-" + filename
			downloadUrl = fmt.Sprintf("%s&&downloadname=%s", downloadUrl, downloadname)
			download_list = append(download_list, downloadUrl)
		}
	}
	return download_list
}
