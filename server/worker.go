package server

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type Engine struct {
	Ctx         context.Context
	WorkerCount int
	Scheduler   Scheduler
}

func (e *Engine) Run(requests ...Request) chan FinalResult {
	parseResultChan := make(chan _parseResult)
	finalResultChan := make(chan FinalResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler, parseResultChan)
	}

	for _, r := range requests {
		e.Scheduler.Submit(r)
	}

	go func() {
		var finalResultQ []FinalResult
		for {
			var activeFinalResult FinalResult
			var activeFinalResultChan chan<- FinalResult
			if len(finalResultQ) > 0 {
				activeFinalResultChan = finalResultChan
				activeFinalResult = finalResultQ[0]
			}
			select {
			case activeFinalResultChan <- activeFinalResult:
				finalResultQ = finalResultQ[1:]
			case parseResult := <-parseResultChan:
				finalResultQ = append(finalResultQ, FinalResult{
					Symbol:   parseResult.symbol,
					JSONData: parseResult.jsonResult,
					Error:    parseResult.err,
				})

				// 處理額外的 requests
				for _, request := range parseResult.requests {
					e.Scheduler.Submit(request)
				}
			case <-e.Ctx.Done():
				return
			}
		}
	}()

	return finalResultChan
}

func (e *Engine) createWorker(s Scheduler, out chan _parseResult) {
	in := make(chan Request)
	go func() {
		var parseResultQ []_parseResult
		s.WorkerReady(in)
		for {
			var activeResult _parseResult
			var activeResultChan chan<- _parseResult
			if len(parseResultQ) > 0 {
				activeResult = parseResultQ[0]
				activeResultChan = out
			}

			select {
			case request := <-in:
				parseResult, err := worker(request)
				if err != nil {
					log.Printf("worker: %s", err)
				}
				parseResultQ = append(parseResultQ, parseResult)

				s.WorkerReady(in)
			case activeResultChan <- activeResult:
				parseResultQ = parseResultQ[1:]
			case <-e.Ctx.Done():
				return
			}
		}
	}()
}

func worker(req Request) (_parseResult, error) {
	log.Printf("fetch %s\n", req.Symbol)
	parseResult, err := req.ParseFunc(req.Symbol)
	if err != nil {
		return _parseResult{symbol: req.Symbol, err: err},
			errors.WithMessage(err, fmt.Sprintf("req.ParseFunc(%s)", req.Symbol))
	}

	return _parseResult{
		symbol:     req.Symbol,
		jsonResult: parseResult,
	}, err
}
