package retroengine

type CountDownWebSocketResponse struct {
	MsgType string `json:"msgType"`
	Tick    string `json:"tick"`
}

type GameStateWebSocketResponse struct {
	MsgType string `json:"msgType"`
	Game    *Game  `json:"game"`
}
