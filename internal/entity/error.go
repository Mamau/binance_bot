package entity

type LeaderBoardError struct {
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          interface{} `json:"data"`
	Success       bool        `json:"success"`
}
