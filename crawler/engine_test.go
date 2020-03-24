package crawler

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestConcurrentEngine_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		seeds []Request
	}
	tests := []struct {
		name string
		e    *ConcurrentEngine
		args args
		want int
	}{
		// TODO: Add test cases.
		{"test",
			&ConcurrentEngine{
				Scheduler:       &QueueScheduler{Ctx: ctx},
				WorkerCount:     10,
				Ctx:             ctx,
				CheckExistOrAdd: func(interface{}) bool { return true },
			},
			args{seeds: []Request{
				Request{
					Item: "https://www.google.com/",
					ParseFunc: func(req Request) (ParseResult, error) {
						resp, err := http.Get(req.Item.(string))
						if err != nil {
							return ParseResult{}, err
						}
						defer resp.Body.Close()
						return ParseResult{
							Item: resp.StatusCode,
							Done: true,
						}, nil
					},
				}}},
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataChan := tt.e.Run(tt.args.seeds...)
		loop:
			for {
				select {
				case data, more := <-dataChan:
					if !more {
						break loop
					}
					got := data.(int)
					if got != tt.want {
						t.Errorf("ConcurrentEngine.Run() = %+v, want %v", got, tt.want)
					}
				case <-time.After(time.Second * 10):
					t.Fatal("Timeout")
				}
			}
		})
	}
}
