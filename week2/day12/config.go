package main

import "sync"

//定义 Config 结构体，包含 AppName string 和 Port int。
//定义 ConfigManager，包含 sync.Once 和 *Config。
//提供 GetConfig() 方法。
//初始化逻辑只允许执行一次。
//初始化时设置：
//- AppName = "GoWorker"
//- Port = 8080
//暂时不使用 Mutex 或 Atomic。
//保证多个 Goroutine 并发调用时安全。

type Config struct {
	AppName string
	Port    int
}

type ConfigManager struct {
	once   sync.Once
	config *Config
}

func (cm *ConfigManager) loadConfig() {
	cm.config = &Config{
		AppName: "GoWorker",
		Port:    8080,
	}
}

func (cm *ConfigManager) GetConfig() Config {
	cm.once.Do(func() {
		cm.loadConfig()
	})
	return *cm.config
}
