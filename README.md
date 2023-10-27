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

- 编译当前平台可执行文件：

```bash
go build main.go
```

- **交叉编译**指定平台可执行文件：

```bash
# 适用于Linux AArch64平台
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build main.go
```

```bash
# 适用于macOS amd64平台
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
```

```bash
# 适用于Windows amd64平台
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
