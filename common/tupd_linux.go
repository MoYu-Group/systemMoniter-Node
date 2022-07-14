//go:build linux
// +build linux

package common

import (
	"os/exec"
	"strconv"
	"strings"
	"systemMoniter-Server/common"
)

func tupd() {
	byte1, err := exec.Command("bash", "-c", "ss -t|wc -l").Output()
	if err != nil {
		GetBasic.TCP = 0
		zap.L().Error("Get TCP count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get TCP count error:",err)
	} else {
		result := bytes2str(byte1)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindStringSubmatch(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch[0])
		if err != nil {
			zap.L().Error("Get TCP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get TCP count error::",err)
		}
		GetBasic.TCP = uint64(intNum)
	}
	byte2, err := exec.Command("bash", "-c", "ss -u|wc -l").Output()
	if err != nil {
		GetBasic.UDP = 0
		zap.L().Error("Get UDP count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get UDP count error:",err)
	} else {
		result := bytes2str(byte2)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindString(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch)
		if err != nil {
			zap.L().Error("Get UDP count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get UDP count error:",err)
		}
		GetBasic.UDP = uint64(intNum)
	}
	byte3, err := exec.Command("bash", "-c", "ps -ef|wc -l").Output()
	if err != nil {
		GetBasic.Process = 0
		zap.L().Error("Get process count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get process count error:",err)
	} else {
		result := bytes2str(byte3)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindString(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch)
		if err != nil {
			zap.L().Error("Get process count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get process count error:",err)
		}
		GetBasic.Process = uint64(intNum)
	}
	byte4, err := exec.Command("bash", "-c", "ps -eLf|wc -l").Output()
	if err != nil {
		GetBasic.Process = 0
		zap.L().Error("Get threads count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get threads count error:",err)
	} else {
		result := bytes2str(byte4)
		pattern := regexp.MustCompile(`[0-9]+`)
		strmatch := pattern.FindString(result)

		//fmt.Println(strmatch[0])
		// result = strings.Replace(result, "\r", "", -1)
		// result = strings.Replace(result, "\n", "", -1)
		intNum, err := strconv.Atoi(strmatch)
		if err != nil {
			zap.L().Error("Get threads count error:", zap.Error(err))
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get threads count error:",err)
		}
		GetBasic.Thread = uint64(intNum)
	}
}
