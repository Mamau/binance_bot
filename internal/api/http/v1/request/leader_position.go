package request

type LeaderPosition struct {
	EncryptedUid string `form:"encryptedUid"`
	TradeType    string `form:"tradeType"`
}
