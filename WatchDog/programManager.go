package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ProgramManager 管理多个程序的运行和监控
type ProgramManager struct {
	programs     []*ProgramInfo
	stopCh       chan struct{}
	restartQueue chan *ProgramInfo
	sigCh        chan os.Signal
	mutex        sync.RWMutex
}

// ProgramInfo 存储程序信息
type ProgramInfo struct {
	Name        string
	Path        string
	Args        []string
	Cmd         *exec.Cmd
	Pid         int
	Restarts    int
	MaxRestarts int
	Enabled     bool
}

// NewProgramManager 创建新的程序管理器
func NewProgramManager() *ProgramManager {
	return &ProgramManager{
		stopCh:       make(chan struct{}),
		restartQueue: make(chan *ProgramInfo, 100),
		sigCh:        make(chan os.Signal, 56),
	}
}

// AddProgram 添加要监控的程序
func (pm *ProgramManager) AddProgram(name, path string, args []string, enable bool) {
	program := &ProgramInfo{
		Name:        name,
		Path:        path,
		Args:        args,
		MaxRestarts: -1,
		Enabled:     enable,
	}

	if program.Path[len(program.Path)-1] != '/' {
		program.Path += "/"
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.programs = append(pm.programs, program)
	log.Printf("添加程序: %s, param: %s", program.Path+program.Name, args)
}

// StartProgram 启动单个程序
func (pm *ProgramManager) StartProgram(program *ProgramInfo) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var fullpath = program.Path + program.Name
	if !program.Enabled {
		return fmt.Errorf("程序 %s 已被禁用", fullpath)
	}

	// 检查程序文件是否存在
	fileInfo, err := os.Stat(fullpath)
	if os.IsNotExist(err) {
		return fmt.Errorf("程序文件不存在: %s", fullpath)
	}
	if err != nil {
		return fmt.Errorf("检查程序文件时出错: %v", err)
	}

	// 判断路径是否是目录
	if fileInfo.IsDir() {
		return fmt.Errorf("路径是目录而不是文件: %s", fullpath)
	}

	// 检查当前是否有执行权限
	currentMode := fileInfo.Mode()
	if currentMode&0111 == 0 {
		newMode := currentMode | 0111
		err := os.Chmod(fullpath, newMode)
		if err != nil {
			return fmt.Errorf("添加执行权限失败: %v", err)
		}
		fmt.Printf("当前文件没有执行权限，已添加: %s\n", fullpath)
	}

	// 创建命令
	program.Cmd = exec.Command(fullpath, program.Args...)
	program.Cmd.Stdout = os.Stdout
	program.Cmd.Stderr = os.Stderr
	program.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 启动程序
	if err := program.Cmd.Start(); err != nil {
		return fmt.Errorf("启动程序 %s 失败: %v", fullpath, err)
	}

	program.Pid = program.Cmd.Process.Pid
	log.Printf("程序 %s 已启动 (PID: %d)", fullpath, program.Pid)

	// 启动goroutine等待进程退出
	go pm.waitForProgramExit(program)

	return nil
}

// waitForProgramExit 等待程序退出并处理
func (pm *ProgramManager) waitForProgramExit(program *ProgramInfo) {
	err := program.Cmd.Wait()

	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if err != nil {
		log.Printf("程序 %s (PID: %d) 退出: %v", program.Name, program.Pid, err)
	} else {
		log.Printf("程序 %s (PID: %d) 正常退出", program.Name, program.Pid)
	}

	// 重置PID
	program.Pid = 0

	// 如果程序启用且未超过最大重启次数，加入重启队列
	if program.Enabled && (program.MaxRestarts == -1 || program.Restarts < program.MaxRestarts) {
		program.Restarts++
		log.Printf("程序 %s 将在1秒后重启 (重启次数: %d/%d)",
			program.Name, program.Restarts, program.MaxRestarts)

		// 非阻塞方式加入重启队列
		select {
		case pm.restartQueue <- program:
		default:
			log.Printf("重启队列已满，程序 %s 的重试将被延迟", program.Name)
			// 如果队列满，稍后重试
			go func(p *ProgramInfo) {
				time.Sleep(1 * time.Second)
				select {
				case pm.restartQueue <- p:
				case <-pm.stopCh:
				}
			}(program)
		}
	} else if program.Restarts >= program.MaxRestarts && program.MaxRestarts != -1 {
		log.Printf("程序 %s 已达到最大重启次数 (%d)，不再重启",
			program.Name, program.MaxRestarts)
	}
}

// StopProgram 停止单个程序
func (pm *ProgramManager) StopProgram(program *ProgramInfo) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if program.Cmd != nil && program.Cmd.Process != nil && program.Pid > 0 {
		log.Printf("正在停止程序 %s (PID: %d)", program.Name, program.Pid)

		// 发送SIGTERM信号
		if err := program.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("发送终止信号失败: %v", err)
		}

		// 设置超时强制终止
		done := make(chan error, 1)
		go func() {
			done <- program.Cmd.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			log.Printf("程序 %s 未正常退出，强制终止", program.Name)
			program.Cmd.Process.Kill()
		case err := <-done:
			if err != nil {
				log.Printf("程序 %s 退出: %v", program.Name, err)
			} else {
				log.Printf("程序 %s 已正常退出", program.Name)
			}
		}
	}

	program.Pid = 0
	program.Enabled = false
}

// setupSignalHandler 设置信号处理器
func (pm *ProgramManager) setupSignalHandler() {
	signal.Notify(pm.sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	go func() {
		for {
			select {
			case sig := <-pm.sigCh:
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					log.Printf("ProgramManager 收到信号 %v，正在停止所有程序...", sig)
					pm.StopAll()
					close(pm.stopCh)
					os.Exit(0)

				case syscall.SIGHUP:
					log.Println("收到 SIGHUP 信号，重新加载配置")
					// 这里可以添加配置重载逻辑
				}

			case <-pm.stopCh:
				return
			}
		}
	}()
}

// StartAll 启动所有程序
func (pm *ProgramManager) StartAll() error {
	var start_failed_count = 0
	for _, program := range pm.programs {
		if err := pm.StartProgram(program); err != nil {
			log.Printf("启动程序 %s 失败: %v", program.Name, err)
			start_failed_count++
		}
		time.Sleep(100 * time.Millisecond) // 稍微延迟，避免同时启动太多程序
	}

	if start_failed_count == len(pm.programs) {
		return fmt.Errorf("所有的程序都启动失败，监控程序退出")
	}

	return nil
}

// StopAll 停止所有程序
func (pm *ProgramManager) StopAll() {
	for _, program := range pm.programs {
		pm.StopProgram(program)
	}
}

// RestartProcessor 处理重启队列（用于非SIGCHLD的重启）
func (pm *ProgramManager) RestartProcessor() {
	for {
		select {
		case program := <-pm.restartQueue:
			time.Sleep(1 * time.Second)
			if err := pm.StartProgram(program); err != nil {
				log.Printf("重启程序 %s 失败: %v", program.Name, err)
			}

		case <-pm.stopCh:
			return
		}
	}
}

// Monitor 监控程序状态
func (pm *ProgramManager) Monitor() {
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-pm.stopCh:
			log.Println("停止监控")
			return

		case <-ticker.C:
			// 定期日志输出
			pm.logStatus()
		}
	}
}

// logStatus 输出当前状态日志
func (pm *ProgramManager) logStatus() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	log.Println("=== 程序状态报告 ===")
	for _, program := range pm.programs {
		status := "停止"
		if program.Pid > 0 {
			// 检查进程是否真的在运行
			if program.Cmd != nil && program.Cmd.Process != nil {
				if err := program.Cmd.Process.Signal(syscall.Signal(0)); err == nil {
					status = fmt.Sprintf("运行中 (PID: %d)", program.Pid)
				} else {
					status = "进程不存在"
				}
			}
		}

		log.Printf("程序: %s, 状态: %s, 重启次数: %d",
			program.Name, status, program.Restarts)
	}
}
