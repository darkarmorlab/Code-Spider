package core

import (
	"code-spider/util"
	"fmt"
	"github.com/tidwall/gjson"
)

func Jenkins(target string, cookie string) []string {
	download_list := []string{}
	url_api := target + "/api/json?pretty=true"
	resp, err := util.DoGet(url_api, cookie)
	if err != nil {
		return nil
	}
	json := gjson.Parse(resp.String())
	result := json.Get("jobs")
	for _, tar := range result.Array() {
		url := fmt.Sprintf("%sws/*zip*/%s.zip", tar.Get("url"), tar.Get("name"))
		download_list = append(download_list, url)
	}

	return download_list
}
