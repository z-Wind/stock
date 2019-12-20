package types

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJson(t *testing.T) {
	path := "./data_test.json"
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	origin, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	var dataTest struct {
		Order                    SavedOrder      `json:"Order"`
		Orders                   SavedOrders     `json:"Orders"`
		SavedOrder               SavedOrder      `json:"SavedOrder"`
		SavedOrders              SavedOrders     `json:"SavedOrders"`
		Account                  Account         `json:"Account"`
		Accounts                 Accounts        `json:"Accounts"`
		Instructment             []Instrument    `json:"Instrument"`
		Instructments            Instruments     `json:"Instruments"`
		Hours                    MarketHours     `json:"Hours"`
		Hour                     MarketHours     `json:"Hour"`
		Movers                   Movers          `json:"Movers"`
		OptionChain              OptionChain     `json:"OptionChain"`
		PriceHistory             CandleList      `json:"PriceHistory"`
		Quotes                   Quotes          `json:"Quotes"`
		Quote                    Quotes          `json:"Quote"`
		Transactions             Transactions    `json:"Transactions"`
		Preferences              Preferences     `json:"Preferences"`
		StreamerSubscriptionKeys SubscriptionKey `json:"StreamerSubscriptionKeys"`
		Principals               UserPrincipal   `json:"Principals"`
		Watchlist                Watchlist       `json:"Watchlist"`
		Watchlists               Watchlists      `json:"Watchlists"`
	}

	err = json.Unmarshal(origin, &dataTest)
	if err != nil {
		t.Fatal(err)
	}

	after, err := json.MarshalIndent(dataTest, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(string(origin), string(after)) != 0 {
		//t.Errorf("get = %v, want %v", string(after), string(origin))
		t.Log("Different")
		// 存檔比較
		// f, err := os.Create("after.json")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// defer f.Close()
		// _, err = f.Write(after)
		// if err != nil {
		// 	t.Fatal(err)
		// }
	}
}
