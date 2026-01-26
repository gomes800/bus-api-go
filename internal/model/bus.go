package model

type BusPosition struct {
	Ordem            string `json:"ordem"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	DataHora         string `json:"datahora"`
	Velocidade       string `json:"velocidade"`
	Linha            string `json:"linha"`
	DataHoraEnvio    string `json:"datahoraenvio"`
	DataHoraServidor string `json:"datahoraservidor"`
}
