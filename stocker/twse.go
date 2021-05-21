package stocker

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/twse"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"golang.org/x/time/rate"
)

// TWSE stocker
type TWSE struct {
	Service *twse.Service

	limit            *rate.Limiter
	symbolsTWSE      map[string]string
	symbolsTWSE_Path string
	symbolsTPEx      map[string]string
	symbolsTPEx_Path string
}

const (
	TWSE_url = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2"
	TPEx_url = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
)

// NewTWSE 建立 twse Service
func NewTWSE(csvFolderPath string) (*TWSE, error) {
	client := twse.GetClient()
	twse, err := twse.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewTWSE")
	}

	return &TWSE{
		Service:          twse,
		limit:            rate.NewLimiter(rate.Every(time.Second*5/3), 1),
		symbolsTWSE_Path: filepath.Join(csvFolderPath, "TWSE.csv"),
		symbolsTPEx_Path: filepath.Join(csvFolderPath, "TPEx.csv"),
	}, nil
}

func (t *TWSE) getSymbolList(url string) (map[string]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http.Get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	// goquery 限定需為 UTF-8
	b, _ = DecodeBig5(b)
	r := bytes.NewReader(b)

	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrapf(err, "goquery.NewDocumentFromReader")
	}

	symbols := make(map[string]string)
	// Find the link items
	dom.Find("tr td").Not(":only-child").Filter(":first-of-type").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		text := s.Text()
		arr := strings.SplitN(text, "\u3000", 2)
		if len(arr) == 2 {
			symbols[arr[0]] = arr[1]
		}
	})

	return symbols, nil
}

func (t *TWSE) loadData(filepath, url string) (map[string]string, error) {
	symbols := make(map[string]string)
	var err error

	if _, err = os.Stat(filepath); err == nil {
		symbols, err = readFromCSV(filepath)
		if err != nil {
			return nil, errors.Wrapf(err, "readFromCSV")
		}

		return symbols, nil

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		symbols, err = t.getSymbolList(url)
		if err != nil {
			return nil, errors.Wrapf(err, "t.getSymbolList")
		}

		err = writeToCSV(symbols, filepath)
		if err != nil {
			return nil, errors.Wrapf(err, "writeToCSV")
		}

		return symbols, nil

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		symbols, err = t.getSymbolList(url)
		if err != nil {
			return nil, errors.Wrapf(err, "t.getSymbolList")
		}

		err = writeToCSV(symbols, filepath)
		if err != nil {
			return nil, errors.Wrapf(err, "writeToCSV")
		}

		return symbols, nil

	}
}

func (t *TWSE) isInTPEx(symbol string) bool {
	symbol = strings.SplitN(symbol, ".", 2)[0]
	if t.symbolsTPEx != nil {
		_, ok := t.symbolsTPEx[symbol]
		return ok
	}

	var err error
	t.symbolsTPEx, err = t.loadData(t.symbolsTPEx_Path, TPEx_url)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "t.loadData"))
		return false
	}

	_, ok := t.symbolsTPEx[symbol]
	return ok
}

func (t *TWSE) isInTWSE(symbol string) bool {
	symbol = strings.SplitN(symbol, ".", 2)[0]
	if t.symbolsTWSE != nil {
		_, ok := t.symbolsTWSE[symbol]
		return ok
	}

	var err error
	t.symbolsTWSE, err = t.loadData(t.symbolsTWSE_Path, TWSE_url)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "t.loadData"))
		return false
	}

	_, ok := t.symbolsTWSE[symbol]
	return ok
}

// Quote 得到股票價格
func (t *TWSE) Quote(ctx context.Context, symbol string) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ErrorFatal{errors.Wrapf(ctx.Err(), "twse").Error()}
	default:
	}

	var quote *twse.StockInfo
	var err error
	switch {
	case t.isInTWSE(symbol):
		if err := t.limit.Wait(ctx); err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
		}

		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoTWSE(id)
		call.Context(ctx)
		quote, err = call.Do()
		if err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: GetStockInfoTWSE.Do").Error()}
		}
	case t.isInTPEx(symbol):
		if err := t.limit.Wait(ctx); err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
		}

		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoOTC(id)
		call.Context(ctx)
		quote, err = call.Do()
		if err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: GetStockInfoOTC.Do").Error()}
		}
	default:
		return 0, ErrorNoSupport{fmt.Sprintf("twse: Quote: %s is not supported", symbol)}
	}

	price := quote.TradePrice
	if quote.TradePrice == 0 {
		price = quote.YesterdayPrice
	}
	return float64(price), nil
}

func (t *TWSE) priceHistoryTWSE(ctx context.Context, symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyTWSE(id, time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		return nil, ErrorFatal{errors.Wrapf(err, "twse: MonthlyTWSE.Do").Error()}
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   float64(trade.Open),
			High:   float64(trade.High),
			Low:    float64(trade.Low),
			Close:  float64(trade.Close),
			Volume: float64(trade.Volume),
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"twse: Empty List"}
	}

	return timeSeries, nil
}
func (t *TWSE) priceHistoryOTC(ctx context.Context, symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyOTC(id, time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		return nil, ErrorFatal{errors.Wrapf(err, "twse: MonthlyOTC.Do").Error()}
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   float64(trade.Open),
			High:   float64(trade.High),
			Low:    float64(trade.Low),
			Close:  float64(trade.Close),
			Volume: float64(trade.Volume),
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"twse: Empty List"}
	}

	return timeSeries, nil
}

// PriceHistory 得到股票歷史價格
func (t *TWSE) PriceHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	select {
	case <-ctx.Done():
		return nil, ErrorFatal{errors.Wrapf(ctx.Err(), "twse").Error()}
	default:
	}

	switch {
	case t.isInTWSE(symbol):
		if err := t.limit.Wait(ctx); err != nil {
			return nil, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
		}

		return t.priceHistoryTWSE(ctx, symbol)
	case t.isInTPEx(symbol):
		if err := t.limit.Wait(ctx); err != nil {
			return nil, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
		}

		return t.priceHistoryOTC(ctx, symbol)
	default:
		return nil, ErrorNoSupport{fmt.Sprintf("twse: PriceHistory: %s is not supported", symbol)}
	}
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (t *TWSE) PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	return nil, ErrorNoSupport{"TWSE does not support PriceAdjHistory"}
}

//DecodeBig5 convert BIG5 to UTF-8
func DecodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	b, err := ioutil.ReadAll(O)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadAll")
	}
	return b, nil
}

func writeToCSV(dict map[string]string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return errors.Wrapf(err, "os.Create")
	}
	defer file.Close()

	w := csv.NewWriter(file)

	for key, val := range dict {
		if err := w.Write([]string{key, val}); err != nil {
			return errors.Wrapf(err, "w.Write")
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return errors.Wrapf(err, "w.Error")
	}

	return nil
}

func readFromCSV(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open")
	}
	defer file.Close()

	r := csv.NewReader(file)

	dict := make(map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "r.Read")
		}

		dict[record[0]] = record[1]
	}

	return dict, nil
}
