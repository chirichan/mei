# pwdgen

## Install

```bash
go install github.com/chirichan/mei/cmd/pwdgen@latest
go install github.com/chirichan/mei/cmd/bing15@latest
```

## Usage

pwdgen 用法

```
生成随机密码

Usage:
  pwdgen [flags]
  pwdgen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  csv2xykey   浏览器导出的 csv 格式的密码转为 xykey 格式。
  help        Help about any command

Flags:
  -h, --help         help for pwdgen
  -n, --length int   生成的密码长度, [6, 2048] (default 16)
  -l, --level int    生成的密码强度等级, 数字越大, 强度越高, [1, 4] (default 4)
  -o, --output int   输出方式, 1: 剪贴板, 2: 控制台 (default 1)
  -v, --version      版本

Use "pwdgen [command] --help" for more information about a command.
```
