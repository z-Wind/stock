package gotd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/z-Wind/stock/api/gotd/api"
	"github.com/z-Wind/stock/api/gotd/api/types"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	AccessTokenExpiresIn  = 1800 * time.Second
	RefreshTokenExpiresIn = 7776000 * time.Second
)

type TDAmeritrade struct {
	apiKey            string
	certPath, keyPath string
	qps               int

	code        string
	redirectURI string
	Token       // 不用 pointer 以便自動初始化
}

type Token struct {
	AccessToken       string
	AccessTokenStart  time.Time
	RefreshToken      string
	RefreshTokenStart time.Time
}

// NewTD default td
func NewTD(apiKey, redirectURL, authDir string) (*TDAmeritrade, error) {
	td := &TDAmeritrade{
		apiKey:   apiKey,
		certPath: filepath.Join(authDir, "cert.pem"),
		keyPath:  filepath.Join(authDir, "key.pem"),
	}
	err := td.LoadAuth(authDir)
	if err != nil {
		log.Printf("Load TD Info Fail: %T %v\n", errors.Cause(err), err)
	}

	err = td.RefreshAccessToken()
	switch e := errors.Cause(err).(type) {
	case *api.ErrorAuth, *api.ErrorGrant:
		err = td.Authentication(redirectURL)
		log.Printf("Not Authorized Update RefreshToken, Error:%s\n", e)
		if err != nil {
			return nil, errors.WithMessage(err, "Authentication")
		}

		err = td.SaveAuth(authDir)
		if err != nil {
			return nil, errors.WithMessage(err, "SaveAuth")
		}
	}

	return td, nil
}

func (td *TDAmeritrade) openAuthServer() error {
	var wg sync.WaitGroup

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		code, ok := r.URL.Query()["code"]
		if !ok {
			log.Println("No code")
			return
		}
		fmt.Fprintf(w, "Get code\n%s", code)
		td.code = code[0]
	})
	srv := &http.Server{Addr: "localhost:8080", Handler: mux}

	wg.Add(1)
	go func() {
		certPath, _ := filepath.Abs(td.certPath)
		keyPath, _ := filepath.Abs(td.keyPath)

		// openssl genrsa -out key.pem 2048
		// openssl req -new -x509 -key key.pem -out cert.pem -days 3650
		if err := srv.ListenAndServeTLS(td.certPath, td.keyPath); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("srv.ListenAndServeTLS error: %T %+v\n", errors.Cause(err), err)
			log.Printf("certPath: %s\n", certPath)
			log.Printf("keyPath: %s\n", keyPath)
		}
		log.Println("Server for code Stop")
	}()

	wg.Wait()

	if err := srv.Close(); err != nil {
		return errors.Wrap(err, "srv.Close")
	}

	log.Printf("sever get code. exiting")
	return nil
}

// Authentication 初步認證
func (td *TDAmeritrade) Authentication(url string) error {
	authURL := api.Authentication(td.apiKey, url)
	err := api.OpenURI(authURL)
	if err != nil {
		return errors.WithMessage(err, "api.Authentication")
	}
	td.redirectURI = url

	err = td.openAuthServer()
	if err != nil {
		return errors.WithMessage(err, "openAuthServer")
	}

	err = td.authorizationToken()
	if err != nil {
		return errors.WithMessage(err, "authorizationToken")
	}

	return nil
}

func (td *TDAmeritrade) authorizationToken() error {
	config := &api.TokenConfig{}
	config.SetAuthorization(td.code, td.redirectURI)
	log.Printf("Get code: %s", td.code)

	token, err := api.PostAccessToken(td.apiKey, "", config)
	if err != nil {
		return errors.WithMessage(err, "api.PostAccessToken")
	}
	td.AccessToken = token.AccessToken
	td.AccessTokenStart = time.Now()
	td.RefreshToken = token.RefreshToken
	td.RefreshTokenStart = time.Now()

	log.Printf("Get all token")
	return nil
}

// SaveAuth 儲存認證
func (td *TDAmeritrade) SaveAuth(dirPath string) error {
	b, err := json.MarshalIndent(td.Token, "", "    ")
	if err != nil {
		return errors.Wrap(err, "json.MarshalIndent")
	}

	path := filepath.Join(dirPath, "secrets.properties")

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "os.Create")
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return errors.Wrap(err, "f.Write")
	}

	return nil
}

// LoadAuth 讀取認證
func (td *TDAmeritrade) LoadAuth(dirPath string) error {
	path := filepath.Join(dirPath, "secrets.properties")
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	return json.Unmarshal(b, td)
}

func (td *TDAmeritrade) IsAccessTokenFail() bool {
	if td.AccessToken == "" {
		return true
	}
	_, err := td.GetQuote("VTI")
	if _, ok := errors.Cause(err).(*api.ErrorAuth); ok {
		return true
	}
	return time.Since(td.AccessTokenStart) > AccessTokenExpiresIn
}

func (td *TDAmeritrade) IsRefreshTokenExpired() bool {
	if td.RefreshToken == "" {
		return true
	}
	return time.Since(td.RefreshTokenStart) > RefreshTokenExpiresIn
}

func (td *TDAmeritrade) RefreshAccessTokenOrNot() error {
	if !td.IsAccessTokenFail() {
		return nil
	}

	err := td.RefreshAccessToken()
	if err != nil {
		return errors.WithMessage(err, "RefreshAccessToken")
	}

	return nil
}

func (td *TDAmeritrade) RefreshAccessToken() error {
	config := &api.TokenConfig{}
	config.SetRefresh(td.RefreshToken)

	token, err := api.PostAccessToken(td.apiKey, "", config)
	if err != nil {
		return errors.WithMessage(err, "api.PostAccessToken")
	}
	td.AccessToken = token.AccessToken
	td.AccessTokenStart = time.Now()

	log.Printf("Refresh access token")

	return nil
}

func (td *TDAmeritrade) GetAccountIDs() ([]int64, error) {
	info, err := td.GetUserPrincipals([]string{})
	if err != nil {
		return nil, err
	}
	accountIDs := []int64{}
	for _, a := range info.Accounts {
		accountIDs = append(accountIDs, a.AccountID)
	}

	return accountIDs, nil
}

func (td *TDAmeritrade) GetQuotes(symbols []string) (*types.Quotes, error) {
	return api.GetQuotes(td.apiKey, td.AccessToken, symbols)
}

func (td *TDAmeritrade) GetQuote(symbol string) (*types.Quote, error) {
	return api.GetQuote(td.apiKey, td.AccessToken, symbol)
}

func (td *TDAmeritrade) GetAccounts(fields []string) (*types.Accounts, error) {
	return api.GetAccounts(td.AccessToken, fields)
}

func (td *TDAmeritrade) GetAccount(accountID int64, fields []string) (*types.Account, error) {
	return api.GetAccount(td.AccessToken, accountID, fields)
}

func (td *TDAmeritrade) SearchInstruments(symbol string, projection string) (*types.Instruments, error) {
	return api.SearchInstruments(td.apiKey, td.AccessToken, symbol, projection)
}

func (td *TDAmeritrade) GetInstrument(cusip string) (*types.Instrument, error) {
	return api.GetInstrument(td.apiKey, td.AccessToken, cusip)
}

func (td *TDAmeritrade) GetHoursforMultipleMarkets(markets []string, date time.Time) (*types.MarketHours, error) {
	return api.GetHoursforMultipleMarkets(td.apiKey, td.AccessToken, markets, date)
}

func (td *TDAmeritrade) GetHoursforSingleMarket(market string, date time.Time) (*types.MarketHours, error) {
	return api.GetHoursforSingleMarket(td.apiKey, td.AccessToken, market, date)
}

func (td *TDAmeritrade) GetMovers(index, direction, change string) (*types.Movers, error) {
	return api.GetMovers(td.apiKey, td.AccessToken, index, direction, change)
}

func (td *TDAmeritrade) GetOptionChain(symbol string) (*types.OptionChain, error) {
	config := &api.OptionChainConfig{}
	config.SetSingle()
	return api.GetOptionChain(td.apiKey, td.AccessToken, symbol, config)
}

func (td *TDAmeritrade) GetPriceHistory(symbol string) (*types.CandleList, error) {
	config := &api.PriceHistoryConfig{}
	config.SetPeriod(api.PriceHistoryPeriodTypeYear, 5)
	config.SetFrequency(api.PriceHistoryFrequencyTypeDaily, 1)
	return api.GetPriceHistory(td.apiKey, td.AccessToken, symbol, config)
}

func (td *TDAmeritrade) GetTransactions(accountID int64) (*types.Transactions, error) {
	config := &api.TransactionsConfig{}
	config.SetDate(time.Now().AddDate(0, -2, -20), time.Now())
	return api.GetTransactions(td.AccessToken, accountID, config)
}

func (td *TDAmeritrade) GetTransaction(accountID, transactionID int64) (*types.Transaction, error) {
	return api.GetTransaction(td.AccessToken, accountID, transactionID)
}

func (td *TDAmeritrade) GetPreferences(accountID int64) (*types.Preferences, error) {
	return api.GetPreferences(td.AccessToken, accountID)
}

func (td *TDAmeritrade) GetStreamerSubscriptionKeys(accountIDs []int64) (*types.SubscriptionKey, error) {
	return api.GetStreamerSubscriptionKeys(td.AccessToken, accountIDs)
}

func (td *TDAmeritrade) GetUserPrincipals(fields []string) (*types.UserPrincipal, error) {
	return api.GetUserPrincipals(td.AccessToken, fields)
}

func (td *TDAmeritrade) UpdatePreferences(accountID int64, preferences *types.Preferences) error {
	return api.UpdatePreferences(td.AccessToken, accountID, preferences)
}

func (td *TDAmeritrade) GetWatchlistsforMultipleAccounts() (*types.Watchlists, error) {
	return api.GetWatchlistsforMultipleAccounts(td.AccessToken)
}

func (td *TDAmeritrade) GetWatchlistsforSingleAccount(accountID int64) (*types.Watchlists, error) {
	return api.GetWatchlistsforSingleAccount(td.AccessToken, accountID)
}

func (td *TDAmeritrade) GetWatchlist(accountID int64, watchlistID string) (*types.Watchlist, error) {
	return api.GetWatchlist(td.AccessToken, accountID, watchlistID)
}

func (td *TDAmeritrade) CreateWatchlist(accountID int64, w *types.WatchlistBasic) error {
	return api.CreateWatchlist(td.AccessToken, accountID, w)
}

func (td *TDAmeritrade) DeleteWatchlist(accountID int64, watchlistID string) error {
	return api.DeleteWatchlist(td.AccessToken, accountID, watchlistID)
}

func (td *TDAmeritrade) ReplaceWatchlist(accountID int64, watchlistID string, w *types.WatchlistBasic) error {
	return api.ReplaceWatchlist(td.AccessToken, accountID, watchlistID, w)
}

func (td *TDAmeritrade) UpdateWatchlist(accountID int64, watchlistID string, w *types.WatchlistBasic) error {
	return api.UpdateWatchlist(td.AccessToken, accountID, watchlistID, w)
}

func (td *TDAmeritrade) GetOrdersByPath(accountID int64) (*types.Orders, error) {
	config := &api.OrdersConfig{}
	return api.GetOrdersByPath(td.AccessToken, accountID, config)
}

func (td *TDAmeritrade) GetOrdersByQuery(accountID int64, maxResults int64, fromEnteredTime, toEnteredTime time.Time, status string) (*types.Orders, error) {
	config := &api.OrdersConfig{}
	return api.GetOrdersByQuery(td.AccessToken, accountID, config)
}

func (td *TDAmeritrade) GetOrder(accountID, orderID int64) (*types.Order, error) {
	return api.GetOrder(td.AccessToken, accountID, orderID)
}

func (td *TDAmeritrade) PlaceOrder(accountID int64, order *types.Order) error {
	return api.PlaceOrder(td.AccessToken, accountID, order)
}

func (td *TDAmeritrade) ReplaceOrder(accountID, orderID int64, order *types.Order) error {
	return api.ReplaceOrder(td.AccessToken, accountID, orderID, order)
}

func (td *TDAmeritrade) CancelOrder(accountID, orderID int64) error {
	return api.CancelOrder(td.AccessToken, accountID, orderID)
}

func (td *TDAmeritrade) GetSavedOrdersbyPath(accountID int64) (*types.SavedOrders, error) {
	return api.GetSavedOrdersbyPath(td.AccessToken, accountID)
}

func (td *TDAmeritrade) GetSavedOrder(accountID, savedOrderID int64) (*types.SavedOrder, error) {
	return api.GetSavedOrder(td.AccessToken, accountID, savedOrderID)
}

func (td *TDAmeritrade) CreateSavedOrder(accountID int64, order *types.Order) error {
	return api.CreateSavedOrder(td.AccessToken, accountID, order)
}

func (td *TDAmeritrade) DeleteSavedOrder(accountID, savedOrderID int64) error {
	return api.DeleteSavedOrder(td.AccessToken, accountID, savedOrderID)
}

func (td *TDAmeritrade) ReplaceSavedOrder(accountID int64, order *types.SavedOrder) error {
	return api.ReplaceSavedOrder(td.AccessToken, accountID, order)
}
