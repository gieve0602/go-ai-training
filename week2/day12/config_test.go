package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigManagerConcurrent(t *testing.T) {
	cm := &ConfigManager{}
	var wg sync.WaitGroup

	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			config := cm.GetConfig()
			assert.Equal(t, "GoWorker", config.AppName)
			assert.Equal(t, 8080, config.Port)
		}()
	}
	wg.Wait()
}
