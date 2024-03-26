package main

import (
	"code-spider/core"
	"code-spider/util"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	repo   = flag.String("repo", "", "代码仓库或知识文档 . 目前支持: gitea, gogs, gitlab, gitblit, jenkins, confluence, yuque_open, shimo")
	tar    = flag.String("tar", "", "目标 . example: http://1.1.1.1")
	cookie = flag.String("cookie", "", "Cookie")
	file   = flag.String("file", "", "URL File")
	user   = flag.String("user", "", "confluence登陆账号")
	pwd    = flag.String("pwd", "", "confluence登陆密码")
	q      = flag.String("q", "", "keyword 需要搜索的语雀内容")
	p      = flag.Int("p", 1, "PageSize 需要爬取的语雀公开搜索页数, 默认第一页")
	gb     = flag.Bool("gb", false, "是否拉取分支代码")
)

func main() {
	flag.Parse()

	download_list := []string{}

	// 输出目录
	var folder string

	if len(*file) != 0 {
		syncfilelist := make(chan string)
		getdownloadlist := make(chan string)
		downloadlist := make(chan string)
		wg := sync.WaitGroup{}
		//download_list = core.SyncFile(*file)
		folder = *file
		folder = strings.Replace(folder, ".", "-", -1)
		folder = "result/" + folder + "-file"

		wg.Add(3)

		go func() {
			core.SyncFile(*file, syncfilelist, &wg)
		}()

		go func() {
			core.GetDownloadList(syncfilelist, getdownloadlist, &wg)
		}()

		go func() {
			for getdownload := range getdownloadlist {
				downloadlist <- getdownload
			}
			wg.Done()
		}()
		// 异步接收数据并合到要下载的文件里
		go func() {
			for download := range downloadlist {
				download_list = append(download_list, download)
			}
		}()

		// 等待所有 goroutine 完成
		wg.Wait()

		// 所有 goroutine 完成后，关闭通道
		close(downloadlist)
	} else {
		if !strings.Contains(*tar, "http") {
			fmt.Println("未输入url或未以http开头")
			os.Exit(0)
		}
		switch *repo {
		case "gitea":
			download_list = core.Gitea(*tar, *cookie, *gb)
		case "gogs":
			download_list = core.Gogs(*tar, *cookie)
		case "gitlab":
			download_list = core.Gitlab(*tar, *cookie, *gb)
		case "jenkins":
			download_list = core.Jenkins(*tar, *cookie)
		case "gitblit":
			download_list = core.Gitblit(*tar, *cookie, *gb)
		case "confluence":
			download_list = core.Confluence(*tar, *cookie, *user, *pwd)
		//case "confluence1":
		//	download_list = core.Confluence1(*tar, *cookie, *user, *pwd)
		case "yuque_open":
			download_list = core.YuqueOpen(*q, *p, *cookie)
		case "shimo":
			download_list = core.Shimo(*cookie)
		}

		folder = *tar

		if strings.Contains(folder, "http") {
			folder = strings.Trim(folder, "http://")
			folder = strings.Trim(folder, "https://")
			folder = strings.Trim(folder, "/")
		}
		if strings.Contains(folder, ":") {
			folder = strings.Split(folder, ":")[0]
		}
		if *repo == "yuque_open" {
			folder = "result/" + strings.Replace(*q, " ", "-", -1) + "-yuque-open"
		} else if *repo == "shimo" {
			folder = "result/" + "shimo"
		} else {
			folder = "result/" + folder + "-" + *repo
		}
	}

	download_list = util.RemoveDuplicate(download_list)

	os.MkdirAll(folder, os.ModePerm) //创建目录

	//每个goroutine处理一个文件的下载
	ch := make(chan string)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, file := range download_list {
		go func(file string) {
			switch *repo {
			case "confluence":
				ch <- core.ConfluenceDownload(*tar, file, folder, *user, *pwd)
			default:
				ch <- core.Download(file, folder, *cookie)
			}
		}(file)
	}
	// 等待每个文件下载的完成，并检查超时
	timeout := time.After(20000 * time.Second)
	for idx := 0; idx < len(download_list); idx++ {
		select {
		case res := <-ch:
			nt := time.Now().Format("2006-01-02 15:04:05")
			if res != "" {
				fmt.Printf("[+] [%s] Finish download %s\n", nt, res)
			}
		case <-timeout:
			fmt.Println("Timeout...")
			break
		}
	}
}
