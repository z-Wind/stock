package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/z-Wind/stock/api/alphavantage"
	"github.com/z-Wind/stock/server"

	"github.com/pkg/errors"
)

type priceAdjHistoryFunc func(string) (datePrices, error)

type PriceAdjHistory struct {
	funcs []priceAdjHistoryFunc
}

// RegisterPriceAdjHistory 各種歷史詢價來源 adj
func (p *PriceAdjHistory) RegisterPriceAdjHistory(f priceAdjHistoryFunc) {
	p.funcs = append(p.funcs, f)
}

func (p *PriceAdjHistory) makeParseFunc(f priceAdjHistoryFunc) server.ParseFunc {
	return func(symbol string) ([]byte, error) {
		result, err := f(symbol)
		if err != nil {
			return nil, errors.WithMessage(err, "priceAdjHistoryFunc")
		}
		resultJson, err := json.Marshal(result)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal")
		}
		return resultJson, nil
	}
}

func (p *PriceAdjHistory) HandlerFunc(w http.ResponseWriter, req *http.Request) {
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
		for _, f := range p.funcs {
			count++
			requests = append(requests, server.Request{
				Symbol:    symbol,
				ParseFunc: p.makeParseFunc(f),
			})
		}
	}

	m := make(map[string]datePrices, len(symbols))
	resultChan := e.Run(requests...)
	for result := range resultChan {
		count--
		if result.JSONData == nil {
			log.Printf("%s result.JSONData is nil\n", result.Symbol)
			if strings.Contains(result.Error.Error(), "Invalid API call") {
				m[result.Symbol] = nil
			}
		} else {
			var temp datePrices
			err := json.Unmarshal(result.JSONData, &temp)
			if err != nil {
				log.Printf("json.Unmarshal error: %s\n", err)
				if strings.Contains(result.Error.Error(), "Invalid API call") {
					m[result.Symbol] = nil
				}
			} else {
				sort.Slice(temp, func(i, j int) bool {
					return temp[i].Date < temp[j].Date
				})
				m[result.Symbol] = temp
			}
		}

		// 請求全部完成 或 資料全拿到
		if count == 0 || len(m) == len(symbols) {
			break
		}
	}
	cancel()
	if len(m) > 0 {
		err := json.NewEncoder(w).Encode(m)
		if err != nil {
			log.Printf("json.NewEncoder error: %s\n", err)
		}
	}
}

func priceAdjHistoryAlpha(symbol string) (datePrices, error) {
	p, err := alpha.DailyAdj(symbol, alphavantage.OutputSizeFull)
	if err != nil {
		return nil, errors.WithMessage(err, "alpha.DailyAdj")
	}

	timeSeries := make(datePrices, len(p.TimeSeriesDaily))
	i := 0
	for key, trade := range p.TimeSeriesDaily {
		t := datePrice{
			Date:     key,
			Open:     trade.Open,
			High:     trade.High,
			Low:      trade.Low,
			Close:    trade.Close,
			CloseAdj: trade.AdjustedClose,
			Volume:   trade.Volume_,
		}
		timeSeries[i] = &t
		i++
	}

	return timeSeries, nil
}
