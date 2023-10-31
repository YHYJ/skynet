# README

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-10-27 14:44:26 -->

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Usage](#usage)
* [Compile](#compile)
  * [当前平台](#当前平台)
  * [交叉编译](#交叉编译)
    * [Linux](#linux)
    * [macOS](#macos)
    * [Windows](#windows)

<!-- vim-markdown-toc -->

---

<!------------------------------------->
<!--      _                     _    -->
<!--  ___| | ___   _ _ __   ___| |_  -->
<!-- / __| |/ / | | | '_ \ / _ \ __| -->
<!-- \__ \   <| |_| | | | |  __/ |_  -->
<!-- |___/_|\_\\__, |_| |_|\___|\__| -->
<!--           |___/                 -->
<!------------------------------------->

---

一个网络管理器

## Usage

- `gui`子命令

    启动GUI版Skynet

- `httpserver`子命令

    启动CLI版Skynet

- `version`子命令

    查看程序版本信息

- `help`

    查看程序帮助信息

## Compile

### 当前平台

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/function.buildTime=`date +%s` -X github.com/yhyj/skynet/function.buildBy=$USER" -o skynet main.go
```

### 交叉编译

使用命令`go tool dist list`查看支持的平台

#### Linux

```bash
# 适用于Linux平台
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/function.buildTime=`date +%s` -X github.com/yhyj/skynet/function.buildBy=$USER" -o skynet main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64

#### macOS

```bash
# 适用于macOS平台
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/function.buildTime=`date +%s` -X github.com/yhyj/skynet/function.buildBy=$USER" -o skynet main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64

#### Windows

```powershell
# 适用于Windows平台
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -H windowsgui -X github.com/yhyj/skynet/function.buildTime=`date +%s` -X github.com/yhyj/skynet/function.buildBy=$USER" -o skynet main.go
```

> 使用`echo %PROCESSOR_ARCHITECTURE%`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64
