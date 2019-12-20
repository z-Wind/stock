package server

import "context"

type QueueScheduler struct {
	workerChan  chan chan Request
	requestChan chan Request
	Ctx         context.Context
}

func (s *QueueScheduler) Submit(r Request) {
	select {
	case s.requestChan <- r:
	case <-s.Ctx.Done():
	}
}

func (s *QueueScheduler) WorkerReady(w chan Request) {
	select {
	case s.workerChan <- w:
	case <-s.Ctx.Done():
	}
}

func (s *QueueScheduler) Run() {
	s.requestChan = make(chan Request)
	s.workerChan = make(chan chan Request)
	go func() {
		var requestQ []Request
		var workerQ []chan Request
		for {
			var activeRequest Request
			var activeWorker chan Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			case <-s.Ctx.Done():
				return
			}
		}
	}()
}
