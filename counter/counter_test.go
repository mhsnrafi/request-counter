package counter

import (
	"testing"
	"time"
)

const testPersistenceFile = "test_counts.json"

func TestFindFirstIndexInWindow(t *testing.T) {

	counts := []int64{1, 2, 2, 2, 2, 2, 2, 2, 3, 4, 5, 6, 7, 7, 7, 7, 7, 8, 9, 10, 11, 12, 13}
	index := findFirstIndexInWindow(counts, 7, 0, int64(len(counts)-1))
	if index != 12 {
		t.Errorf("Expected index to be 12 after loading, got %d", index)
		t.Errorf("%d", int64(len(counts))-index)
	}
}

func TestCounter_Count(t *testing.T) {

	tests := []struct {
		name         string
		countsOffset []int64
		want         int64
		countWindow  time.Duration
	}{
		{
			name:         "Three out of six times match the window",
			countsOffset: []int64{-3, -2, -1, 1, 2, 3},
			want:         4,
			countWindow:  20 * time.Second,
		},
		{
			name:         "Worst case scenario only last time matches the window",
			countsOffset: []int64{-3, -2, -1, -1, -1, 19},
			want:         2,
			countWindow:  20 * time.Second,
		},
		{
			name:         "No time matches the window",
			countsOffset: []int64{-3, -2, -1, -1, -1, -1},
			want:         1,
			countWindow:  20 * time.Second,
		},
		{
			name:         "Best case scenario every time is in the window",
			countsOffset: []int64{1, 2, 3, 4, 5, 6, 7, 19},
			want:         9,
			countWindow:  20 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().Unix()
			minWindowTime := now - int64(tt.countWindow.Seconds())
			var counts []int64
			for _, c := range tt.countsOffset {
				counts = append(counts, c+minWindowTime)
			}
			c := &Counter{
				Config: &Config{
					TimeWindow:      tt.countWindow,
					PersistenceFile: testPersistenceFile,
				},
				Timestamps:  counts,
				LastUpdated: now,
			}
			if got := c.Count(); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
