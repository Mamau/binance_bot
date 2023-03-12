package entity

type LeaderBoard struct {
	Code          string  `json:"code"`
	Message       *string `json:"message"`
	MessageDetail *string `json:"messageDetail"`
	Data          struct {
		OtherPositionRetList []struct {
			Symbol          string  `json:"symbol"`
			EntryPrice      float64 `json:"entryPrice"`
			MarkPrice       float64 `json:"markPrice"`
			Pnl             float64 `json:"pnl"`
			Roe             float64 `json:"roe"`
			UpdateTime      []int   `json:"updateTime"`
			Amount          float64 `json:"amount"`
			UpdateTimeStamp int64   `json:"updateTimeStamp"`
			Yellow          bool    `json:"yellow"`
			TradeBefore     bool    `json:"tradeBefore"`
			Leverage        int     `json:"leverage"`
		} `json:"otherPositionRetList"`
		UpdateTime      []int `json:"updateTime"`
		UpdateTimeStamp int64 `json:"updateTimeStamp"`
	} `json:"data"`
	Success bool `json:"success"`
}
