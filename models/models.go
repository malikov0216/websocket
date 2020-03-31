package models

type Bid []string

// Ask info with price and amount
type Ask []string

// WsDepth depth info
type Depth struct {
	Event         string `json:"e"`
	Time          int64  `json:"E"`
	Symbol        string `json:"s"`
	UpdateID      int64  `json:"u"`
	FirstUpdateID int64  `json:"U"`
	Bids          []Bid  `json:"b"`
	Asks          []Ask  `json:"a"`
}

type LiveStream struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

type DepthWebSocket func(message []byte)
type ErrorWebSocket func(err error)
