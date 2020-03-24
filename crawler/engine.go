package crawler

import (
	"context"
	"fmt"

	"github.com/z-Wind/getNovel/util"
)

// ELog engine log
var ELog util.Log

func New(ctx context.Context, workerNum int, isProcessedOrAdd func(interface{}) bool, isDone func(interface{}) bool) *ConcurrentEngine {
	return &ConcurrentEngine{
		Scheduler:        &QueueScheduler{Ctx: ctx},
		WorkerCount:      workerNum,
		Ctx:              ctx,
		isProcessedOrAdd: isProcessedOrAdd,
		isDone:           isDone,
	}
}

// ConcurrentEngine 負責處理對外與建立 worker
type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	Ctx              context.Context
	NumTasks         int
	isProcessedOrAdd func(interface{}) bool
	isDone           func(interface{}) bool
}

// Run 開始運作
func (e *ConcurrentEngine) Run(seeds ...Request) chan interface{} {
	parseResultChan := make(chan ParseResult)
	dataChan := make(chan interface{})

	e.Scheduler.Run()
	e.NumTasks = len(seeds)
	ELog.Printf("Initial Tasks: %d\n", e.NumTasks)

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(parseResultChan, e.Scheduler)
	}

	for _, req := range seeds {
		e.Scheduler.Submit(req)
	}

	go func() {
		// 確認是否有任務
		if e.NumTasks == 0 {
			close(dataChan)
		}
		// 用 queue 先存起來，防止阻塞
		var dataQ []interface{}

		for {
			var activeData interface{}
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeDataChan chan<- interface{}
			if len(dataQ) > 0 {
				activeData = dataQ[0]
				activeDataChan = dataChan
			}
			if e.NumTasks == 0 && len(dataQ) == 0 {
				ELog.Printf("Finish =============================================\n")
				close(dataChan)
			}

			select {
			case activeDataChan <- activeData:
				ELog.LPrintf("%-30s DataChan <- Data\n", "Engine.Run")
				dataQ = dataQ[1:]
			case parseResult := <-parseResultChan:
				ELog.LPrintf("%-30s parseResult := <-parseResultChan", "Engine.Run")
				if parseResult.Item != nil {
					ELog.LPrintf("%-30s Get Result\n", "Engine.Run")
					dataQ = append(dataQ, parseResult.Item)
				}

				ELog.LPrintf("%-30s Done: %v\n", "Engine.Run", parseResult.Done)
				if parseResult.Done {
					e.NumTasks--
				} else {
					lastIndex := len(parseResult.Requests) - 1
					req := parseResult.Requests[lastIndex]
					if !e.isDone(req) {
						e.Scheduler.Submit(req)
						ELog.LPrintf("%-30s Recover Task: %+v", "Engine.Run", req.Item)
					} else {
						e.NumTasks--
						ELog.LPrintf("%-30s Done by other worker Task: %+v", "Engine.Run", req.Item)
					}

					// remove original request
					parseResult.Requests = parseResult.Requests[:lastIndex]
				}

				// 排入新增的 requests
				for _, req := range parseResult.Requests {
					if !e.isProcessedOrAdd(req) {
						e.Scheduler.Submit(req)
						e.NumTasks++
						ELog.LPrintf("%-30s Add Task: %+v", "Engine.Run", req.Item)
					}
				}
				ELog.Printf("%-30s To Do Tasks: %d\n", "Engine.Run", e.NumTasks)
			case <-e.Ctx.Done():
				ELog.Printf("%-30s ConcurrentEngine.Run.Done\n", "Engine.Run")
				return
			}
		}
	}()

	return dataChan
}

func (e *ConcurrentEngine) createWorker(parseResultChan chan<- ParseResult, s Scheduler) {
	requestChan := make(chan Request)

	go func() {
		// 用 queue 先存起來，防止阻塞
		var parseResultQ []ParseResult

		s.WorkerReady(requestChan)

		for {
			var activeResult ParseResult
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeResultChan chan<- ParseResult
			if len(parseResultQ) > 0 {
				activeResult = parseResultQ[0]
				activeResultChan = parseResultChan
			}

			select {
			case activeResultChan <- activeResult:
				ELog.LPrintf("%-30s ResultChan <- Result\n", fmt.Sprintf("worker(%v)", requestChan))
				parseResultQ = parseResultQ[1:]
			case request := <-requestChan:
				ELog.LPrintf("%-30s request := <-requestChan\n", fmt.Sprintf("worker(%v)", requestChan))
				ELog.LPrintf("%-30s Process Request: %+v", fmt.Sprintf("worker(%v)", requestChan), request.Item)
				result := worker(request)
				parseResultQ = append(parseResultQ, result)
				s.WorkerReady(requestChan)
			case <-e.Ctx.Done():
				ELog.Printf("%-30s ConcurrentEngine.createWorker.Done\n", fmt.Sprintf("worker(%v)", requestChan))
				return
			}
		}
	}()
}

func worker(req Request) ParseResult {
	parseResult, err := req.ParseFunc(req)
	if err != nil {
		ELog.Printf("worker: req.ParseFunc: err:%s\n", err)
		ELog.LPrintf("%+v\n", parseResult)
		return parseResult
	}

	return parseResult
}
