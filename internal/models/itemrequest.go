package models

type ItemRequest struct {
    IDMenu   int `json:"id_menu"`
    Jumlah   int `json:"jumlah"`
    SubTotal int `json:"sub_total"`
}