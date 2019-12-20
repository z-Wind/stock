package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"reflect"
	"runtime"
	"github.com/z-Wind/stock/api/gotd/api/types"
	"github.com/z-Wind/stock/api/setting"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

type newRequest func(string, interface{}) (*http.Request, error)

func OpenURI(uri string) error {
	var run []string

	switch runtime.GOOS {
	case "windows":
		run = []string{"cmd", "/c", "start"}
		uri = strings.Replace(uri, "&", "^&", -1)
	case "linux":
		run = []string{"xdg-open"}
	//"darwin":  []string{"open"},
	default:
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	run = append(run, uri)
	cmd := exec.Command(run[0], run[1:]...)
	return cmd.Start()
}

// Authentication https://developer.tdameritrade.com/content/simple-auth-local-apps
func Authentication(apiKey, redirectURI string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "auth.tdameritrade.com",
		Path:   "auth",
	}
	paras := u.Query()
	paras.Add("response_type", "code")
	paras.Add("redirect_uri", redirectURI)
	paras.Add("client_id", apiKey)

	u.RawQuery = paras.Encode()

	return u.String()
}

// PostAccessToken https://developer.tdameritrade.com/authentication/apis/post/token-0
func PostAccessToken(apiKey, accessToken string, config *TokenConfig) (*types.EASObject, error) {
	reqURL := "https://api.tdameritrade.com/v1/oauth2/token"

	dataSaved := new(types.EASObject)
	config.Client_id = apiKey
	paras := Struct2URLValue(*config)
	err := api(reqURL, accessToken, newReqPost, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetQuotes https://developer.tdameritrade.com/quotes/apis/get/marketdata/quotes
func GetQuotes(apiKey, accessToken string, symbols []string) (*types.Quotes, error) {
	reqURL := "https://api.tdameritrade.com/v1/marketdata/quotes"

	dataSaved := new(types.Quotes)
	paras := map[string]string{
		"apikey": apiKey,
		"symbol": strings.Join(symbols, ","),
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetQuote https://developer.tdameritrade.com/quotes/apis/get/marketdata/%7Bsymbol%7D/quotes
func GetQuote(apiKey, accessToken, symbol string) (*types.Quote, error) {
	symbol = strings.ToUpper(symbol)
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/quotes", symbol)

	dataSaved := new(types.Quotes)
	paras := map[string]string{
		"apikey": apiKey,
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	m, ok := (*dataSaved)[symbol]
	if !ok {
		return nil, fmt.Errorf("GetQuote:%s not Found", symbol)
	}

	return m, nil
}

// GetAccounts https://developer.tdameritrade.com/account-access/apis/get/accounts-0
func GetAccounts(accessToken string, fields []string) (*types.Accounts, error) {
	reqURL := "https://api.tdameritrade.com/v1/accounts"

	dataSaved := new(types.Accounts)
	paras := map[string]string{
		"fields": strings.Join(fields, ","), //positions,orders
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetAccount https://developer.tdameritrade.com/account-access/apis/get/accounts/%7BaccountId%7D-0
func GetAccount(accessToken string, accountID int64, fields []string) (*types.Account, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d", accountID)

	dataSaved := new(types.Account)
	paras := map[string]string{
		"fields": strings.Join(fields, ","), //positions,orders
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// SearchInstruments https://developer.tdameritrade.com/instruments/apis/get/instruments
func SearchInstruments(apiKey, accessToken, symbol, projection string) (*types.Instruments, error) {
	reqURL := "https://api.tdameritrade.com/v1/instruments"

	dataSaved := new(types.Instruments)
	paras := map[string]string{
		"apikey":     apiKey,
		"symbol":     symbol,
		"projection": projection, //symbol-search, symbol-regex, desc-search, desc-regex, fundamental
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetInstrument https://developer.tdameritrade.com/instruments/apis/get/instruments/%7Bcusip%7D
func GetInstrument(apiKey, accessToken, cusip string) (*types.Instrument, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/instruments/%s", cusip)

	dataSaved := new([]*types.Instrument)
	paras := map[string]string{
		"apikey": apiKey,
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	if len(*dataSaved) == 0 {
		return nil, fmt.Errorf("GetInstrument(%s,%s,%s) not found", apiKey, accessToken, cusip)
	}
	return (*dataSaved)[0], nil
}

// GetHoursforMultipleMarkets https://developer.tdameritrade.com/market-hours/apis/get/marketdata/hours
func GetHoursforMultipleMarkets(apiKey, accessToken string, markets []string, date time.Time) (*types.MarketHours, error) {
	reqURL := "https://api.tdameritrade.com/v1/marketdata/hours"

	dataSaved := new(types.MarketHours)
	paras := map[string]string{
		"apikey":  apiKey,
		"markets": strings.Join(markets, ","), //EQUITY, OPTION, FUTURE, BOND, or FOREX
		"date":    time.Now().UTC().Format("2006-01-02"),
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetHoursforSingleMarket https://developer.tdameritrade.com/market-hours/apis/get/marketdata/%7Bmarket%7D/hours
func GetHoursforSingleMarket(apiKey, accessToken string, market string, date time.Time) (*types.MarketHours, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/hours", market)

	dataSaved := new(types.MarketHours)
	paras := map[string]string{
		"apikey": apiKey,
		"date":   date.UTC().Format("2006-01-02"),
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetMovers https://developer.tdameritrade.com/movers/apis/get/marketdata/%7Bindex%7D/movers
func GetMovers(apiKey, accessToken string, index, direction, change string) (*types.Movers, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/movers", index)

	dataSaved := new(types.Movers)
	paras := map[string]string{
		"apikey":    apiKey,
		"direction": direction, // up, down
		"change":    change,    // percent, value
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetOptionChain https://developer.tdameritrade.com/option-chains/apis/get/marketdata/chains
func GetOptionChain(apiKey, accessToken, symbol string, config *OptionChainConfig) (*types.OptionChain, error) {
	reqURL := "https://api.tdameritrade.com/v1/marketdata/chains"

	dataSaved := new(types.OptionChain)
	config.Apikey = apiKey
	config.Symbol = symbol
	paras := Struct2URLValue(*config)
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetPriceHistory https://developer.tdameritrade.com/price-history/apis/get/marketdata/%7Bsymbol%7D/pricehistory
func GetPriceHistory(apiKey, accessToken, symbol string, config *PriceHistoryConfig) (*types.CandleList, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/pricehistory", symbol)

	dataSaved := new(types.CandleList)
	config.Apikey = apiKey
	paras := Struct2URLValue(*config)
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetTransactions https://developer.tdameritrade.com/transaction-history/apis/get/accounts/%7BaccountId%7D/transactions-0
func GetTransactions(accessToken string, accountID int64, config *TransactionsConfig) (*types.Transactions, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/transactions", accountID)

	dataSaved := new(types.Transactions)
	paras := Struct2URLValue(*config)
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetTransaction https://developer.tdameritrade.com/transaction-history/apis/get/accounts/%7BaccountId%7D/transactions/%7BtransactionId%7D-0
// 無法獲得，猜測未完成
func GetTransaction(accessToken string, accountID, transactionID int64) (*types.Transaction, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/transactions/%d", accountID, transactionID)

	dataSaved := new(types.Transaction)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetPreferences https://developer.tdameritrade.com/user-principal/apis/get/accounts/%7BaccountId%7D/preferences-0
func GetPreferences(accessToken string, accountID int64) (*types.Preferences, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/preferences", accountID)

	dataSaved := new(types.Preferences)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetStreamerSubscriptionKeys https://developer.tdameritrade.com/user-principal/apis/get/userprincipals/streamersubscriptionkeys-0
func GetStreamerSubscriptionKeys(accessToken string, accountIDs []int64) (*types.SubscriptionKey, error) {
	reqURL := "https://api.tdameritrade.com/v1/userprincipals/streamersubscriptionkeys"

	dataSaved := new(types.SubscriptionKey)
	paras := map[string]string{
		"accountIds": strings.Trim(strings.Replace(fmt.Sprint(accountIDs), " ", ",", -1), "[]"),
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetUserPrincipals https://developer.tdameritrade.com/user-principal/apis/get/userprincipals-0
func GetUserPrincipals(accessToken string, fields []string) (*types.UserPrincipal, error) {
	reqURL := "https://api.tdameritrade.com/v1/userprincipals"

	dataSaved := new(types.UserPrincipal)
	paras := map[string]string{
		"fields": strings.Join(fields, ","),
	}
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// UpdatePreferences https://developer.tdameritrade.com/user-principal/apis/put/accounts/%7BaccountId%7D/preferences-0
func UpdatePreferences(accessToken string, accountID int64, preferences *types.Preferences) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/preferences", accountID)

	content, err := json.Marshal(preferences)
	if err != nil {
		return errors.Wrap(err, "UpdatePreferences")
	}

	return api(reqURL, accessToken, newReqPut, nil, content)
}

// GetWatchlistsforMultipleAccounts https://developer.tdameritrade.com/watchlist/apis/get/accounts/watchlists-0
func GetWatchlistsforMultipleAccounts(accessToken string) (*types.Watchlists, error) {
	reqURL := "https://api.tdameritrade.com/v1/accounts/watchlists"

	dataSaved := new(types.Watchlists)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetWatchlistsforSingleAccount https://developer.tdameritrade.com/watchlist/apis/get/accounts/%7BaccountId%7D/watchlists-0
func GetWatchlistsforSingleAccount(accessToken string, accountID int64) (*types.Watchlists, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists", accountID)

	dataSaved := new(types.Watchlists)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetWatchlist https://developer.tdameritrade.com/watchlist/apis/get/accounts/%7BaccountId%7D/watchlists/%7BwatchlistId%7D-0
func GetWatchlist(accessToken string, accountID int64, watchlistID string) (*types.Watchlist, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists/%s", accountID, watchlistID)

	dataSaved := new(types.Watchlist)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// CreateWatchlist https://developer.tdameritrade.com/watchlist/apis/post/accounts/%7BaccountId%7D/watchlists-0
func CreateWatchlist(accessToken string, accountID int64, w *types.WatchlistBasic) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists", accountID)

	content, err := json.Marshal(w)
	if err != nil {
		return errors.Wrap(err, "CreateWatchlist")
	}

	return api(reqURL, accessToken, newReqPost, nil, content)
}

// DeleteWatchlist https://developer.tdameritrade.com/watchlist/apis/delete/accounts/%7BaccountId%7D/watchlists/%7BwatchlistId%7D-0
func DeleteWatchlist(accessToken string, accountID int64, watchlistID string) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists/%s", accountID, watchlistID)

	return api(reqURL, accessToken, newReqDelete, nil, nil)
}

// ReplaceWatchlist https://developer.tdameritrade.com/watchlist/apis/put/accounts/%7BaccountId%7D/watchlists/%7BwatchlistId%7D-0
func ReplaceWatchlist(accessToken string, accountID int64, watchlistID string, w *types.WatchlistBasic) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists/%s", accountID, watchlistID)

	content, err := json.Marshal(w)
	if err != nil {
		return errors.Wrap(err, "ReplaceWatchlist")
	}

	return api(reqURL, accessToken, newReqPut, nil, content)
}

// UpdateWatchlist https://developer.tdameritrade.com/watchlist/apis/patch/accounts/%7BaccountId%7D/watchlists/%7BwatchlistId%7D-0
func UpdateWatchlist(accessToken string, accountID int64, watchlistID string, w *types.WatchlistBasic) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/watchlists/%s", accountID, watchlistID)

	content, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		return errors.Wrap(err, "UpdateWatchlist")
	}

	return api(reqURL, accessToken, newReqPatch, nil, content)
}

// GetOrdersByPath https://developer.tdameritrade.com/account-access/apis/get/accounts/%7BaccountId%7D/orders-0
func GetOrdersByPath(accessToken string, accountID int64, config *OrdersConfig) (*types.Orders, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/orders", accountID)

	dataSaved := new(types.Orders)
	paras := Struct2URLValue(*config)
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetOrdersByQuery https://developer.tdameritrade.com/account-access/apis/get/orders-0
func GetOrdersByQuery(accessToken string, accountID int64, config *OrdersConfig) (*types.Orders, error) {
	reqURL := "https://api.tdameritrade.com/v1/orders"

	dataSaved := new(types.Orders)
	paras := Struct2URLValue(*config)
	paras.Add("accountId", strconv.FormatInt(accountID, 10))
	err := api(reqURL, accessToken, newReqGet, dataSaved, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetOrder https://developer.tdameritrade.com/account-access/apis/get/accounts/%7BaccountId%7D/orders/%7BorderId%7D-0
func GetOrder(accessToken string, accountID, orderID int64) (*types.Order, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/orders/%d", accountID, orderID)

	dataSaved := new(types.Order)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// PlaceOrder https://developer.tdameritrade.com/account-access/apis/post/accounts/%7BaccountId%7D/orders-0
func PlaceOrder(accessToken string, accountID int64, order *types.Order) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/orders", accountID)

	content, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "PlaceOrder")
	}

	return api(reqURL, accessToken, newReqPost, nil, content)
}

// ReplaceOrder https://developer.tdameritrade.com/account-access/apis/put/accounts/%7BaccountId%7D/orders/%7BorderId%7D-0
func ReplaceOrder(accessToken string, accountID, orderID int64, order *types.Order) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/orders/%d", accountID, orderID)

	content, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "ReplaceOrder")
	}

	return api(reqURL, accessToken, newReqPut, nil, content)
}

// CancelOrder https://developer.tdameritrade.com/account-access/apis/delete/accounts/%7BaccountId%7D/orders/%7BorderId%7D-0
func CancelOrder(accessToken string, accountID, orderID int64) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/orders/%d", accountID, orderID)

	return api(reqURL, accessToken, newReqDelete, nil, nil)
}

// GetSavedOrdersbyPath https://developer.tdameritrade.com/account-access/apis/get/accounts/%7BaccountId%7D/savedorders-0
func GetSavedOrdersbyPath(accessToken string, accountID int64) (*types.SavedOrders, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/savedorders", accountID)

	dataSaved := new(types.SavedOrders)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// GetSavedOrder https://developer.tdameritrade.com/account-access/apis/get/accounts/%7BaccountId%7D/savedorders/%7BsavedOrderId%7D-0
func GetSavedOrder(accessToken string, accountID, savedOrderID int64) (*types.SavedOrder, error) {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/savedorders/%d", accountID, savedOrderID)

	dataSaved := new(types.SavedOrder)
	err := api(reqURL, accessToken, newReqGet, dataSaved, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return dataSaved, nil
}

// CreateSavedOrder https://developer.tdameritrade.com/account-access/apis/post/accounts/%7BaccountId%7D/savedorders-0
func CreateSavedOrder(accessToken string, accountID int64, order *types.Order) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/savedorders", accountID)

	content, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "CreateSavedOrder")
	}

	return api(reqURL, accessToken, newReqPost, nil, content)
}

// DeleteSavedOrder https://developer.tdameritrade.com/account-access/apis/delete/accounts/%7BaccountId%7D/savedorders/%7BsavedOrderId%7D-0
func DeleteSavedOrder(accessToken string, accountID, savedOrderID int64) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/savedorders/%d", accountID, savedOrderID)

	err := api(reqURL, accessToken, newReqDelete, nil, nil)
	if err != nil {
		return errors.WithMessage(err, "api")
	}

	return nil
}

// ReplaceSavedOrder https://developer.tdameritrade.com/account-access/apis/put/accounts/%7BaccountId%7D/savedorders/%7BsavedOrderId%7D-0
func ReplaceSavedOrder(accessToken string, accountID int64, order *types.SavedOrder) error {
	reqURL := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%d/savedorders/%d", accountID, order.SavedOrderID)

	content, err := json.Marshal(order.Order)
	if err != nil {
		return errors.Wrap(err, "ReplaceSavedOrder")
	}

	err = api(reqURL, accessToken, newReqPut, nil, content)
	if err != nil {
		return errors.WithMessage(err, "api")
	}

	return nil
}

func api(reqURL, accessToken string, newReq newRequest, dataSaved, content interface{}) error {
	req, err := newReq(reqURL, content)
	if err != nil {
		return errors.WithMessagef(err, "newReq(%v,%v)", reqURL, content)
	}

	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	body, err := doReq(req)
	if err != nil {
		return errors.WithMessagef(err, "doReq(%v)", req)
	}

	if dataSaved != nil {
		err = json.Unmarshal(body, dataSaved)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal")
		}
	}

	return nil
}

func newReq(reqURL string, content interface{}, method string) (req *http.Request, err error) {
	switch content := content.(type) {
	case url.Values:
		switch method {
		case http.MethodGet:
			req, err = http.NewRequest(method, reqURL, nil)
			req.URL.RawQuery = content.Encode()
		case http.MethodPost:
			req, err = http.NewRequest(method, reqURL, strings.NewReader(content.Encode()))
		default:
			return nil, fmt.Errorf("Not Support Method %s", method)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case map[string]string:
		content = removeEmptyParas(content)
		q := url.Values{}
		for k, v := range content {
			q.Add(k, v)
		}

		switch method {
		case http.MethodGet:
			req, err = http.NewRequest(method, reqURL, nil)
			req.URL.RawQuery = q.Encode()
		case http.MethodPost:
			req, err = http.NewRequest(method, reqURL, strings.NewReader(q.Encode()))
		default:
			return nil, fmt.Errorf("Not Support Method %s", method)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case []byte:
		req, err = http.NewRequest(method, reqURL, bytes.NewBuffer(content))
		req.Header.Set("Content-Type", "application/json")
	case nil:
		req, err = http.NewRequest(method, reqURL, nil)
	default:
		return nil, fmt.Errorf("unexpected type %T: %v", content, content)
	}
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}

	req = setting.WrapRequest(req)

	return req, nil
}

func newReqPost(reqURL string, content interface{}) (req *http.Request, err error) {
	return newReq(reqURL, content, http.MethodPost)
}

func newReqGet(reqURL string, content interface{}) (req *http.Request, err error) {
	return newReq(reqURL, content, http.MethodGet)
}

func newReqPut(reqURL string, body interface{}) (*http.Request, error) {
	return newReq(reqURL, body, http.MethodPut)
}

func newReqDelete(reqURL string, _ interface{}) (*http.Request, error) {
	return newReq(reqURL, nil, http.MethodDelete)
}

func newReqPatch(reqURL string, content interface{}) (*http.Request, error) {
	return newReq(reqURL, content, http.MethodPatch)
}

func doReq(req *http.Request) ([]byte, error) {
	//log.Printf("request URL: %s", req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer resp.Body.Close()
	defer log.Printf("td response %s:%s", resp.Status, req.URL)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s", resp.Status)
	case http.StatusForbidden,
		http.StatusInternalServerError:
		m := new(types.Message)
		err = json.Unmarshal(body, m)
		if err != nil {
			return nil, errors.Wrap(err, "ErrorMessage json.Unmarshal")
		}
		return nil, errors.New(m.ErrorMessage)
	case http.StatusBadRequest:
		m := new(types.Message)
		err = json.Unmarshal(body, m)
		if err != nil {
			return nil, errors.Wrap(err, "ErrorMessage json.Unmarshal")
		}
		switch m.ErrorMessage {
		case "invalid_grant":
			return nil, &ErrorGrant{message: m.ErrorMessage}
		default:
			return nil, errors.New(m.ErrorMessage)
		}
	case http.StatusUnauthorized:
		m := new(types.Message)
		err = json.Unmarshal(body, m)
		if err != nil {
			return nil, errors.Wrap(err, "ErrorMessage json.Unmarshal")
		}
		return nil, &ErrorAuth{message: m.ErrorMessage}
	}

	return body, nil
}

func removeEmptyParas(paras map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range paras {
		if v != "" && v != "0" && v != "<nil>" && v != "nil" {
			m[k] = v
		}
	}
	return m
}

// Lcfirst 首字小寫
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Struct2URLValue struct to url.Values
func Struct2URLValue(obj interface{}) url.Values {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	paras := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		str, err := encode(v.Field(i))
		if err != nil {
			panic(err)
		}
		if str != "" && str != "0" && str != "nil" {
			paras.Add(Lcfirst(t.Field(i).Name), str)
		}
	}
	return paras
}

func encode(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return "nil", nil

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", v.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%v", v.Uint()), nil

	case reflect.String:
		return v.String(), nil

	case reflect.Ptr:
		return encode(v.Elem())

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v.Float()), nil

	case reflect.Bool:
		return fmt.Sprintf("%v", v.Bool()), nil

	default: // complex, map, struct, array, slice, chan, func, interface
		return "", fmt.Errorf("unsupported type: %s", v.Type())
	}
}
