package models

type OrderRequest struct {
    IDCustomer       int            `json:"id_customer"`
    IDMeja     		 int          `json:"id_meja"`
    Items            []ItemRequest  `json:"items"`
    TotalHarga       int            `json:"total_harga"`
    MetodePembayaran string         `json:"metode_pembayaran"`
}