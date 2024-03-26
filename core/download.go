package core

import (
	"crypto/tls"
	"fmt"
	"github.com/FFFFFFStatic/resty"
	"os"
	"strings"
	"time"
)

//单个文件下载
func Download(url string, folder string, cookie string) string {
	nt := time.Now().Format("2006-01-02 15:01:05")
	var filename string
	if strings.Contains(url, "&") && !strings.Contains(url, ".action?pageId=") && !strings.Contains(url, "exportword?pageId=") && !strings.Contains(url, "description=") { // 处理url中含&的
		if strings.Contains(url, "www.yuque.com") {
			filename = "yuque-open-search-" + strings.Split(strings.Split(url, "&&")[1], "=")[1] + ".md"
		} else if strings.Contains(url, "shimo.im") {
			filename = "shimo-" + strings.Split(strings.Split(url, "&&")[1], "=")[1] + ".md"
		} else {
			// http://110.241.221.59:2001/zip/?r=hlw-ui.git&format=zip
			// filename = strings.Join(strings.Split(url, "/")[4:], "-")
			// filename = strings.Split(filename, "&")[0]
			// filename = strings.Replace(filename, "?r=", "", -1)
			// filename = strings.Replace(filename, "%2F", "-", -1)
			// filename = strings.Replace(filename, ".git", "", -1)
			// filename = strings.Replace(filename, "?", "-", -1) + ".zip"
			// http://110.241.221.59:2001/zip/?r=ziran%2Fziran-web3.git&format=zip&&name=3.4.4版本架构-单体  自然资源管理-54.zip
			filename = strings.Split(strings.Split(url, "&&")[1], "=")[1]
			filename = strings.Replace(filename, "%2F", "-", -1)
		}
	} else if strings.Contains(url, "&") && strings.Contains(url, ".action?pageId=") { // confluence
		filename = strings.Split(strings.Split(url, "&")[1], "=")[1] + ".pdf"
	} else if strings.Contains(url, "&") && strings.Contains(url, "exportword?pageId=") { // confluence
		filename = strings.Replace(url, "/", "-", -1)
		filename = filename[strings.Index(filename, ":")+3:]
		filename = strings.Replace(filename, ":", "-", -1)
		filename = strings.Replace(filename, "?", "-", -1)
		filename = filename + "-" + strings.Split(strings.Split(url, "&")[1], "=")[1] + ".doc"
	} else if strings.Contains(url, "&") && strings.Contains(url, "description=") { // gitlab
		filename = strings.Split(strings.Split(url, "&")[1], "=")[1] + ".zip"
		//url_s := strings.Split(url, "/")
		//filename = strings.Join(url_s[3:], "-")
	} else if url != "" {
		// 去除下载文件名中的域名或ip信息
		filename = strings.Join(strings.Split(url, "/")[3:], "-")
		//filename = strings.Replace(url, "/", "-", -1)
		//filename = filename[strings.Index(filename, ":")+3:]
		filename = strings.Replace(filename, ":", "-", -1)
		filename = strings.Replace(filename, "?", "-", -1)
	} else {
		return ""
	}
	var fpath string
	//fmt.Println(url, filename)
	if strings.Contains(url, ".action?pageId=") || strings.Contains(url, "exportword?pageId=") || strings.Contains(url, "description=") { // confluence && gitlab
		fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, filename)
		fpath = fmt.Sprintf("%s/%s", folder, filename)
	} else {
		fmt.Printf("[!] [%s] DownloadFIle: %s\n", nt, url)
		fpath = fmt.Sprintf("%s/%s", folder, filename)
	}
	newFile, err := os.Create(fpath)
	if err != nil {
		fmt.Println(err.Error())
		return "[-] process failed for " + filename
	}
	defer newFile.Close()

	/*
		// 创建自定义的http.Transport以跳过TLS验证
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

			// 创建一个HTTP客户端，指定自定义Transport和超时时间
			client := &http.Client{
				Timeout:   9000 * time.Second,
				Transport: tr,
			}

			url = strings.Split(url, "&")[0]
			// 创建一个GET请求
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("创建请求时出错:", err)
			}

			// 添加包含多个Cookie的字符串到请求头
			req.Header.Add("Cookie", cookie)

			// 发送请求
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(err.Error())
				return "download failed for" + filename
			}
			defer resp.Body.Close()

			_, err = io.Copy(newFile, resp.Body)

			if err != nil {
				fmt.Println(err.Error())
			}

	*/
	// 创建 Resty 客户端
	client := resty.New()
	client.
		SetRetryCount(1).
		SetRetryWaitTime(time.Duration(500) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(5000) * time.Millisecond).
		SetTimeout(20000 * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	client.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	client.Header.Set("X-Forwarded-For", "127.0.0.1")
	client.Header.Set("Accept-Encoding", "gzip, deflate")
	client.Header.Set("Cookie", cookie)

	if strings.Contains(url, "&&") {
		url = strings.Split(url, "&&")[0]
	} else {
		url = strings.Split(url, "&")[0]
	}

	resp, err := client.R().Get(url)
	if err != nil {
		fmt.Println("[-] "+url+" 文件下载失败:", err)
		fmt.Println("[!] " + url + " 尝试重新下载.tar.gz")
		if strings.Contains(url, ".zip") {
			url = strings.Replace(url, ".zip", ".tar.gz", -1)
			resp, _ = client.R().Get(url)

			if err != nil {
				fmt.Println("[-] "+url+" 文件下载失败:", err)
			}
			if resp.StatusCode() != 200 {
				fmt.Println("[-] "+url+" 文件下载失败，状态码:", resp.StatusCode())
			}
		}
	}

	if resp.StatusCode() != 200 {
		fmt.Println("[-] "+url+" 文件下载失败，状态码:", resp.StatusCode())
		fmt.Println("[!] " + url + " 尝试重新下载.tar.gz")
		if strings.Contains(url, ".zip") {
			url = strings.Replace(url, ".zip", ".tar.gz", -1)
			resp, _ = client.R().Get(url)

			if err != nil {
				fmt.Println("[-] "+url+" 文件下载失败:", err)
			}
			if resp.StatusCode() != 200 {
				fmt.Println("[-] "+url+" 文件下载失败，状态码:", resp.StatusCode())
			}
		}
	}
	_, err = newFile.Write(resp.Body())

	if err != nil {
		fmt.Println("[-] "+filename+" 写入文件失败:", err)
	}

	return filename
}
