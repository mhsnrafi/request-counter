package counter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Counter struct {
	sync.Mutex
	Config      *Config
	Timestamps  []int64
	LastUpdated int64
}

type Config struct {
	TimeWindow      time.Duration
	PersistenceFile string
}

func NewCounter(config *Config) *Counter {
	return &Counter{
		Config: config,
	}
}

func (c *Counter) Count() int64 {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	minTime := now - int64(c.Config.TimeWindow.Seconds())

	// Best Case Scenario 1: If the last called time is less than min time we know there is no request in last 60 seconds
	if c.LastUpdated < minTime {
		c.LastUpdated = now
		// Clean the array
		c.Timestamps = []int64{now}
		return 1
	}

	c.LastUpdated = now

	// Best case scenario number 2: If the first index is within the range of window, all times are in the window
	if c.Timestamps[0] >= minTime {
		c.Timestamps = append(c.Timestamps, now)
		return int64(len(c.Timestamps))
	}

	// Try divide and conquer to find the index in log(n) time
	ind := findFirstIndexInWindow(c.Timestamps, minTime, 0, int64(len(c.Timestamps))-1)
	if ind == -1 {
		// A worst case scenario, but it wont happen because last updated case above will handle this
		// Clean the array
		c.Timestamps = []int64{now}
		return 1
	}
	// Slice the timestamp array to cleanup times outside window
	c.Timestamps = c.Timestamps[ind:]
	// Increment/Append the latest timestamp
	c.Timestamps = append(c.Timestamps, now)

	return int64(len(c.Timestamps))
}

func findFirstIndexInWindow(counts []int64, minTime int64, startIndex int64, endIndex int64) int64 {
	mid := ((endIndex - startIndex) / 2) + startIndex
	// base case
	if endIndex-startIndex == 0 {
		if counts[startIndex] >= minTime {
			return startIndex
		} else {
			return -1
		}
	}
	if counts[mid] >= minTime {
		return findFirstIndexInWindow(counts, minTime, startIndex, mid)
	} else {
		return findFirstIndexInWindow(counts, minTime, mid+1, endIndex)
	}
}

func (c *Counter) Save() error {
	c.Lock()
	defer c.Unlock()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Config.PersistenceFile, data, 0644)
}

func (c *Counter) Load() error {
	_, err := os.Stat(c.Config.PersistenceFile)
	if err != nil {
		err = c.createDefaultPersistenceFile()
		if err != nil {
			return err
		}
	}

	c.Lock()
	defer c.Unlock()

	data, err := ioutil.ReadFile(c.Config.PersistenceFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *Counter) createDefaultPersistenceFile() error {
	// Create the file with default values
	file, err := os.Create(c.Config.PersistenceFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the default JSON object (e.g., empty object)
	_, err = file.WriteString("{}")
	return err
}
