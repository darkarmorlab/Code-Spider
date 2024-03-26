env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_darwin_amd64
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_386
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_amd64
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_arm
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_arm64
env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_windows_386.exe
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_windows_amd64.exe
env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_mips64
env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_mips64le
env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_mips
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -trimpath -ldflags "-s -w" -o ./code-spider-release/code_spider_linux_mipsle
