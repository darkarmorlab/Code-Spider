package core

import (
	"code-spider/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func YuqueOpen(q string, p int, cookie string) []string {
	var download_list []string
	for i := 1; i < p+1; i++ {
		zsearchurl := fmt.Sprintf("https://www.yuque.com/api/zsearch?q=%s&type=content&scope=/&tab=public&p=%s&sence=searchPage&time_horizon=", strings.Replace(q, " ", "+", -1), strconv.Itoa(p))
		resp1, err := util.DoGet(zsearchurl, cookie)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s 请求失败", zsearchurl))
		}
		// 使用 map[string]interface{} 解析 JSON
		var respjson1 map[string]interface{}
		err = json.Unmarshal(resp1.Body(), &respjson1)
		if err != nil {
			fmt.Println("解析 JSON 时出错:", err)
		}

		// 访问解析后的数据
		data := respjson1["data"].(map[string]interface{})
		hits := data["hits"].([]interface{})
		for i := 0; i < len(hits); i++ {
			hit := hits[i].(map[string]interface{})
			hiturl := hit["url"].(string)
			book_name := hit["book_name"].(string)
			record := hit["_record"].(map[string]interface{})
			book_id := fmt.Sprintf("%.f", record["book_id"].(float64))
			//fmt.Println(hiturl, book_id)
			markdown_download_url := fmt.Sprintf("https://www.yuque.com%s/markdown?attachment=true&latexcode=true&anchor=false&linebreak=false&&book_name=%s", hiturl, book_name)
			//fmt.Println(markdown_download_url)
			download_list = append(download_list, markdown_download_url)

			docsurl := fmt.Sprintf("https://www.yuque.com/api/docs?book_id=%s", book_id)
			resp2, err := util.DoGet(docsurl, cookie)
			if err != nil {
				fmt.Println(fmt.Sprintf("%s 请求失败", docsurl))
			}
			//fmt.Println(string(resp.Body()))

			// 使用 map[string]interface{} 解析 JSON
			var respjson2 map[string]interface{}
			err = json.Unmarshal(resp2.Body(), &respjson2)
			if err != nil {
				fmt.Println("解析 JSON 时出错:", err)
			}
			// 访问解析后的数据
			datas := respjson2["data"].([]interface{})
			for i := 0; i < len(datas); i++ {
				data_ := datas[i].(map[string]interface{})
				slug := data_["slug"].(string)
				hiturlsplit := strings.Split(hiturl, "/")
				hiturl_ := strings.Join(hiturlsplit[:len(hiturlsplit)-1], "/")
				book_name1 := book_name + fmt.Sprintf("-%s", data_["title"])
				//fmt.Println(hiturl_, slug)
				markdown_download_url = fmt.Sprintf("https://www.yuque.com%s/%s/markdown?attachment=true&latexcode=true&anchor=false&linebreak=false&&book_name=%s", hiturl_, slug, book_name1)
				//fmt.Println(markdown_download_url)
				download_list = append(download_list, markdown_download_url)
			}
		}
	}
	return download_list
}
