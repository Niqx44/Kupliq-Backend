package models

type Costumer struct {
	IDCostumer      int    `json:"id_user"`
	NamaCostumer 	string `json:"nama_costumer"`
	Email 			string `json:"email"`
	Password 		string `json:"password"`
	IDRole     		string `json:"id_role"`
}
