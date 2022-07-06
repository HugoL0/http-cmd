package ipfs

import (
	"fmt"
	"github.com/ipfs/go-log/v2"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var logs = log.Logger("ipfsHelper/ipfs")

var Bin = "ipfs"

var Env []string

func getIpfsCmd(sub ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: Bin,
		Args: append([]string{Bin}, sub...),
		Env:  Env,
	}
}

func getIpfsCmdOS(sub ...string) *exec.Cmd {
	return &exec.Cmd{
		Path:   Bin,
		Args:   append([]string{Bin}, sub...),
		Env:    Env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func AddFile2Ipfs(fileName string) error {
	return getIpfsCmdOS("blockchain", "file", "add", fileName).Run()
}

func StartIpfs() error {
	c := getIpfsCmd("daemon")
	err := c.Start()
	if err != nil {
		logs.Error("ipfs daemon启动失败：%s", err.Error())
		return err
	}
	logs.Info("ipfs daemon启动成功")
	return c.Wait()
}

func StartIpfsOnce() error {
	return getIpfsCmd("daemon").Start()
}

func RunIpfs(reCh, doneCh chan struct{}) {
	err := StartIpfs()
	if err != nil {
		logs.Errorf("ipfs daemon 异常退出：%s", err.Error())
		close(doneCh)
		return
	}
	logs.Warn("ipfs daemon 正常退出")
	reCh <- struct{}{}
}

func Daemon(errCh chan error, reCh chan struct{}, doneCh chan struct{}) {
	go RunIpfs(reCh, doneCh)
	m := map[time.Time]struct{}{}
	interval := 10 * time.Second
	limit := 3
	// 定时重启
	go func(i int, t int) {
		if i <= 0 {
			logs.Infof("未启用定时重启")
		}
		logs.Infof("每%v天%v点重启一次", i, t)
		var now, next time.Time
		var timer *time.Timer
		for {
			now = time.Now()
			next = time.Date(now.Year(), now.Month(), now.Day()+i, t, 0, 0, 0, now.Location())
			timer = time.NewTimer(next.Sub(now))
			fmt.Println(next.Sub(now))
			<-timer.C
			err := Shutdown()
			if err != nil {
				logs.Errorf("重启失败：%v", err)
			}
		}
	}(1, 2)

	for {
		select {
		case <-reCh:
			// 短时间频繁重启  比如 10s内重启三次 直接退出
			if ok := ComputeFailFrequency(m, interval, limit); ok {
				go RunIpfs(reCh, doneCh)
			} else {
				err := Shutdown()
				if err != nil {
					logs.Errorf(" ipfs shutdown failed：%s", err.Error())
				}
				return
			}
		case <-doneCh:
			err := Shutdown()
			if err != nil {
				logs.Errorf(" ipfs shutdown failed：%s", err.Error())
			}
			return
		}
	}
}

func ComputeFailFrequency(m map[time.Time]struct{}, interval time.Duration, limit int) bool {
	now := time.Now()
	cnt := 0
	for t := range m {
		if now.Sub(t).Seconds() > interval.Seconds() {
			delete(m, t)
		} else {
			cnt++
			if cnt >= limit {
				return false
			}
		}
	}
	m[now] = struct{}{}
	return true
}

func Shutdown() error {
	return getIpfsCmd("shutdown").Run()
}

func ReStart() error {
	err := Shutdown()
	if err != nil {
		logs.Warnf("关闭ipfs失败：%s", err.Error())
	}
	// 如果ipfs守护进程没有运行
	return getIpfsCmd("daemon").Start()
}

func Stop(pid int) error {
	if pid < 0 {
		return fmt.Errorf("pid wrong")
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGKILL)
}

func Update(old string, ipfsBin string) error {
	err := os.Rename(old, ipfsBin)
	if err != nil {
		return err
	}
	f, err := os.Open(ipfsBin)
	if err != nil {
		return err
	}
	err = f.Chmod(0777)

	return getIpfsCmd("shutdown").Run()
}
