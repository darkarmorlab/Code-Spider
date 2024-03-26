package core

import (
	"code-spider/util"
	"fmt"
	"strings"
	"sync"
)

/*
func SyncFile(file string) []string {
	UrlList, err := util.Readfile(file)
	if err != nil {
		return nil
	}
	download_list := []string{}
	for _, url := range UrlList {
		if !strings.Contains(url, "http") {
			url = "http://" + url
		}
		resp, err := util.DoGet(url, "")
		if err != nil {
			return nil
		}
		repo := Finger(string(resp.Body()), util.GetRespHeader(resp.Header()))
		switch repo {
		case "gitea":
			fmt.Println(fmt.Sprintf("%s 识别为 gitea", url))
			download_list = append(download_list, Gitea(url, "")...)
		case "gogs":
			fmt.Println(fmt.Sprintf("%s 识别为 gogs", url))
			download_list = append(download_list, Gogs(url, "")...)
		case "gitlab":
			fmt.Println(fmt.Sprintf("%s 识别为 gitlab", url))
			download_list = append(download_list, Gitlab(url, "")...)
		case "jenkins":
			fmt.Println(fmt.Sprintf("%s 识别为 jenkins", url))
			download_list = append(download_list, Jenkins(url, "")...)
		case "gitblit":
			fmt.Println(fmt.Sprintf("%s 识别为 gitblit", url))
			download_list = append(download_list, Gitblit(url, "")...)
		}
	}
	return download_list
}
*/

func SyncFile(file string, syncfilelist chan string, wg *sync.WaitGroup) {
	defer close(syncfilelist)
	UrlList, err := util.Readfile(file)
	if err != nil {
		return
	}
	for _, url := range UrlList {
		if !strings.Contains(url, "http") {
			url = "http://" + url
		}
		syncfilelist <- url
	}

	wg.Done()
}

func GetDownloadList(syncfilelist chan string, getdownloadlist chan string, wg *sync.WaitGroup) {
	for url := range syncfilelist {
		var download_list []string
		resp, err := util.DoGet(url, "")
		if err != nil {
			return
		}
		repo := Finger(string(resp.Body()), util.GetRespHeader(resp.Header()))
		switch repo {
		case "gitea":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 gitea", url))
			download_list = append(download_list, Gitea(url, "", false)...)
		case "gogs":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 gogs", url))
			download_list = append(download_list, Gogs(url, "")...)
		case "gitlab":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 gitlab", url))
			download_list = append(download_list, Gitlab(url, "", false)...)
		case "jenkins":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 jenkins", url))
			download_list = append(download_list, Jenkins(url, "")...)
		case "gitblit":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 gitblit", url))
			download_list = append(download_list, Gitblit(url, "", false)...)
		case "confluence":
			fmt.Println(fmt.Sprintf("[!] %s 识别为 confluence", url))
			download_list = append(download_list, Confluence(url, "", "", "")...)
		}
		if len(download_list) != 0 {
			for _, download := range download_list {
				getdownloadlist <- download
			}
		} else {
			getdownloadlist <- ""
		}
	}

	close(getdownloadlist)
	wg.Done()
}

func Finger(body string, header string) string {
	if IsGitea(body, header) != "" {
		return "gitea"
	} else if IsGogs(body, header) != "" {
		return "gogs"
	} else if IsGitblit(body, header) != "" {
		return "gitblit"
	} else if IsGitLab(body, header) != "" {
		return "gitlab"
	} else if IsJenkins(body, header) != "" {
		return "jenkins"
	} else if IsConfluence(body, header) != "" {
		return "confluence"
	}
	return ""
}

func IsGitea(body string, header string) string {
	colomn1 := strings.Contains(body, "href=\"https://docs.gitea.io\">Help</a>")
	column2 := strings.Contains(header, "Set-Cookie: i_like_gitea=")
	if colomn1 || column2 {
		return "gitea"
	} else {
		return ""
	}
}

func IsGogs(body string, header string) string {
	column1 := strings.Contains(body, "<a class=\"item\" target=\"_blank\" rel=\"noopener noreferrer\" href=\"https://gogs.io/docs\" rel=\"noreferrer\">帮助</a>")
	column2 := strings.Contains(body, "content=\"Gogs\"")
	if column1 || column2 {
		return "gogs"
	} else {
		return ""
	}
}

func IsGitLab(body string, header string) string {
	column1 := strings.Contains(body, "class=\"col-sm-7 brand-holder pull-left\"") || strings.Contains(body, "<a href=\"https://about.gitlab.com/\">About GitLab</a>")
	column2 := strings.Contains(body, "content=\"GitLab\"")
	column3 := strings.Contains(body, "sign in · gitlab")
	column4 := strings.Contains(body, "content=\"Gitlab Community Edition\"")
	column5 := strings.Contains(body, "gon.default_issues_tracker")
	column6 := strings.Contains(header, "_gitlab_session")
	if ((((column1 || column2) || column3) || column4) || column5) || column6 {
		return "gitlab"
	} else {
		return ""
	}
}

func IsGitblit(body string, header string) string {
	column1 := strings.Contains(body, "<a title=\"gitblit homepage\" href=\"http://gitblit.com/\">")
	if column1 {
		return "gitblit"
	} else {
		return ""
	}
}

func IsJenkins(body string, header string) string {
	column1 := strings.Contains(header, "Content-Type: text/plain") && strings.Contains(body, "jenkins-agent-protocols")
	column2 := strings.Contains(header, "X-Required-Permission: hudson.model.hudson.read")
	column3 := strings.Contains(header, "X-Hudson")
	column4 := strings.Contains(header, "X-Jenkins")
	if column1 || column2 || column3 || column4 {
		return "jenkins"
	} else {
		return ""
	}
}

func IsConfluence(body string, header string) string {
	column1 := strings.Contains(body, "Atlassian Confluence")
	column2 := strings.Contains(body, "id=\"com-atlassian-confluence\"") && strings.Contains(body, "name=\"confluence-base-url\"")
	column3 := !strings.Contains(header, "tp-link router upnp") && strings.Contains(header, "X-Confluence-")
	if (column1 || column2) || column3 {
		return "confluence"
	} else {
		return ""
	}
}
