package models

type QuetionGPT struct {
	PreguntaID  string   `json:preguntaID`
	RespuestaID string   `json:respuestaID`
	Contexto    []string `json:context`
}

type ContinueChat struct {
	Context []string `json:context`
}
