* 代码来源：https://github.com/bigwhite/experiments/tree/master/tcp-stream-proto
* 文章：https://tonybai.com/2021/07/28/classic-blocking-network-tcp-stream-protocol-parsing-practice-in-go/

## 服务说明
* demo1：Go 经典阻塞式 TCP 协议流解析的实践
* demo1-with-metrics:
* demo2：
  * server 端使用 bufio，减少直接 read conn.Conn 的次数
* demo3：
  * server 端使用 bufio，减少直接 read conn.Conn 的次数
  * frame 和 packet 中使用 mcache 减少gc
* demo3-with-metrics:
* demo4:
  * I/O 多路复用
* demo4-with-metrics:
* demo5:
  * I/O 多路复用
  * 异步回应答
* demo5-with-metrics:

