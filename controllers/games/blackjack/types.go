package bj_controller

type Message struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
	Score string      `json:"score"`
}

type receivedMsg struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
