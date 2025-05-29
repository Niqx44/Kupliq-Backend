package models

type Item struct {
	IDMenu    int     `json:"id_menu"`
	NamaMenu  string  `json:"nama_menu"`
	Jumlah    int     `json:"jumlah"`
	SubTotal  float64 `json:"sub_total"`
}