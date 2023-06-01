package entity

type OtherPosition struct {
	EncryptedUid string `json:"encryptedUid"`
	TradeType    string `json:"tradeType"`
}
type LeaderBoardPayload struct {
	TradeType      string `json:"tradeType"`
	StatisticsType string `json:"statisticsType"`
	PeriodType     string `json:"periodType"`
	IsShared       bool   `json:"isShared"`
	IsTrader       bool   `json:"isTrader"`
}
