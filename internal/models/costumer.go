package models

type Costumer struct {
	IDCostumer      int    `json:"id_user"`
	NamaCostumer 	string `json:"nama_costumer"`
	Email 			string `json:"email"`
	NoTelp 			string `json:"notelp_costumer"`
	Password 		string `json:"password"`
	IDRole     		string `json:"id_role"`
	Foto_Profile 	string `json:"foto_profile"`
}
