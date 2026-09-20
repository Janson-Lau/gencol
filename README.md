# 目标：把 `gencol` 做成独立可复用插件，任意Go项目直接拿来用

现状：现在是 `go run ./internal/tools/gencol/main.go`，只能在当前项目跑。
目标：**全局安装命令** **`gencol`，任何项目在任意目录执行** **`gencol`** **就能生成字段常量文件**，类似 `gentool`。

## 整体方案分3步

1. 调整项目结构，把代码改成可编译为独立命令的 `main` 包
2. 发布/本地编译，用 `go install` 全局安装 `gencol`
3. 在别的项目直接调用，支持参数：扫描目录、输出目录

> 原理：`go install` 会编译 `main` 包，把二进制放到 `GOPATH/bin`，加到系统PATH后，powershell直接敲 `gencol` 就能执行。

## 第一步：调整目录结构（改成独立工具）

推荐单独仓库（也可以放在当前仓库，建议独立仓库方便跨项目复用）

```
gencol/
├── go.mod
├── go.sum
└── main.go   // 全部代码放这里，package main
```

> 关键点：文件头部必须是 `package main`，不能是别的包，否则无法编译成命令行程序。

把我们完整代码放到 `main.go`，包名改成 `package main`（已经是了）。

## 第二步：本地全局安装（重点！）

进入gencol工具目录，执行 go install

```
# 进入gencol目录（main.go所在目录）
cd D:\Code\Tools\Go\codeup.aliyun.com\gencol
go install .
```

> `go install .` 编译当前目录main包，生成gencol.exe，放到 `%GOPATH%\bin`

### 验证

新开PowerShell，直接敲：

```
gencol
```

出现帮助信息 `Usage: gencol <modelDir> <outDir>` 代表安装成功。

> 若提示找不到命令：
>
> 1. 查看GOPATH：`go env GOPATH`
> 2. 把 `C:\Users\susul\go\bin` 添加到系统环境变量PATH，重启终端

## 第三步：任意其他项目使用

### 新项目目录示例

```
new-project/
├── internal/models/role.go
└── internal/meta/  // 生成文件输出目录
```

进入新项目根目录，直接执行：

```
gencol ./internal/models ./internal/meta
```

自动扫描models下所有go文件，生成 `role_col.go` 到meta包。

### 搭配 go generate（推荐，项目内一键生成）

在项目任意一个go文件头部加一行注释（不用写在model结构体上！）

```
//go:generate gencol ./internal/models ./internal/meta
package models
```

之后只需要执行：

```
go generate ./...
```

自动调用全局gencol命令生成代码。

## 第四步：跨项目分发（两种方式）

### 方式1：推到git仓库（推荐，像gentool一样）

把gencol推到git仓库，例如 `codeup.aliyun.com/xxx/gencol`
之后任何电脑可以直接安装：

```
go install codeup.aliyun.com/xxx/gencol@latest
```

> 私有仓库需要配置git鉴权（ssh或者.netrc），和你之前拉私有module一样。

### 方式2：直接拷贝exe（简单，适合内网）

编译exe分发给团队：

```
go build -o gencol.exe .
```

放到所有人的 `go/bin`。

## 可选增强功能（后续迭代）

1. 增加参数 `-pkg`：自定义输出包名（现在固定meta）
2. 增加参数 `-parse-table-method`：解析 `TableName()` 函数获取表名（解决结构体tag语法坑）
3. 增加忽略标记：`//gencol:ignore`，跳过某些struct不生成
4. 增加注释：读取结构体comment写入生成代码注释
5. 增加版本参数：`gencol -v`

## 重要注意点

1. 工具本身**独立**，不依赖你的业务项目代码；它只做AST解析源码文本，不需要编译业务代码。
2. 业务项目只需要Go环境，不需要在项目内放gencol源码。
3. Windows / Linux / Mac 都可以编译，跨平台。

***

# 卸载 gencol（两种方式）

&#x20;

## 方式1：go install 安装的，直接删除二进制文件（Windows）

`go install` 只是把编译好的 `gencol.exe` 放到 `GOPATH/bin`，**没有别的依赖，删掉exe就等于卸载**

1. 查看你的 GOPATH

```
go env GOPATH

```

一般是 `C:\Users\susul\go` 2. 进入 bin 目录，删除 `gencol.exe`

```
# 直接删除
Remove-Item C:\Users\susul\go\bin\gencol.exe

```

1. 验证：新开 PowerShell，输入

```
gencol

```

提示 `'gencol' 不是内部或外部命令` → 卸载成功。

> 源码还在你本地文件夹（D:\Code\Tools\Go\codeup.aliyun.com\gencol），这个是源码，**不会自动删除**；你想删源码就手动删掉这个文件夹即可。

## 方式2：如果是远程仓库安装（类似 go install codeup.xxx/gencol\@latest）

Go 的 module 缓存不会自动清理二进制，依然是**删除 bin 下的 exe**。 如果想清理模块缓存（可选，不删也不影响）

```
go clean -modcache

```

⚠️ 这个命令会清空全部go模块缓存，所有下载过的依赖包都会被删掉，谨慎执行。

# 补充

- Go 的 `go install` **没有专门的** **`go uninstall`** **命令**！这是Go本身设计，只能手动删二进制。
- 如果你后面更新 gencol：直接重新 `go install .`，新的exe会直接覆盖旧的，不用手动卸载。

