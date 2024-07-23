package main

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"CatMi-devops/utils/easyqueue"

	"github.com/gin-gonic/gin"
)

// 创建队列
var g_Queue = easyqueue.CreateEasyQueue(easyqueue.SetQueueParttion(1), easyqueue.SetQueueCapacity(100), easyqueue.SetWorkerNum(1))

var cacheNum atomic.Int64
var errCount atomic.Int64
var sucessCount atomic.Int64
var selloutCount atomic.Int64

var pingCount atomic.Int64

func main() {

	// 20个库存的商品
	cacheNum.Store(20)

	r := gin.Default()

	// 结果统计
	r.GET("/cache", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"商品剩余数量":      strconv.Itoa(int(cacheNum.Load())),
			"售罄->没买到的用户":  strconv.Itoa(int(selloutCount.Load())),
			"队列满->丢弃的用户数": strconv.Itoa(int(errCount.Load())),
			"成功->抢购的用户数":  strconv.Itoa(int(sucessCount.Load())),
			"抢购链接->总调用次数": strconv.Itoa(int(pingCount.Load())),
		})
	})

	// 增加商品库存
	r.GET("/add", func(c *gin.Context) {
		cacheNum.Add(20)

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// 秒杀商品
	r.GET("/ping", func(c *gin.Context) {

		pingCount.Add(1)

		// 如果已经售罄直接返回
		if cacheNum.Load() == 0 {
			selloutCount.Add(1)
			c.JSON(http.StatusOK, gin.H{
				"message": "商品售罄",
			})

			return
		}

		// 进入队列
		waitJob := g_Queue.Push(func() {

			// 写http请求中自己的业务逻辑....

			//1. 先判断是否售罄
			if cacheNum.Load() == 0 {
				selloutCount.Add(1)
				c.JSON(http.StatusOK, gin.H{
					"message": "商品售罄",
				})
			} else { //2.未售罄

				// 3.扣减库存
				cacheNum.Add(-1)
				// 4.下单逻辑..(模拟业务有点慢)
				time.Sleep(100 * time.Millisecond)
				// 5.告诉客户抢购成功
				sucessCount.Add(1)
				c.JSON(http.StatusOK, gin.H{
					"message": "抢购成功",
				})
			}

		})

		// 阻塞等待
		<-waitJob.Done()

		// 如果阻塞返回，err不为nil，说明上面的队列任务没有执行，这里做补偿逻辑
		if waitJob.Err() == easyqueue.ErrOverFlow { // 表示进入队列失败（队列满了)
			errCount.Add(1)
			c.JSON(http.StatusOK, gin.H{
				"message": "抢购人数过多，请重试...",
			})
		}
	})
	r.Run(":8089")
}
