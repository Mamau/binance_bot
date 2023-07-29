package entity

import (
	"math/big"
	"time"
)

type TransactionType string

const (
	External TransactionType = "EXTERNAL"
	Internal TransactionType = "INTERNAL"
)
const DividerETH = 1000000000000000000 // 1 Ether = 1000000000000000000 Wei

type Whale struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Wallet   string `json:"wallet"`
}

type WhaleAction struct {
	WhaleName     string    `json:"whaleName"`
	WhalePosition string    `json:"whalePosition"`
	Type          string    `json:"type"`
	ValueETH      float64   `json:"valueETH"`
	Balance       float64   `json:"balance"`
	Hash          string    `json:"hash"`
	Date          time.Time `json:"date"`
}

type ByDate []WhaleAction

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func NewWhaleAction(whaleName, wp string, transactionType TransactionType, value *big.Int, date time.Time, hash string, balance *big.Int) WhaleAction {
	weiPerEth := big.NewInt(DividerETH)
	ethPerWei := new(big.Float).SetInt(weiPerEth) // Конвертируем в тип big.Float для точности
	ethAmount := new(big.Float).SetInt(value)     // Конвертируем в тип big.Float для точности
	ethAmount.Quo(ethAmount, ethPerWei)           // Выполняем деление, чтобы получить значение в ETH
	ethAmountFloat, _ := ethAmount.Float64()      // Конвертируем big.Float в float64

	balanceAmount := new(big.Float).SetInt(balance) // Конвертируем в тип big.Float для точности
	ethAmount.Quo(balanceAmount, ethPerWei)         // Выполняем деление, чтобы получить значение в ETH
	balanceAmountFloat, _ := ethAmount.Float64()    // Конвертируем big.Float в float64

	return WhaleAction{
		WhaleName:     whaleName,
		WhalePosition: wp,
		Type:          GetTransactionName(transactionType),
		ValueETH:      ethAmountFloat,
		Date:          date,
		Hash:          hash,
		Balance:       balanceAmountFloat,
	}
}

func GetTransactionName(transactionType TransactionType) string {
	switch transactionType {
	case External:
		return "Вывел"
	case Internal:
		return "Ввод"
	default:
		return ""
	}
}

// AccountTransaction
// пример адреса https://api.blockchain.info/v2/eth/data/account/0xe8c19db00287e3536075114b2576c70773e039bd/wallet?page=0&size=20
type AccountTransaction struct {
	Hash                     string `json:"hash"`
	Nonce                    string `json:"nonce"`
	Balance                  string `json:"balance"`
	TransactionCount         string `json:"transactionCount"`
	InternalTransactionCount string `json:"internalTransactionCount"`
	TotalSent                string `json:"totalSent"`
	TotalReceived            string `json:"totalReceived"`
	TotalFees                string `json:"totalFees"`
	AccountTransactions      []struct {
		Hash             string `json:"hash"`
		BlockHash        string `json:"blockHash"`
		BlockNumber      string `json:"blockNumber"`
		From             string `json:"from"`
		To               string `json:"to"`
		Value            string `json:"value"`
		Nonce            string `json:"nonce"`
		GasPrice         string `json:"gasPrice"`
		GasLimit         string `json:"gasLimit"`
		GasUsed          string `json:"gasUsed"`
		Success          bool   `json:"success"`
		FirstSeen        string `json:"firstSeen"`
		TransactionIndex string `json:"transactionIndex"`
		State            string `json:"state"`
		Type             string `json:"type"`
		Timestamp        string `json:"timestamp"`
	} `json:"accountTransactions"`
}

func FindWhale(wallet string) Whale {
	for _, whale := range WhaleList {
		if whale.Wallet == wallet {
			return whale
		}
	}
	return Whale{}
}

var WhaleList = []Whale{
	{
		Name:     "RuneKek",
		Position: "соучредитель MakerDAO",
		Wallet:   "0x0f89d54b02ca570de82f770d33c7b7cf7b3c3394",
	},
	{
		Name:     "Джастин Сан",
		Position: "CEO компании Tron",
		Wallet:   "0x3ddfa8ec3052539b6c9549f12cea2c295cff5296",
	},
	{
		Name:     "Bankless",
		Position: "хз",
		Wallet:   "0x844e211e291077b11221c0f18615a64f2ff19c26",
	},
	{
		Name:     "degenspartan",
		Position: "генеральный директор и Web3 lovoor",
		Wallet:   "0x4e60be84870fe6ae350b563a121042396abe1eaf",
	},
	{
		Name:     "0xSisyphus",
		Position: "генеральный директор Unibot",
		Wallet:   "0x4dbe965abcb9ebc4c6e9d95aeb631e5b58e70d5b",
	},
	{
		Name:     "Стани Кулечев",
		Position: "CEO Lens Protocol и Aave",
		Wallet:   "0xf5fb27b912d987b5b6e02a1b1be0c1f0740e2c6f",
	},
	{
		Name:     "Патрисио Ворталтер",
		Position: "ранний инвестор $RPL и CEO poapxyz",
		Wallet:   "0x57757e3d981446d585af0d9ae4d7df6d64647806",
	},
	{
		Name:     "Михаил Егоров",
		Position: "генеральный директор Curve Finance",
		Wallet:   "0x7a16ff8270133f063aab6c9977183d9e72835428",
	},
	{
		Name:     "Каин Ворвик",
		Position: "генеральный директор компании Synthetix",
		Wallet:   "0x27cc4d6bc95b55a3a981bf1f1c7261cda7bb0931",
	},
	{
		Name:     "Виталик Бутерин",
		Position: "сооснователь Ethereum",
		Wallet:   "0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
	},
	{
		Name:     "Алекс Сваневик",
		Position: "CEO Nansen",
		Wallet:   "0xc006544b93e86e8999623c1d12d2e352c61c8123",
	},
	{
		Name:     "Эндрю Канг",
		Position: "трейдер DEX (GMX, Mux Protocol)",
		Wallet:   "0xe8c19db00287e3536075114b2576c70773e039bd",
	},
	{
		Name:     "Артур Хейс",
		Position: "один из крупнейших владельцев $GMX",
		Wallet:   "0x534a0076fb7c2b1f83fa21497429ad7ad3bd7587",
	},
	{
		Name:     "Тетранод",
		Position: "PenilessWassie",
		Wallet:   "0x9c5083dd4838e120dbeac44c052179692aa5dac5",
	},
}
