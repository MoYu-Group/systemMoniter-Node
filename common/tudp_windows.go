package common

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap"
)

func tupd() {
	cmd, err := Command("cmd", "/c netstat -an|find \"TCP\" /c")
	if err != nil {
		GetBasic.TCP = 0
		zap.L().Error("Get TCP count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet TCP count error:",err)
	} else {
		byte1, err := cmd.Output()
		if err != nil {
			zap.L().Error("Get TCP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet TCP count error:",err)
		}
		result := Bytes2str(byte1)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindString(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch)
		if err != nil {
			zap.L().Error("Get TCP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet TCP count error:",err)
		}
		GetBasic.TCP = uint64(intNum)
	}
	cmd2, err := Command("cmd", "/c netstat -an|find \"UDP\" /c")
	if err != nil {
		GetBasic.UDP = 0
		zap.L().Error("Get UDP count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet UDP count error:",err)
	} else {
		byte2, err := cmd2.Output()
		if err != nil {
			zap.L().Error("Get TCP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet TCP count error:",err)
		}
		result := Bytes2str(byte2)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindString(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch)
		if err != nil {
			zap.L().Error("Get UDP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet UDP count error:",err)
		}
		GetBasic.UDP = uint64(intNum)
	}
	pids, err := process.Processes()
	if err != nil {
		zap.L().Error("Get process count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," zGet process count error:",err)
	} else {
		GetBasic.Process = uint64(len(pids))
	}
	GetBasic.Thread = 0

}

func Command(name, args string) (*exec.Cmd, error) {
	// golang 使用 exec.Comand 运行含管道的 cmd 命令会产生问题（如 netstat -an | find "TCP" /c），因此使用此办法调用
	// 参考：https://studygolang.com/topics/10284
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err != nil {
			return nil, err
		}
		name = lp
	}
	return &exec.Cmd{
		Path:        name,
		SysProcAttr: &syscall.SysProcAttr{CmdLine: name + " " + args},
	}, nil
}
