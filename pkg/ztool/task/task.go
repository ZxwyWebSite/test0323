// 计划任务(beta)

package task

import (
	"sync"
	"time"

	"github.com/ZxwyWebSite/ztool/logs"
)

type (
	TaskFunc func(*logs.Logger, int64) error
	// 任务对象
	Task struct {
		// 下次执行时间
		//  为什么是<下次>而不是<上次>：<上次>每次检测都需要加上间隔，<下次>可直接使用无需计算
		Next int64
		Wait int64    // 执行间隔
		Do   TaskFunc // 执行任务
		Err  []error  // 报错
	}
	// 控制器
	TaskWrapper struct {
		wait time.Duration    // 检测间隔
		task map[string]*Task // 任务列表
		maxe int              // 最多错误
	}
)

// 创建计划任务 <检测间隔, 最多错误>
/*
 注：在已开始执行后添加任务可能导致map并发读写错误，建议在初始化阶段完成添加
*/
func New(wait time.Duration, maxe int) *TaskWrapper {
	if wait < time.Second {
		panic(`task: 检测间隔必须大于等于1秒`)
	}
	if maxe < 1 {
		panic(`task: 最多错误数必须大于等于1`)
	}
	return &TaskWrapper{
		wait: wait,
		task: make(map[string]*Task),
		maxe: maxe,
	}
}

// 添加计划任务 <名称, 任务, 间隔, 立即执行>
/*
 注：
  为节省性能，任务并非实时检测，可能导致执行延后
  如需在最开始执行一次可将立即执行设为true (间隔时间过短可能导致误差，但不影响正常执行)
*/
func (w *TaskWrapper) Add(name string, do TaskFunc, wait int64, now bool) {
	// if wait < 0 {
	// 	panic(`task: 执行间隔不能小于0`)
	// }
	if _, ok := w.task[name]; ok {
		panic(`同名任务已存在`)
	}
	// w.task = append(w.task, task)
	next := time.Now().Unix()
	if !now {
		next += wait
	}
	w.task[name] = &Task{
		Next: next,
		Wait: wait,
		Do:   do,
		// Err:  make([]error, 0),
	}
}
func (w *TaskWrapper) AddTsk(name string, task *Task) {
	w.task[name] = task
}

// 执行计划任务 <日志, 保证首次运行> <同步>
/*
 注：同步可保证在程序结束前执行完此次检测，如无需求可忽略 (不能保证未开始的检测)
 同步函数等同于调用 (*sync.WaitGroup).Wait()
*/
func (w *TaskWrapper) Run(loger *logs.Logger /*, must bool*/) func() {
	if loger == nil {
		panic(`task: 日志输出不可为空`)
	}
	wg := new(sync.WaitGroup)
	// if must {
	// 	wg.Add(1)
	// }
	// var waiting bool
	l := loger.NewGroup(`CronTab`)
	if len(w.task) > 0 {
		go func() {
			for i := 0; true; i++ {
				wg.Add(1)
				// if waiting {
				// 	break
				// }
				time.Sleep(time.Millisecond * 300)
				l.Debug(`第 %v 次检测开始`, i)
				for name, tsk := range w.task {
					if now := time.Now().Unix(); tsk.Next <= now {
						// 重新计时以防"时间倒流"问题 //tsk.Next += tsk.Wait
						/*
						 [LOGS] [Info]  2024-02-17 00:05:46 [CronTab] [tsk2] Now:  1708099546
						 [LOGS] [Debug] 2024-02-17 00:05:46 [CronTab] [tsk2] Next: 1708099538
						*/
						tsk.Next = now + tsk.Wait
						l := l.AppGroup(name)
						l.Debug(`任务 %v 执行开始`, name)
						// 注：使用函数包裹以防panic时跳出循环导致无法正常输出
						func() {
							defer func() {
								if fatal := recover(); fatal != nil {
									l.Error(`致命错误: %s`, fatal)
									delete(w.task, name)
								}
							}()
							if e := tsk.Do(l, now); e != nil {
								l.Warn(`发生错误: %s`, e)
								tsk.Err = append(tsk.Err, e)
								if leng := len(tsk.Err); leng >= w.maxe {
									l.Error(`错误过多(%v=>%v): %+s`, leng, w.maxe, tsk.Err)
									delete(w.task, name)
									return
								}
								tsk.Next = now
							}
							l.Debug(`Next: %v`, tsk.Next)
						}()
						l.Debug(`执行结束`).Free()
					}
				}
				l.Debug(`检测完成`)
				wg.Done()
				time.Sleep(w.wait)
			}
			// if must {
			// 	wg.Done()
			// }
		}()
	}
	return wg.Wait //wg
	// return func() {
	// 	waiting = true
	// 	wg.Wait()
	// }
}

// func (w *TaskWrapper) Run2(loger *logs.Logger) func() {
// 	wg := new(sync.WaitGroup)
// 	l := loger.NewGroup(`CronTab`)
// 	if len(w.task) <= 0 {
// 		return func() {}
// 	}
// 	ticker := time.NewTicker(w.wait)
// 	go func() {
// 		for tick := range ticker.C {
// 			wg.Add(1)
// 			l.Debug(`检测开始`)
// 			for name, tsk := range w.task {
// 				if now := tick.Unix(); tsk.Next <= now {
// 					tsk.Next += tsk.Wait
// 					t := l.AppGroup(name)
// 					t.Debug(`任务 %v 执行开始`, name)
// 					func() {
// 						defer func() {
// 							if f := recover(); f != nil {
// 								t.Error(`致命错误: %s`, f)
// 								delete(w.task, name)
// 							}
// 						}()
// 						if e := tsk.Do(t, now); e != nil {
// 							t.Warn(`发生错误: %s`, e)
// 							tsk.Err = append(tsk.Err, e)
// 						}
// 						if leng := len(tsk.Err); leng >= w.maxe {
// 							t.Error(`错误过多(%v=>%v): %+s`, leng, w.maxe, tsk.Err)
// 							delete(w.task, name)
// 						}
// 					}()
// 					t.Debug(`执行结束 Next: %v`, tsk.Next).Free()
// 				}
// 			}
// 			l.Debug(`检测完成`)
// 			wg.Done()
// 		}
// 	}()
// 	return func() {
// 		wg.Wait()
// 		ticker.Stop()
// 	}
// }
