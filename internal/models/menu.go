package models

type Menu struct {
	IDMenu        int    `json:"id_menu"`
	NamaMenu      string `json:"nama_menu"`
	HargaMenu 	  string `json:"harga_menu"`
	Kategori      string `json:"kategori"`
	Deskripsi     string `json:"deskripsi"`
}
