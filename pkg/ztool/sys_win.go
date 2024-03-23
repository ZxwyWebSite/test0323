//go:build windows

package ztool

import (
	"golang.org/x/sys/windows"
)

var sys_Prioritys_win = map[Sys_PriorityLev]uint32{
	Sys_PriorityLowest:   windows.IDLE_PRIORITY_CLASS,
	Sys_PriorityLower:    windows.BELOW_NORMAL_PRIORITY_CLASS,
	Sys_PriorityNormal:   windows.NORMAL_PRIORITY_CLASS,
	Sys_PriorityHigher:   windows.ABOVE_NORMAL_PRIORITY_CLASS,
	Sys_PriorityHighest:  windows.HIGH_PRIORITY_CLASS,
	Sys_PriorityRealtime: windows.REALTIME_PRIORITY_CLASS,
}

var Sys_GetPid = windows.CurrentProcess

func Sys_SetPriorityLev(pid windows.Handle, lev Sys_PriorityLev) error {
	return windows.SetPriorityClass(pid, sys_Prioritys_win[lev])
}
