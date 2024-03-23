// 系统调用 (beta); 分类：Sys_(系统)

package ztool

type (
	Sys_PriorityLev uint8
)

const (
	// 系统优先级 (https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setpriorityclass)

	Sys_PriorityLowest   Sys_PriorityLev = iota // 低 (其线程仅在系统空闲时运行的进程)
	Sys_PriorityLower                           // 低于正常
	Sys_PriorityNormal                          // 正常 (没有特殊调度需求的流程)
	Sys_PriorityHigher                          // 高于正常
	Sys_PriorityHighest                         // 高 (执行必须立即执行的时间关键型任务的流程)
	Sys_PriorityRealtime                        // 实时 (windows暂不支持)
)

/*
 Sys_GetPid() 获取适用于当前系统的Pid

 Sys_SetPriorityLev(pid, lev) 设置优先级

*/
