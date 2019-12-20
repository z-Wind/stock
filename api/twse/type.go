package twse

import (
	"fmt"
)

type Data interface {
	GetError() error
}

type Trade struct {
	MilliSecond string `json:"tlong"` //資料時間（毫秒）
	Date        string `json:"d"`     //今日日期
	Time        string `json:"t"`     //資料時間

	FullName string `json:"nf"` //全名
	Name     string `json:"n"`  //名字
	Symbol   string `json:"c"`  //股要代碼
	Channel  string `json:"ch"` //1101.tw

	BestBuyAmount      string  `json:"f"`         //最佳五檔賣出數量
	BestBuyPrice       string  `json:"b"`         //最佳五檔買入價格
	BestSellAmount     string  `json:"g"`         //最佳五檔買入數量
	BestSellPrice      string  `json:"a"`         //最佳五檔賣出價格
	TradePrice         float64 `json:"z,string"`  //最近成交價
	PreviousTradePrice float64 `json:"pz,string"` //前一個成交價
	YesterdayPrice     float64 `json:"y,string"`  //昨天收價
	Open               float64 `json:"o,string"`  //開盤價
	DayLow             float64 `json:"l,string"`  //今日最低
	DayHigh            float64 `json:"h,string"`  //今日最高

	Ex string `json:"ex"` //上市上櫃

	Flag             float64 `json:"ip,string"` //好像是一個 flag，3 是暫緩收盤股票, 2 是趨漲, 1 是趨跌
	DownLimit        float64 `json:"w,string"`  //跌停點
	UpLimit          float64 `json:"u,string"`  //漲停點
	CumulativeVolume int64   `json:"v,string"`  //當日累計成交量

	TemporalVolume float64 `json:"tv,string"` //當盤成交量
	//   TemporalVolume float64 `json:"ps"` //當盤成交量?
	//   TemporalVolume float64 `json:"s"` //當盤成交量?

	AfterHourTime       string  `json:"ot"`        //盤後定價時間
	AfterHourTradePrice float64 `json:"oz,string"` //盤後成交價
	AfterHourBuyPrice   float64 `json:"oa,string"` //盤後賣出價格
	AfterHourSellPrice  float64 `json:"ob,string"` //盤後買入價格
	AfterHourVolume     int64   `json:"ov,string"` //盤後成交量

	//   "ts":"0",
	//   "fv":"7",
	//   "tk0":"1101.tw_tse_20181018_B_9999154936",
	//   "tk1":"1101.tw_tse_20181018_B_9999132899",
	//   "it":"12",
	//   "mt":"000000",
	//   "i":"01",
	//   "p":"0",

}

type Message struct {
	Trades []Trade `json:"msgArray"`
	Infomation
}

func (m *Message) GetError() error {
	if m.Infomation.RTcode != "0000" {
		return fmt.Errorf("Error code %s", m.Infomation.RTmessage)
	}
	if len(m.Trades) == 0 {
		return fmt.Errorf("no Data")
	}
	return nil
}

type Infomation struct {
	UserDelay int64     `json:"userDelay"`
	RTmessage string    `json:"rtmessage"`
	Referer   string    `json:"referer"`
	QueryTime QueryTime `json:"queryTime"`
	RTcode    string    `json:"rtcode"`
}

type QueryTime struct {
	SysTime           string `json:"sysTime"`
	SessionLatestTime int64  `json:"sessionLatestTime"`
	SysDate           string `json:"sysDate"`
	SessionFromTime   int64  `json:"sessionFromTime"`
	StockInfoItem     int64  `json:"stockInfoItem"`
	ShowChart         bool   `json:"showChart"`
	SessionStr        string `json:"sessionStr"`
	StockInfo         int64  `json:"stockInfo"`
}

type QuoteHistory struct {
	Stat   string     `json:"stat"`
	Date   string     `json:"date"`
	Title  string     `json:"title"`
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data"`
	Notes  []string   `json:"notes"`
}

func (q *QuoteHistory) GetError() error {
	if q.Stat != "OK" {
		return fmt.Errorf("Error %s %v", q.Stat, q.Notes)
	}
	return nil
}
