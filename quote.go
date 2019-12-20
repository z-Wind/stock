package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/z-Wind/stock/api/gotd/api/types"
	"github.com/z-Wind/stock/server"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type quoteFunc func(string) (float64, error)

type Quote struct {
	funcs []quoteFunc
}

// RegisterQuote 各種詢價來源
func (q *Quote) RegisterQuote(f quoteFunc) {
	q.funcs = append(q.funcs, f)
}

func (q *Quote) makeParseFunc(f quoteFunc) server.ParseFunc {
	return func(symbol string) ([]byte, error) {
		result, err := f(symbol)
		if err != nil {
			return nil, errors.WithMessage(err, "quoteFunc")
		}
		resultJson, err := json.Marshal(result)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal")
		}
		return resultJson, nil
	}
}

func (q *Quote) HandlerFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	symbolQ, ok := req.URL.Query()["symbols"]
	if !ok || len(symbolQ) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolQ[0], ",")
	symbols = removeDuplicates(symbols)

	ctx, cancel := context.WithCancel(context.Background())
	e := server.Engine{
		Scheduler:   &server.QueueScheduler{Ctx: ctx},
		WorkerCount: 10,
		Ctx:         ctx,
	}

	count := 0
	requests := []server.Request{}
	for _, symbol := range symbols {
		for _, f := range q.funcs {
			count++
			requests = append(requests, server.Request{
				Symbol:    symbol,
				ParseFunc: q.makeParseFunc(f),
			})
		}
	}

	m := make(map[string]float64, len(symbols))
	resultChan := e.Run(requests...)
	for result := range resultChan {
		count--
		if result.JSONData == nil {
			log.Printf("%s result.JSONData is nil\n", result.Symbol)
		} else {
			var temp float64
			err := json.Unmarshal(result.JSONData, &temp)
			if err != nil {
				log.Printf("json.Unmarshal error: %s\n", err)
			} else {
				m[result.Symbol] = temp
			}
		}
		// 請求全部完成 或 資料全拿到
		if count == 0 || len(m) == len(symbols) {
			break
		}
	}
	cancel()
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Printf("json.NewEncoder error: %s\n", err)
	}
}
func quoteTD(symbol string) (float64, error) {
	err := td.RefreshAccessTokenOrNot()
	if err != nil {
		return 0.0, errors.WithMessage(err, "td.RefreshAccessTokenOrNot")
	}

	q, err := td.GetQuote(symbol)
	if err != nil {
		return 0.0, errors.WithMessage(err, "td.GetQuote")
	}

	switch data := q.Data.(type) {
	case *types.MutualFundQ:
		return data.NetChange, nil
	case *types.Future:
		return data.Mark, nil
	case *types.FutureOptions:
		return data.Mark, nil
	case *types.Index:
		return data.NetChange, nil
	case *types.Option:
		return data.Mark, nil
	case *types.Forex:
		return data.Mark, nil
	case *types.ETF:
		return data.Mark, nil
	default:
		return 0.0, fmt.Errorf("Not support type %T", data)
	}
}

func quoteAlpha(symbol string) (float64, error) {
	q, err := alpha.LatestPrice(symbol)
	if err != nil {
		return 0.0, errors.WithMessage(err, "alpha.LatestPrice")
	}

	return q.Price, nil
}

func quoteTWSE(symbol string) (float64, error) {
	q, err := tw.Quote(symbol)
	if err != nil {
		return 0.0, errors.WithMessage(err, "tw.Quote")
	}
	if len(q) == 0 {
		return 0.0, fmt.Errorf("tw.Quote: No Data")
	}

	return q[0].TradePrice, nil
}
