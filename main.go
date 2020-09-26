package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	engine "github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/stock/instance"
	"github.com/z-Wind/stock/stocker"
)

var (
	buildstamp = ""
	githash    = ""
	goversion  = ""

	// flag
	addr      string
	accountID string

	exePath  string
	stockers map[string]stocker.Stocker
)

func init() {
	flag.StringVar(&addr, "addr", "", "host:port, like localhost:6060 or 127.0.0.1:8090")
	flag.StringVar(&accountID, "accountID", "", "(option) TDAmeritrade account id")
}

func setting() {
	stockers = make(map[string]stocker.Stocker)

	var err error
	var path string
	exePath, err = getCurExePath()
	if err != nil {
		path = "./instance"
	} else {
		path = filepath.Join(exePath, "instance")
	}
	log.Printf("Current Path:%s", path)

	td, err := stocker.NewTDAmeritradeTLS(
		filepath.Join(path, "client_secret.json"),
		"TDAmeritrade-go.json", filepath.Join(path, "cert.pem"),
		filepath.Join(path, "key.pem"),
	)
	if err != nil {
		panic(err)
	}
	Register("TDAmeritrade", td)

	av, err := stocker.NewAlphavantage(instance.AlphaVantageKey)
	if err != nil {
		panic(err)
	}
	Register("alphavantage", av)

	twse, err := stocker.NewTWSE()
	if err != nil {
		panic(err)
	}
	Register("twse", twse)
}

func main() {
	flag.Parse()
	if addr == "" {
		fmt.Printf("addr is empty\n")

		flag.PrintDefaults()
		return
	}
	if accountID == "" {
		accountID = instance.AccountID
	}

	setting()

	engine.ELog.Start("engine.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	engine.ELog.SetFlags(0)
	defer engine.ELog.Stop()

	fmt.Println("=========================================")
	fmt.Printf("Git Commit Hash: %s\n", githash)
	fmt.Printf("Build Time : %s\n", buildstamp)
	fmt.Printf("Golang Version : %s\n", goversion)
	fmt.Println("=========================================")

	http.Handle("/", http.HandlerFunc(handleIndex))
	http.Handle("/quote", http.HandlerFunc(handleGet))
	http.Handle("/priceHistory", http.HandlerFunc(handleGet))
	http.Handle("/priceAdjHistory", http.HandlerFunc(handleGet))
	// http.Handle("/savedOrder", http.HandlerFunc(handleSavedOrder))

	fmt.Printf("start stock server: http://%s\n", addr)
	fmt.Println("=========================================")
	fmt.Printf("accountID : %q\n", accountID)
	fmt.Println("=========================================")
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Register 註冊可用 stocker
func Register(name string, s stocker.Stocker) {
	stockers[name] = s
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	template, err := parseTemplate(filepath.Join(exePath, "templates/index.html"), nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(template))
}

func makeQuoteParseFunc(f func(string) (float64, error)) func(engine.Request) (engine.ParseResult, error) {
	return func(req engine.Request) (engine.ParseResult, error) {
		parseResult := engine.ParseResult{
			Item:          nil,
			ExtraRequests: []engine.Request{},
			RedoRequests:  []engine.Request{},
			Done:          false,
		}

		symbol := req.Item.(string)

		price, err := f(symbol)
		if err != nil {
			switch err.(type) {
			case stocker.ErrorNoSupport, stocker.ErrorNoFound, stocker.ErrorFatal:
				parseResult.Done = true
			default:
				parseResult.RedoRequests = append(parseResult.RedoRequests, req)
			}

			return parseResult, err
		}

		rsp := Response{
			symbol: symbol,
			item:   price,
		}
		parseResult.Item = rsp
		parseResult.Done = true

		return parseResult, nil
	}
}

func makePriceHistoryParseFunc(f func(string) ([]*stocker.DatePrice, error)) func(engine.Request) (engine.ParseResult, error) {
	return func(req engine.Request) (engine.ParseResult, error) {
		parseResult := engine.ParseResult{
			Item:          nil,
			ExtraRequests: []engine.Request{},
			RedoRequests:  []engine.Request{},
			Done:          false,
		}

		symbol := req.Item.(string)

		history, err := f(symbol)
		if err != nil {
			switch err.(type) {
			case stocker.ErrorNoSupport, stocker.ErrorNoFound, stocker.ErrorFatal:
				parseResult.Done = true
			default:
				parseResult.RedoRequests = append(parseResult.RedoRequests, req)
			}

			return parseResult, err
		}

		// 日期由小到大
		sort.Slice(history, func(i, j int) bool {
			return time.Time(history[i].Date).Unix() < time.Time(history[j].Date).Unix()
		})

		rsp := Response{
			symbol: symbol,
			item:   history,
		}
		parseResult.Item = rsp
		parseResult.Done = true

		return parseResult, nil
	}
}

func reqToKey(req engine.Request) interface{} {
	key := req.Item.(string)

	return strings.ToUpper(key)
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := req.URL.Query()
	symbolQ := query.Get("symbols")
	if symbolQ == "" {
		http.Error(w, "symbols is empty", http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolQ, ",")
	symbols = removeDuplicates(symbols)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := engine.New(ctx, 10, reqToKey)

	requests := []engine.Request{}
	for _, symbol := range symbols {
		symbol = strings.ToUpper(symbol)
		for _, stk := range stockers {
			var parseFunc func(engine.Request) (engine.ParseResult, error)
			switch req.URL.Path {
			case "/quote":
				parseFunc = makeQuoteParseFunc(stk.Quote)
			case "/priceHistory":
				parseFunc = makePriceHistoryParseFunc(stk.PriceHistory)
			case "/priceAdjHistory":
				parseFunc = makePriceHistoryParseFunc(stk.PriceAdjHistory)
			default:
				http.Error(w, fmt.Sprintf("%s\n not support", req.URL.Path), http.StatusBadRequest)
				return
			}

			requests = append(requests, engine.Request{
				Item:      symbol,
				ParseFunc: parseFunc,
			})
		}
	}

	// 初始化
	prices := make(map[string]interface{}, len(symbols))
	for _, symbol := range symbols {
		prices[symbol] = nil
	}

	rspChan := e.Run(requests...)
	for rsp := range rspChan {
		result := rsp.(Response)

		prices[result.symbol] = result.item
		e.Recorder.Done(result.symbol)
	}

	err := json.NewEncoder(w).Encode(prices)
	if err != nil {
		log.Printf("json.NewEncoder error: %s\n", err)
	}
}
