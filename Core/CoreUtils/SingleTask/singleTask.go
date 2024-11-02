package SingleTask

import "time"

type ProcessorRunner func(interface{})

var EmptyDearDrop ProcessorRunner = func(i interface{}) {}

type Processor struct {
	taskChan         chan interface{}
	Runner           ProcessorRunner
	DearDrop         ProcessorRunner
	sleepMillisecond int
}

func NewProcessor(runner ProcessorRunner, dearDrop ProcessorRunner, sleepMillisecond int) *Processor {
	cp := &Processor{
		taskChan:         make(chan interface{}, 1), // 带缓冲区的通道，确保仅保留最新任务
		Runner:           runner,
		DearDrop:         dearDrop,
		sleepMillisecond: sleepMillisecond,
	}
	//go cp.run()
	return cp
}

func (cp *Processor) AddCoin(task interface{}) {
	// 尝试将任务发送到通道，覆盖旧任务
	select {
	case cp.taskChan <- task:
		// 成功发送，旧任务被替换
	default:
		// 通道已满，丢弃旧任务，再发送新任务
		cp.DearDrop(<-cp.taskChan)
		cp.taskChan <- task
	}
}

func (cp *Processor) Run() {
	for task := range cp.taskChan {
		cp.Runner(task)
		time.Sleep(time.Millisecond * time.Duration(cp.sleepMillisecond))
	}
}
