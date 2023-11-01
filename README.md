<h1 align="center">Skynet</h1>

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-10-27 14:44:26 -->

---


<p align="center">
  <a href="https://github.com/YHYJ/skynet/actions/workflows/release.yml"><img src="https://github.com/YHYJ/skynet/actions/workflows/release.yml/badge.svg" alt="Go build and release by GoReleaser"></a>
</p>

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Usage](#usage)
* [Screenshot](#screenshot)
  * [CLI version](#cli-version)
  * [GUI version](#gui-version)
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

  启动 GUI 版 Skynet

- `http`子命令

  在 CLI 启动 HTTP 服务

- `version`子命令

  查看程序版本信息

- `help`

  查看程序帮助信息

## Screenshot

### CLI version

![Skynet CLI version](resources/screenshots/cli-version.png)

### GUI version

![Skynet GUI version](resources/screenshots/gui-version.png)

## Compile

### 当前平台

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/skynet/general.BuildTime=`date +%s` -X github.com/yhyj/skynet/general.BuildBy=$USER" -o build/skynet main.go
```

### 交叉编译

使用命令`go tool dist list`查看支持的平台

#### Linux

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/skynet/general.BuildTime=`date +%s` -X github.com/yhyj/skynet/general.BuildBy=$USER" -o build/skynet main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### macOS

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/skynet/general.BuildTime=`date +%s` -X github.com/yhyj/skynet/general.BuildBy=$USER" -o build/skynet main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### Windows

```powershell
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -H windowsgui -X github.com/yhyj/skynet/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/skynet/general.BuildTime=`date +%s` -X github.com/yhyj/skynet/general.BuildBy=$USER" -o build/skynet.exe main.go
```

> 使用`echo %PROCESSOR_ARCHITECTURE%`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64
