* 代码来源：https://github.com/bigwhite/experiments/tree/master/tcp-stream-proto
* 文章：https://tonybai.com/2021/07/28/classic-blocking-network-tcp-stream-protocol-parsing-practice-in-go/

## 服务说明
* demo1：Go经典阻塞式TCP协议流解析的实践
* demo1-with-metrics:
* demo2：
  * server端使用 bufio，减少直接read conn.Conn的次数
* demo3：
  * server端使用 bufio，减少直接read conn.Conn的次数
  * frame和packet中使用 mcache 减少gc
* demo3-with-metrics:
* demo4:
  * I/O多路复用
* demo4-with-metrics:
* demo5:
  * I/O 多路复用
  * 异步回应答
* demo5-with-metrics:

