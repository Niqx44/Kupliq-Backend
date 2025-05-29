package models

type Order struct {
	IDPemesanan int     `json:"id_pemesanan"`
	IDCustomer  int     `json:"id_costumer"`
	IDMeja      int     `json:"id_meja"`
	TotalHarga  float64 `json:"total_harga"`
	TanggalPemesanan     string  `json:"tanggal_pemesanan"`
	Status      string  `json:"status"`
	Items       []Item  `json:"items"`
}