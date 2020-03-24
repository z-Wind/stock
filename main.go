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
	"strings"

	"github.com/z-Wind/stock/crawler"
	"github.com/z-Wind/stock/instance"
	"github.com/z-Wind/stock/stocker"
)

var (
	buildstamp = ""
	githash    = ""
	goversion  = ""

	// flag
	addr        string
	redirectURL string
	accountID   string

	stockers []stocker.Stocker
)

func init() {
	flag.StringVar(&addr, "addr", "", "host:port, like localhost:6060 or 127.0.0.1:8090")
	flag.StringVar(&accountID, "accountID", "", "(option) TDAmeritrade account id")
	flag.StringVar(&redirectURL, "redirectURL", "", "TDAmeritrade redirectURL")
}

func setting() {
	path, err := getCurExePath()
	if err != nil {
		path = "./instance"
	} else {
		path = filepath.Join(path, "instance")
	}
	log.Printf("Current Path:%s", path)

	td, err := stocker.NewTDAmeritradeTLS(
		redirectURL,
		filepath.Join(path, "client_secret.json"),
		"TDAmeritrade-go.json", filepath.Join(path, "cert.pem"),
		filepath.Join(path, "key.pem"),
	)
	if err != nil {
		panic(err)
	}
	Register(td)

	av, err := stocker.NewAlphavantage(instance.AlphaVantageKey)
	if err != nil {
		panic(err)
	}
	Register(av)

	twse, err := stocker.NewTWSE()
	if err != nil {
		panic(err)
	}
	Register(twse)
}

func main() {
	flag.Parse()
	if addr == "" {
		fmt.Printf("addr is empty\n")

		flag.PrintDefaults()
		return
	}
	if redirectURL == "" {
		fmt.Printf("redirectURL is empty\n")

		flag.PrintDefaults()
		return
	}
	if accountID == "" {
		accountID = instance.AccountID
	}

	setting()

	crawler.ELog.Start("engine.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	crawler.ELog.SetFlags(0)
	defer crawler.ELog.Stop()

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
	fmt.Printf("redirectURL: %q\n", redirectURL)
	fmt.Printf("accountID : %q\n", accountID)
	fmt.Println("=========================================")
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Register 註冊可用 stocker
func Register(s stocker.Stocker) {
	stockers = append(stockers, s)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	template, err := parseTemplate("templates/index.html", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(template))
}

func makeQuoteParseFunc(f func(string) (float64, error)) func(crawler.Request) (crawler.ParseResult, error) {
	return func(req crawler.Request) (crawler.ParseResult, error) {
		parseResult := crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{},
			Done:     false,
		}

		symbol := req.Item.(string)

		price, err := f(symbol)
		if err != nil {
			if _, ok := err.(stocker.ErrorNoSupport); !ok {
				parseResult.Requests = append(parseResult.Requests, req)
			} else {
				parseResult.Done = true
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

func makePriceHistoryParseFunc(f func(string) ([]*stocker.DatePrice, error)) func(crawler.Request) (crawler.ParseResult, error) {
	return func(req crawler.Request) (crawler.ParseResult, error) {
		parseResult := crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{},
			Done:     false,
		}

		symbol := req.Item.(string)

		history, err := f(symbol)
		if err != nil {
			if _, ok := err.(stocker.ErrorNoSupport); !ok {
				parseResult.Requests = append(parseResult.Requests, req)
			} else {
				parseResult.Done = true
			}
			return parseResult, err
		}

		rsp := Response{
			symbol: symbol,
			item:   history,
		}
		parseResult.Item = rsp
		parseResult.Done = true

		return parseResult, nil
	}
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

	his := newRecord()
	e := crawler.New(ctx, 10, his.isProcessedOrAdd, his.isDone)

	requests := []crawler.Request{}
	for _, symbol := range symbols {
		for _, stk := range stockers {
			var parseFunc func(crawler.Request) (crawler.ParseResult, error)
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

			requests = append(requests, crawler.Request{
				Item:      symbol,
				ParseFunc: parseFunc,
			})
		}
	}

	prices := make(map[string]interface{}, len(symbols))
	rspChan := e.Run(requests...)
	for rsp := range rspChan {
		result := rsp.(Response)

		prices[result.symbol] = result.item
		his.done(result.symbol)
	}

	err := json.NewEncoder(w).Encode(prices)
	if err != nil {
		log.Printf("json.NewEncoder error: %s\n", err)
	}
}
