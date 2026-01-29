package dto

// RejectCallRequest request para rejeitar chamada
type RejectCallRequest struct {
	CallFrom string `json:"CallFrom" example:"5511999999999@s.whatsapp.net"`
	CallID   string `json:"CallID" example:"ABCD1234567890"`
}

// CallActionResponse resposta de acao de chamada
type CallActionResponse struct {
	Details string `json:"Details" example:"Call rejected"`
}
