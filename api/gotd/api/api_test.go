package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"github.com/z-Wind/stock/api/gotd/api/types"
	"github.com/z-Wind/stock/instance"
	"testing"
	"time"
)

var (
	apiKey       = instance.TdAPIKey
	redirectURL  = instance.TdURL
	code         = ""
	accessToken  = ""
	refreshToken = ""
)

func init() {
	path := filepath.Join("../../../instance/", "secrets.properties")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var token struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}

	json.Unmarshal(b, &token)
	if apiKey == "" {
		panic("td is null")
	}

	accessToken = token.AccessToken
	refreshToken = token.RefreshToken

}

func getAccountID() int64 {
	info, err := GetUserPrincipals(accessToken, []string{})
	if err != nil {
		panic(fmt.Sprintf("GetUserPrincipals: Error : %s", err))
	}
	accountIDs := []int64{}
	for _, a := range info.Accounts {
		accountIDs = append(accountIDs, a.AccountID)
	}
	accountID := accountIDs[0]
	return accountID
}

func TestAuthentication(t *testing.T) {
	type args struct {
		apiKey string
		url    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"Test", args{apiKey, redirectURL}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Authentication(tt.args.apiKey, tt.args.url)
			req, err := http.Get(got)
			if err != nil || req.StatusCode != 200 {
				t.Errorf("Authentication() error= %s, status %d", err, req.StatusCode)
			}
		})
	}
}

func TestPostAccessToken(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		config      *TokenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		//{"Access", args{apiKey, "", (&TokenConfig{}).SetAuthorization(code, redirectURL)}, "Bearer", false},
		{"Refresh", args{apiKey, "", (&TokenConfig{}).SetRefresh(refreshToken)}, "Bearer", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostAccessToken(tt.args.apiKey, tt.args.accessToken, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got.TokenType != tt.want {
				t.Errorf("PostAccessToken() = %v, want %v", got.TokenType, tt.want)
			}
		})
	}
}

func TestGetQuotes(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		symbols     []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{apiKey, "", []string{"VTI", "VBR"}}, 2, false},
		{"Test", args{apiKey, "", []string{"bnd"}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQuotes(tt.args.apiKey, tt.args.accessToken, tt.args.symbols)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if len(*got) != tt.want {
				t.Errorf("GetQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetQuote(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		symbol      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"ETF", args{apiKey, "", "VTI"}, "VTI", false},
		{"ETF", args{apiKey, "", "bnd"}, "BND", false},
		{"Error", args{apiKey, "", "abc"}, "ABC", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQuote(tt.args.apiKey, tt.args.accessToken, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if err == nil && got.Data.(*types.ETF).Symbol != tt.want {
				t.Errorf("GetQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccounts(t *testing.T) {
	type args struct {
		accessToken string
		fields      []string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Accounts
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, []string{AccountFieldPositions, AccountFieldOrders}}, &types.Accounts{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccounts(tt.args.accessToken, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got == tt.want {
				t.Errorf("GetAccounts() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	type args struct {
		accessToken string
		accountID   int64
		fields      []string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Account
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, getAccountID(), []string{AccountFieldPositions, AccountFieldOrders}}, &types.Account{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccount(tt.args.accessToken, tt.args.accountID, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got == tt.want {
				t.Errorf("GetAccount() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestSearchInstruments(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		symbol      string
		projection  string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Instruments
		wantErr bool
	}{
		// TODO: Add test cases.
		{"symbol-search", args{apiKey, "", "VTI,BND", InstrumentsProjectionSymbolSearch}, &types.Instruments{
			"BND": {
				Cusip:       "921937835",
				Symbol:      "BND",
				Description: "Vanguard Total Bond Market ETF",
				Exchange:    "NASDAQ",
				AssetType:   "ETF",
			},
			"VTI": {
				Cusip:       "922908769",
				Symbol:      "VTI",
				Description: "Vanguard Total Stock Market ETF",
				Exchange:    "Pacific",
				AssetType:   "ETF",
			},
		}, false},
		{"symbol-regex", args{apiKey, "", "VTI", InstrumentsProjectionSymbolRegex}, &types.Instruments{
			"VTI": {
				Cusip:       "922908769",
				Symbol:      "VTI",
				Description: "Vanguard Total Stock Market ETF",
				Exchange:    "Pacific",
				AssetType:   "ETF",
			},
		}, false},
		{"desc-search", args{apiKey, "", "Spirits", InstrumentsProjectionDescSearch}, &types.Instruments{
			"SRSG": {
				Cusip:       "84861Y107",
				Symbol:      "SRSG",
				Description: "Spirits Time International, Inc. Common Stock (PC)",
				Exchange:    "Pink Sheet",
				AssetType:   "EQUITY",
			},
		}, false},
		{"desc-regex", args{apiKey, "", `Vanguard Total Stock.*`, InstrumentsProjectionDescRegex}, &types.Instruments{
			"VTI": {
				Cusip:       "922908769",
				Symbol:      "VTI",
				Description: "Vanguard Total Stock Market ETF",
				Exchange:    "Pacific",
				AssetType:   "ETF",
			},
		}, false},
		// {"fundamental", args{apiKey, "", "VTI", InstrumentsProjectionFundamental}, &types.Instruments{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchInstruments(tt.args.apiKey, tt.args.accessToken, tt.args.symbol, tt.args.projection)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchInstruments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchInstruments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInstrument(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		cusip       string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Instrument
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{apiKey, "", "861025104"}, &types.Instrument{
			Cusip:       "861025104",
			Symbol:      "SYBT",
			Description: "Stock Yards Bancorp, Inc. - Common Stock",
			Exchange:    "NASDAQ",
			AssetType:   "EQUITY",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstrument(tt.args.apiKey, tt.args.accessToken, tt.args.cusip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstrument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstrument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHoursforMultipleMarkets(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		markets     []string
		date        time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{apiKey, "", []string{HoursMarketsBOND, HoursMarketsEQUITY, HoursMarketsFOREX, HoursMarketsFUTURE, HoursMarketsOPTION}, time.Now()}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHoursforMultipleMarkets(tt.args.apiKey, tt.args.accessToken, tt.args.markets, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHoursforMultipleMarkets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if len(*got) != tt.want {
				t.Errorf("GetHoursforMultipleMarkets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHoursforSingleMarket(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		market      string
		date        time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{apiKey, "", HoursMarketsEQUITY, time.Now()}, "equity", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHoursforSingleMarket(tt.args.apiKey, tt.args.accessToken, tt.args.market, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHoursforSingleMarket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Logf("content:%s", string(content))
			if _, ok := (*got)[tt.want]; !ok {
				t.Errorf("GetHoursforSingleMarket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMovers(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		index       string
		direction   string
		change      string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"COMPX", args{apiKey, "", MoversIndexCOMPX, MoversDirectionUp, MoversChangePercent}, 0, false},
		{"DJI", args{apiKey, "", MoversIndexDJI, MoversDirectionDown, MoversChangeValue}, 0, false},
		{"SPX.X", args{apiKey, "", MoversIndexSPXX, MoversDirectionDown, MoversChangeValue}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMovers(tt.args.apiKey, tt.args.accessToken, tt.args.index, tt.args.direction, tt.args.change)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMovers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if len(*got) <= tt.want {
				t.Errorf("GetMovers() = %v, want more than %v", got, tt.want)
			}
		})
	}
}

func TestGetOptionChain(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		symbol      string
		config      *OptionChainConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"VTI", args{apiKey, "", "VTI", (&OptionChainConfig{}).SetSingle()}, "VTI", false},
		{"VTI", args{apiKey, "", "VTI", (&OptionChainConfig{}).SetVertical()}, "VTI", false},
		//{"V00", args{apiKey, "", "V00", &OptionChainConfig{Strategy: OptionChainStrategyVERTICAL}}, "V00", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOptionChain(tt.args.apiKey, tt.args.accessToken, tt.args.symbol, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOptionChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got.Symbol != tt.want {
				t.Errorf("GetOptionChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPriceHistory(t *testing.T) {
	type args struct {
		apiKey      string
		accessToken string
		symbol      string
		config      *PriceHistoryConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Period", args{apiKey, "", "VTI", (&PriceHistoryConfig{}).SetPeriod("day", 2).SetFrequency("minute", 30)}, "VTI", false},
		{"Date", args{apiKey, "", "VTI", (&PriceHistoryConfig{}).SetDate(time.Now().AddDate(0, 0, -1), time.Now())}, "VTI", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPriceHistory(tt.args.apiKey, tt.args.accessToken, tt.args.symbol, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPriceHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got.Symbol != tt.want {
				t.Errorf("GetPriceHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	type args struct {
		accessToken string
		accountID   int64
		config      *TransactionsConfig
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, getAccountID(), (&TransactionsConfig{}).SetDate(time.Now().AddDate(0, -2, -20), time.Now())}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTransactions(tt.args.accessToken, tt.args.accountID, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if len(*got) <= tt.want {
				t.Errorf("GetTransactions() = %v, want more than %v", got, tt.want)
			}
		})
	}
}

func TestGetTransaction(t *testing.T) {
	type args struct {
		accessToken   string
		accountID     int64
		transactionID int64
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, getAccountID(), 123}, &types.Transaction{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTransaction(tt.args.accessToken, tt.args.accountID, tt.args.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStreamerSubscriptionKeys(t *testing.T) {
	type args struct {
		accessToken string
		accountIDs  []int64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, []int64{getAccountID()}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStreamerSubscriptionKeys(tt.args.accessToken, tt.args.accountIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStreamerSubscriptionKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if len(got.Keys) <= tt.want {
				t.Errorf("GetStreamerSubscriptionKeys() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserPrincipals(t *testing.T) {
	type args struct {
		accessToken string
		fields      []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{accessToken, []string{UserPrincipalsFieldsStreamerSubscriptionKeys, UserPrincipalsFieldsStreamerConnectionInfo, UserPrincipalsFieldsPreferences, UserPrincipalsFieldsSurrogateIds}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserPrincipals(tt.args.accessToken, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPrincipals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			content, _ := json.MarshalIndent(got, "", "    ")
			t.Log(string(content))
			if got.UserID == tt.want {
				t.Errorf("GetUserPrincipals() = %v, not want %v", got, tt.want)
			}
		})
	}
}

func Test_Preferences(t *testing.T) {
	orignal, err := GetPreferences(accessToken, getAccountID())
	if err != nil {
		panic(err)
	}

	// 建立 order
	p := NewPreferences(PreferencesTimeoutTWOHOURS)
	err = UpdatePreferences(accessToken, getAccountID(), p)
	if err != nil {
		panic(err)
	}

	// 獲得建立的 order
	p, err = GetPreferences(accessToken, getAccountID())
	if err != nil {
		panic(err)
	}
	if p.AuthTokenTimeout != PreferencesTimeoutTWOHOURS {
		t.Errorf("want: %s get: %s\n", PreferencesTimeoutTWOHOURS, p.AuthTokenTimeout)
	}

	// 更新回來
	err = UpdatePreferences(accessToken, getAccountID(), orignal)
	if err != nil {
		panic(err)
	}
}

func Test_SavedOrder(t *testing.T) {
	// 建立 save order
	order := NewOrder("BWX", OrderAssetTypeEQUITY, OrderInstructionBuy, 1, 1)
	err := CreateSavedOrder(accessToken, getAccountID(), order)
	if err != nil {
		panic(err)
	}

	// 獲得建立的 save order
	orderSaveds, err := GetSavedOrdersbyPath(accessToken, getAccountID())
	if err != nil {
		panic(err)
	}
	orderSaved := (*orderSaveds)[len(*orderSaveds)-1]
	content, _ := json.MarshalIndent(orderSaved, "", "    ")
	t.Log(string(content))

	// 取代建立的 save order
	orderSaved.Order = NewOrder("BND", OrderAssetTypeEQUITY, OrderInstructionBuy, 1, 1)
	err = ReplaceSavedOrder(accessToken, getAccountID(), orderSaved)
	if err != nil {
		panic(err)
	}

	// 獲得取代的 save order
	orderSaved, err = GetSavedOrder(accessToken, getAccountID(), orderSaved.SavedOrderID)
	if err != nil {
		panic(err)
	}
	content, _ = json.MarshalIndent(orderSaved, "", "    ")
	t.Log(string(content))

	// 刪除取代的 save order
	err = DeleteSavedOrder(accessToken, getAccountID(), orderSaved.SavedOrderID)
	if err != nil {
		panic(err)
	}

	// 獲得刪除的 save order，應為空
	orderSaved, err = GetSavedOrder(accessToken, getAccountID(), orderSaved.SavedOrderID)
	if err == nil {
		content, _ = json.MarshalIndent(orderSaved, "", "    ")
		t.Error(string(content))
	}
}

// 暫不測試，因需戶口有錢
func test_Order(t *testing.T) {
	// 建立 order
	order := NewOrder("VTI", OrderAssetTypeEQUITY, OrderInstructionBuy, 10, 1)
	err := PlaceOrder(accessToken, getAccountID(), order)
	if err != nil {
		panic(err)
	}

	// 獲得建立的 order
	orders, err := GetOrdersByQuery(accessToken, getAccountID(), (&OrdersConfig{MaxResults: 100}).SetEnteredTime(time.Now().AddDate(0, 0, -1), time.Now()))
	if err != nil {
		panic(err)
	}
	order = nil
	for _, o := range *orders {
		if o.Status == "WORKING" {
			order = (*orders)[0]
		}
	}
	if order == nil {
		panic("can not find order")
	}
	content, _ := json.MarshalIndent(order, "", "    ")
	t.Log(string(content))

	orders, err = GetOrdersByPath(accessToken, getAccountID(), (&OrdersConfig{}).SetEnteredTime(time.Now().AddDate(0, 0, -1), time.Now()))
	if err != nil {
		panic(err)
	}
	order = nil
	for _, o := range *orders {
		if o.Status == "WORKING" {
			order = (*orders)[0]
		}
	}
	if order == nil {
		panic("can not find order")
	}
	content, _ = json.MarshalIndent(order, "", "    ")
	t.Log(string(content))

	// 獲得將被取代的 order
	order, err = GetOrder(accessToken, getAccountID(), order.OrderID)
	if err != nil {
		panic(err)
	}
	content, _ = json.MarshalIndent(order, "", "    ")
	t.Log(string(content))

	orderID := order.OrderID
	// 取代建立的 order
	order = NewOrder("BND", OrderAssetTypeEQUITY, OrderInstructionBuy, 10, 1)
	err = ReplaceOrder(accessToken, getAccountID(), orderID, order)
	if err != nil {
		panic(err)
	}

	// 獲得取代的 order
	orders, err = GetOrdersByPath(accessToken, getAccountID(), (&OrdersConfig{}).SetEnteredTime(time.Now().AddDate(0, 0, -1), time.Now()))
	if err != nil {
		panic(err)
	}
	order = nil
	for _, o := range *orders {
		if o.Status == "WORKING" {
			order = (*orders)[0]
		}
	}
	if order == nil {
		panic("can not find order")
	}
	content, _ = json.MarshalIndent(order, "", "    ")
	t.Error(string(content))

	// 取消取代的 order
	err = CancelOrder(accessToken, getAccountID(), order.OrderID)
	if err != nil {
		panic(err)
	}

	// 獲得刪除的 order，狀態為 canceled
	order, _ = GetOrder(accessToken, getAccountID(), order.OrderID)
	if order.Status != "CANCELED" {
		content, _ = json.MarshalIndent(order, "", "    ")
		t.Error(string(content))
	}
}

func Test_Watchlist(t *testing.T) {
	// 建立 watchlist
	wb := NewWatchlistBasic("test", []string{"VTI", "BND"})
	err := CreateWatchlist(accessToken, getAccountID(), wb)
	if err != nil {
		panic(err)
	}

	// 獲得建立的 watchlist
	ws, err := GetWatchlistsforSingleAccount(accessToken, getAccountID())
	if err != nil {
		panic(err)
	}
	var w *types.Watchlist
	for _, v := range *ws {
		if v.Name == "test" {
			w = v
			break
		}
	}
	content, _ := json.MarshalIndent(w, "", "    ")
	t.Log(string(content))

	ws, err = GetWatchlistsforMultipleAccounts(accessToken)
	if err != nil {
		panic(err)
	}

	for _, v := range *ws {
		if v.Name == "test" {
			w = v
			break
		}
	}
	content, _ = json.MarshalIndent(w, "", "    ")
	t.Log(string(content))

	// 取代建立的 watchlist
	wb = NewWatchlistBasic("testReplace", []string{"VTI", "BND", "VBR"})
	err = ReplaceWatchlist(accessToken, getAccountID(), w.WatchlistID, wb)
	if err != nil {
		panic(err)
	}

	// 獲得取代的 watchlist
	ws, err = GetWatchlistsforMultipleAccounts(accessToken)
	if err != nil {
		panic(err)
	}
	for _, v := range *ws {
		if v.Name == "testReplace" {
			w = v
			break
		}
	}
	content, _ = json.MarshalIndent(w, "", "    ")
	t.Log(string(content))

	// 更正 watchlist
	wb = NewWatchlistBasic("testPatch", []string{"VTI", "BND", "VBR", "VGK"})
	err = UpdateWatchlist(accessToken, getAccountID(), w.WatchlistID, wb)
	if err != nil {
		panic(err)
	}

	// 獲得更正的 watchlist
	w, err = GetWatchlist(accessToken, getAccountID(), w.WatchlistID)
	if err != nil {
		panic(err)
	}
	content, _ = json.MarshalIndent(w, "", "    ")
	t.Log(string(content))

	// 刪除取代的 watchlist
	err = DeleteWatchlist(accessToken, getAccountID(), w.WatchlistID)
	if err != nil {
		panic(err)
	}

	// 獲得刪除的 watchlist，應為空
	ws, err = GetWatchlistsforMultipleAccounts(accessToken)
	if err != nil {
		panic(err)
	}
	for _, v := range *ws {
		if v.Name == "testPatch" {
			content, _ = json.MarshalIndent(v, "", "    ")
			t.Error(string(content))
			break
		}
	}
}

func TestOpenURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{"https://golang.org/"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OpenURI(tt.args.uri); (err != nil) != tt.wantErr {
				t.Errorf("OpenURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
