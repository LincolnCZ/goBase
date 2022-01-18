# yylog

提取 https://git.yy.com/golang/gfy 中的日志库，并进行精简

## 使用实例

使用标准参数初始化日志。使用标准打印和日志会话打印

```go
package main

import (
	"context"
	"fmt"
	"time"

	"git.yy.com/media_svr/yylog"
	"go.uber.org/zap"
)

func main() {
	err := yylog.InitYYServerLog()
	if err != nil {
		fmt.Println("init yylog error:", err)
		return
	}
	yylog.Info("init yylog success")

	ctx := yylog.LogStart(context.Background(), zap.String("key1", "val1"), zap.Int("int1", 1))
	yylog.LogAppend(ctx, zap.String("key2", "val2"))
	yylog.LogFlush(ctx, "flush1", zap.String("key3", "val3"))
	yylog.LogFlush(ctx, "flush2", zap.String("key4", "val4"))
	time.Sleep(time.Second)
}
```