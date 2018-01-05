package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const config = "config.json"

type Config struct {
	lock     *sync.RWMutex
	Interval int             `json:"interval"`
	Channels map[string]bool `json:"channels"`
}

func (c *Config) setInterval(i int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Interval = i
	write(c)
}

func (c *Config) getInterval() time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()
	return time.Duration(c.Interval) * time.Second
}

func (c *Config) getChannels() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	i := 0
	chn := make([]string, len(c.Channels))
	for id := range c.Channels {
		chn[i] = id
		i++
	}

	return chn
}

func (c *Config) hasChannel(id string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.Channels[id]
	return ok
}

func (c *Config) addChannel(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Channels[id] = true
	write(c)
}

func (c *Config) removeChannel(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.Channels, id)
	write(c)
}

func getOrCreate() *Config {
	var cfg Config
	cfg.lock = &sync.RWMutex{}
	cfg.Interval = 15
	cfg.Channels = make(map[string]bool)

	data, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println(err)
		write(&cfg)
		return &cfg
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println(err)
		write(&cfg)
		return &cfg
	}

	return &cfg
}

func write(cfg *Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(config, data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
