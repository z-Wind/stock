package main

import (
	"bytes"
	"fmt"
	"github.com/z-Wind/stock/api/alphavantage"
	"github.com/z-Wind/stock/api/gotd"
	"github.com/z-Wind/stock/api/twse"
	"github.com/z-Wind/stock/instance"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

var (
	buildstamp = ""
	githash    = ""
	goversion  = ""

	td         *gotd.TDAmeritrade
	alpha      *alphavantage.AlphaVantage
	tw         *twse.TWSE
	accountIDs []int64

	quote           Quote
	priceHistory    PriceHistory
	priceAdjHistory PriceAdjHistory
)

func getCurExePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", errors.Wrap(err, "exec.LookPath")
	}

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrap(err, "filepath.Abs")
	}

	rst := filepath.Dir(path)

	return rst, nil
}

func getCurScriptPath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller Fail")
	}

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrap(err, "filepath.Abs")
	}

	rst := filepath.Dir(path)

	return rst, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func init() {
	tdAPIKey := instance.TdAPIKey
	tdURL := instance.TdURL
	alphaVantageKey := instance.AlphaVantageKey

	var err error
	path, err := getCurScriptPath()
	if err != nil {
		path = "./instance"
	} else {
		path = filepath.Join(path, "instance")
	}
	log.Printf("Current Path:%s", path)

	//tdAPIKey & tdURL should be added by yourself
	td, err = gotd.NewTD(tdAPIKey, tdURL, path)
	if err != nil {
		log.Printf("gotdNewTD error: %s\n", err)
		td = nil
	} else {
		quote.RegisterQuote(quoteTD)
		priceHistory.RegisterPriceHistory(priceHistoryTD)
		accountIDs, err = td.GetAccountIDs()
		if err != nil {
			log.Printf("td.GetAccountIDs Error:%s", err)
		}
	}

	alpha, err = alphavantage.NewAlphaVantage(alphaVantageKey)
	if err != nil {
		log.Printf("alphavantage.NewAlphaVantage error: %s\n", err)
		alpha = nil
	} else {
		quote.RegisterQuote(quoteAlpha)
		priceHistory.RegisterPriceHistory(priceHistoryAlpha)
		priceAdjHistory.RegisterPriceAdjHistory(priceAdjHistoryAlpha)
	}

	tw, err = twse.NewTWSE()
	if err != nil {
		log.Printf("twse.NewTWSE error: %s\n", err)
		tw = nil
	} else {
		quote.RegisterQuote(quoteTWSE)
	}
}

func main() {
	fmt.Println("=========================================")
	fmt.Printf("Git Commit Hash: %s\n", githash)
	fmt.Printf("Build Time : %s\n", buildstamp)
	fmt.Printf("Golang Version : %s\n", goversion)
	fmt.Println("=========================================")

	http.Handle("/", http.HandlerFunc(handleIndex))
	http.Handle("/quote", http.HandlerFunc(quote.HandlerFunc))
	http.Handle("/savedOrder", http.HandlerFunc(savedOrder))
	http.Handle("/priceHistory", http.HandlerFunc(priceHistory.HandlerFunc))
	http.Handle("/priceAdjHistory", http.HandlerFunc(priceAdjHistory.HandlerFunc))

	log.Fatal(http.ListenAndServe("localhost:6060", nil))
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

func parseTemplate(fileName string, data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	template, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	err = template.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]struct{}{}
	result := []string{}

	for v := range elements {
		if _, ok := encountered[elements[v]]; ok {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = struct{}{}
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
