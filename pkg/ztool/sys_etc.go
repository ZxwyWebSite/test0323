//go:build !windows

package ztool

import "syscall"

var sys_Prioritys_etc = map[Sys_PriorityLev]int{
	Sys_PriorityLowest:   19,
	Sys_PriorityLower:    10,
	Sys_PriorityNormal:   0,
	Sys_PriorityHigher:   -10,
	Sys_PriorityHighest:  -15,
	Sys_PriorityRealtime: -20,
}

var Sys_GetPid = syscall.Getpid

func Sys_SetPriorityLev(pid int, lev Sys_PriorityLev) error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, pid, sys_Prioritys_etc[lev])
}
