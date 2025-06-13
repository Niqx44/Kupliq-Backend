package models

type Reservasi struct {
	IDReservasi      int    `json:"id_reservasi"`
	IDCustomer       int    `json:"id_costumer"`
	TanggalReservasi string `json:"tanggal_reservasi"`
	WaktuReservasi   string `json:"waktu_reservasi"`
	Keterangan       string `json:"keterangan"`
	Status           string `json:"status"`
	NamaCostumer     string `json:"nama_costumer"`
	NoTelpCostumer   string `json:"notelp_costumer"`
}
