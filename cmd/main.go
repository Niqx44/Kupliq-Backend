package main

import (
	"backend-kupliq/config"
	"backend-kupliq/internal/api" // Pastikan path impor sesuai dengan folder proyek kamu
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config.InitS3()
	// Membuat router
	r := mux.NewRouter()

	// Menangani route untuk /guru/{id}
	r.HandleFunc("/costumer", api.GetCostumerHandler).Methods("GET")
	r.HandleFunc("/costumer/{id}", api.GetCostumerByIDHandler).Methods("GET")
	r.HandleFunc("/costumer", api.CreateCostumerHandler).Methods("POST")
	r.HandleFunc("/costumer/{id}", api.UpdateCostumerHandler).Methods("PUT")

	r.HandleFunc("/menu", api.GetMenuHandler).Methods("GET")
	r.HandleFunc("/menu/{id}", api.GetMenuByIDHandler).Methods("GET")
	r.HandleFunc("/menu", api.CreateMenuHandler).Methods("POST")
	r.HandleFunc("/menu/{id}", api.UpdateMenuHandler).Methods("PUT")
	r.HandleFunc("/menu/{id}", api.DeleteMenuHandler).Methods("DELETE")

	r.HandleFunc("/reservasi/{id_costumer}", api.GetReservasiByIDcustomer).Methods("GET")
	r.HandleFunc("/reservasi", api.GetReservasiHandler).Methods("GET")
	r.HandleFunc("/reservasi/by/{id}", api.GetReservasiByIDHandler).Methods("GET")
	r.HandleFunc("/reservasi", api.CreateReservasiHandler).Methods("POST")
	r.HandleFunc("/reservasi/{id}/status", api.UpdateReservasiStatusHandler).Methods("PUT")

	r.HandleFunc("/pemesanan", api.BuatPemesananHandler).Methods("POST")
	r.HandleFunc("/pemesanan/all", api.GetAllPemesananHandler).Methods("GET")

	r.HandleFunc("/pemesanan/status/{id_pemesanan}", api.UpdateStatusPemesanan).Methods("PUT")

	r.HandleFunc("/reservasi/{id}", api.DeleteReservasiHandler).Methods("DELETE")

	r.HandleFunc("/pemesanan/{id_costumer}", api.GetPemesananByCustomerID).Methods("GET")

	r.HandleFunc("/login", api.LoginHandler)
	r.HandleFunc("/upload-foto-profile", api.UploadFotoProfileHandler).Methods("POST")
	r.HandleFunc("/upload-foto-menu", api.UploadFotoMenuHandler).Methods("POST")

	// Menambahkan CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Ganti dengan URL frontend Anda
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Menambahkan middleware CORS
	handler := c.Handler(r)

	// Menjalankan server di port 8080
	log.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
