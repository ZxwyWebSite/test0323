// Shell操作 (beta); 分类: Cmd_(终端)

package ztool

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"

	// "strings"

	"github.com/ZxwyWebSite/ztool/x/bytesconv"
)

// 获取对应系统的Shell调用接口
func getShell(command string) *exec.Cmd {
	if runtime.GOOS == `windows` {
		return exec.Command(`cmd`, `/C`, command)
	}
	return exec.Command(`/usr/bin/env`, `bash`, `-c`, command)
}

// 等待(同步)执行Shell命令
func Cmd_aWaitExec(command string) (string, error) {
	cmd := getShell(command)
	var out bytes.Buffer
	// var ert bytes.Buffer
	cmd.Stdout = &out
	// cmd.Stderr = &ert
	cmd.Stderr = cmd.Stdout
	err := cmd.Start()
	if err != nil {
		// return ert.String(), err
		return out.String(), err
	}
	err = cmd.Wait()
	return out.String(), err
}

// 同步操作 (直接输出到Stdout)
func Cmd_aSyncExec(command string) (err error) {
	cmd := getShell(command)
	// out, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return
	// }
	// defer out.Close()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = cmd.Stdout
	// for k, v := range env {
	// 	cmd.Env = append(cmd.Env, FastStrConcat(k, `=`, v))
	// }
	if err = cmd.Start(); err != nil {
		return
	}
	// buff := make([]byte, 8)
	// for {
	// 	len, err := out.Read(buff)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	fmt.Print(string(buff[:len]))
	// }
	cmd.Wait()
	return
}

// 快速输出字符串
func Cmd_FastPrint(str string) error {
	// if str == `` {
	// 	return nil
	// }
	// _, err := io.Copy(os.Stdout, strings.NewReader(str))
	// return err
	return Cmd_FastFprint(os.Stdout, str)
}
func Cmd_FastFprint(dst io.Writer, src string) error {
	_, err := dst.Write(bytesconv.StringToBytes(src))
	// _, err := io.CopyBuffer(dst, strings.NewReader(src), make([]byte, 512))
	return err
}
func Cmd_FastPrintln(str string) error {
	return Cmd_FastPrint(Str_FastConcat(str, "\n"))
}
