# Code-Spider

## Help

```
 ____                __                   ____                    __                  
/\  _`\             /\ \                 /\  _`\           __    /\ \                 
\ \ \/\_\    ___    \_\ \     __         \ \,\L\_\  _____ /\_\   \_\ \     __   _ __  
 \ \ \/_/_  / __`\  /'_` \  /'__`\ _______\/_\__ \ /\ '__`\/\ \  /'_` \  /'__`\/\`'__\
  \ \ \L\ \/\ \L\ \/\ \L\ \/\  __//\______\ /\ \L\ \ \ \L\ \ \ \/\ \L\ \/\  __/\ \ \/ 
   \ \____/\ \____/\ \___,_\ \____\/______/ \ `\____\ \ ,__/\ \_\ \___,_\ \____\\ \_\ 
    \/___/  \/___/  \/__,_ /\/____/          \/_____/\ \ \/  \/_/\/__,_ /\/____/ \/_/ 
                                                      \ \_\                           
                                                       \/_/                           
Usage of ./code_spider_darwin_amd64:
  -cookie string
    	Cookie
  -file string
    	URL File
  -gb
    	是否拉取gitlab分支代码
  -p int
    	PageSize 需要爬取的语雀公开搜索页数, 默认第一页 (default 1)
  -pwd string
    	confluence登陆密码
  -q string
    	keyword 需要搜索的语雀内容
  -repo string
    	代码仓库或知识文档 . 目前支持: gitea, gogs, gitlab, gitblit, jenkins, confluence, yuque_open, shimo
  -tar string
    	目标 . example: http://1.1.1.1
  -user string
    	confluence登陆账号
```

## 使用须知
```
Gogs
    /explore/repos 存在鉴权，不存在未授权
    /user/sign_up 注册接口，<label for="user_name">判断是否存在
     鉴权需要cookie i_like_gogs=
     输入cookie时直接下源码，未输入cookie时判断注册接口是否存在
     注册存在验证码无法绕过
    
Gitea
    /explore/repos 存在未授权，无需cookie也可以跑
    鉴权需要cookie i_like_gitea=
    未输入cookie时判断注册接口是否存在
    接口存在->自动注册账号->爬取url->自动删除账号
    注册分为需要邮箱验证以及无需邮箱验证

GitLab
    /explore/projects 存在未授权，无需cookie也可以跑
    鉴权需要_gitlab_session=xxx
    注册接口都存在，但注册分为需要管理员验证以及无需管理员验证，而且存在验证token无法自动注册
    
Gitblit
    /repositories/ 存在未授权
    通过<a class="list"查找节点
    通过http://x.x.x.x/zip/?r=xx.git&format=zip下载
    鉴权需要Gitblit=
    登陆可爆破，自动尝试弱口令admin:admin登陆，若登陆成功后复用cookie获取节点
    
Jenkins
    当未设置鉴权时，dashboard存在未授权，即可访问到接口/api/json?pretty=true
    鉴权的话需要JSESSIONID.xxx=xxx
    注册为admin账号后台管理，无法前台注册

Confluence
    /login.action 登陆

语雀
    公开搜索 https://www.yuque.com/api/zsearch?q=%s&type=content&scope=/&tab=public&p=%s&sence=searchPage&time_horizon=
    知识库内文档列表 https://www.yuque.com/api/docs?book_id=%s
    markdown下载 https://www.yuque.com%s/markdown?attachment=true&latexcode=true&anchor=false&linebreak=false&&book_name=%s
石墨
    企业信息列表 https://shimo.im/lizard-api/org/departments/1/users?perPage=100&page=1
    团队空间列表 https://shimo.im/panda-api/file/spaces?orderBy=updatedAt
        空间文件列表 https://shimo.im/lizard-api/files?folder=%s
        导出markdown接口 https://shimo.im/lizard-api/office-gw/files/export?fileGuid=%s&type=md
        获取下载链接 https://shimo.im/lizard-api/office-gw/files/export/progress?taskId=%s
金山文档

```

### 使用指南

#### 指定文件 自动识别指纹后进行相关操作

```
go run main.go -file url.txt
```

#### Gitea

```
实例1:  拉取代码(含分支)
go run main.go --repo gitea --tar http://x.x.x.x

实例2:  拉取代码(不含分支, 仅master/main)

go run main.go --repo gitea --tar http://x.x.x.x -gb
```

#### Gogs

```
go run main.go --tar http://x.x.x.x --repo gogs

go run main.go --tar http://x.x.x.x -cookie i_like_gogs=xxx
```

#### GitLab

```
实例1:  拉取代码(含分支)
go run main.go --repo gitlab --tar http://x.x.x.x

实例2:  拉取代码(不含分支, 仅master)

go run main.go --repo gitlab --tar http://x.x.x.x -gb
```

#### Gitblit

```
实例1:  拉取代码(含分支)
go run main.go --repo gitblit --tar http://x.x.x.x

实例2:  拉取代码(不含分支, 仅master)

go run main.go --repo gitblit --tar http://x.x.x.x -gb
```

#### Jenkins

```
go run main.go --tar http://x.x.x.x --repo jenkins
```

#### Confluence

```
go run main.go --tar http://x.x.x.x --repo confluence

示例: 登陆后爬取
go run main.go --tar http://x.x.x.x --repo confluence -user xx -pwd xx
```
#### 语雀

```
go run main.go --repo yuque_open -q "搜索词" -p pageSize -cookie "cookie"
```

#### 石墨

```
go run main.go -cookie "cookie" -repo shimo
```

### Update

- 2023.12.14 v1.4.3.1 gitlab逻辑更新
- 2023.11.27 v1.4.3 适配gitea分页, 自动获取代码分支以及是否忽略分支
- 2023.11.16 v1.4.2 适配gitblit通用下载接口，自动获取代码分支以及是否忽略分支
- 2023.11.13 v1.4.1 适配gitlab多个版本差异代码逻辑, 自动获取代码分支以及是否忽略分支
- 2023.10.13 v1.4.0 新增石墨企业通讯录获取、团队空间以及个人空间markdown下载
- 2023.10.12 v1.4.0 新增语雀公开搜索下载markdown
- 2023.10.12 v1.3.1 Confluence爬取逻辑修复，新增Confluence登录session后爬取
- 2023.9.18 v1.3 新增Confluence未授权页面爬取
- 2023.9.14 v1.2 异步下载文件
- 2023.9.13 v1.1 新增-file指定文件，自动循环遍历url，以及自动识别url应用
- 2023.9.12 v1.0 发布

### todolist
