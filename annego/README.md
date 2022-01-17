## 简介

YY媒体后台用到的相关组件在Go语言中的封装，方便使用Go进行开发时，可以直接使用相关的组件。

包含的基本功能模块：

- config 配置文件解析，当前包含hostinfo.ini
- logger 日志打印，与C++日志打印相同，打印到syslog
- packet YY协议的封装和解封装，并提供反射方法
- yyserver 基于YY协议的基本网络框架
- s2s S2S节点发现的Go语言封装
- util 杂项

在设计时，尽量减少第三方库的依赖，只依赖标准库和少量轻量级库：

- gopkg.in/mgo.v2/bson bson格式序列化
- gopkg.in/ini.v1 解析ini文件
- github.com/stretchr/testify 单元测试使用

## 基本实例

在`go.mod`文件中添加以下配置，从而导入本包。

```
module "example"

require (
    git.yy.com/media_svr/annego v0.2.0
)
```

go 1.13 版本需要进行设置，支持从私有库查找。

```
go env -w GOPROXY=direct
go env -w GOSUMDB=off
```

实例代码：

```go
package main

import (
	"git.yy.com/media_svr/annego/logger"
)

func main() {
	logger.Info("logger before InitLog print to stdout")

	if err := logger.InitLog(); err != nil {
		logger.Warning("Init logger error %v", err)
		return
	}
	logger.Info("logger to syslog")
}
```

更多的功能使用，可以参考`example`中所含实例。

## 待做列表
