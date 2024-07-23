package easyqueue

type configFunc func(*Config)

// 能够缓存的任务数量 QueuePartition * QueueCapacity
type Config struct {
	QueuePartition int // 队列分区数量
	QueueCapacity  int // 单个分区容量
	WorkersNum     int // 工作协程数量
}

// 设置分区数
func SetQueueParttion(partition int) configFunc {

	return func(c *Config) {
		c.QueuePartition = partition
	}
}

// 设置单个容量
func SetQueueCapacity(cap int) configFunc {

	return func(c *Config) {
		c.QueueCapacity = cap
	}

}

// 设置消费协程数
func SetWorkerNum(num int) configFunc {

	return func(c *Config) {
		c.WorkersNum = num
	}

}

type EasyQueue struct {
	config Config       // 配置
	queue  Queue        // 队列
	wg     *workerGroup // 消费组
}

func (eq *EasyQueue) Push(fn func()) WaitJob {
	job := newJob(fn)
	eq.queue.Push(job)
	return job
}

func CreateEasyQueue(cfs ...configFunc) *EasyQueue {

	conf := Config{}
	for _, cf := range cfs {
		cf(&conf)
	}

	eq := EasyQueue{
		config: conf,
	}

	eq.queue = createMultiJobQueue(conf.QueuePartition, conf.QueueCapacity)
	eq.wg = createWorkerMange(eq.queue, conf.WorkersNum)

	return &eq
}
