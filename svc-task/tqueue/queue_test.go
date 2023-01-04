package tqueue

import (
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

func TestTicker(t *testing.T) {
	type args struct {
		tick  *tmodel.Tick
		delay time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ticker is ticking successfully",
			args: args{
				tick: &tmodel.Tick{
					Ticker: time.NewTicker(time.Millisecond),
				},
				delay: 2 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go Ticker(tt.args.tick)
			time.Sleep(tt.args.delay)
			if tt.args.tick.Commit != true {
				t.Errorf("Ticker() failed to update commit")
			}
		})
	}
}

func TestUpdateTasksWorker(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		tick  *tmodel.Tick
		delay time.Duration
		try   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateTasksWorker update commit and executing flags properly",
			args: args{
				tick: &tmodel.Tick{
					Ticker: time.NewTicker(time.Millisecond),
				},
				delay: time.Millisecond,
				try:   1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				commit    bool
				executing bool
			)
			go UpdateTasksWorker(tt.args.tick)
			time.Sleep(tt.args.delay)
			for i := 0; i < tt.args.try; i++ {
				if tt.args.tick.Commit == true {
					commit = true
				}
				if tt.args.tick.Executing == true {
					executing = true
				}

				if commit && executing {
					break
				}
			}

			if commit && executing {
				t.Errorf("UpdateTasksWorker() failed to update commit and executing flags")
			}

		})
	}
}
