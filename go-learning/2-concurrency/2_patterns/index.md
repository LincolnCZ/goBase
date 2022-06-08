## 说明
* parallelLoop：并行循环
* pipeline：一个涉及三个 stage 的 pipeline
  * done取消机制
  * 并行循环
* du：遍历目录大小
  * done取消机制
  * 并行循环
  * token 控制并发数
* pubsub：发布订阅模式

## 任务编程
* Or-Done：Or-Done模式，如果有多个任务，只要其中任意一个任务完成
* fanIn-fanOut：扇入、扇出
* stream：
* map-reduce：