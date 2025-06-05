package api

import (
	"net/http"
	"backend-kupliq/internal/db"
	"backend-kupliq/internal/models"
	"log"
	"encoding/json"
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
	"time"
	"path/filepath"
	"backend-kupliq/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"io"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}


// GetCostumerHandler - Mendapatkan semua data costumer
func GetCostumerHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_costumer, nama_costumer, password, email, notelp_costumer, id_role, COALESCE(foto_profile, '') FROM costumer")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var costumers []models.Costumer
	for rows.Next() {
		var costumer models.Costumer
		if err := rows.Scan(&costumer.IDCostumer, &costumer.NamaCostumer, &costumer.Password, &costumer.Email, &costumer.NoTelp, &costumer.IDRole, &costumer.Foto_Profile); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		costumers = append(costumers, costumer)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(costumers)
}

// GetCostumerByIDHandler - Mendapatkan data costumer berdasarkan ID
func GetCostumerByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data costumer berdasarkan ID
	var costumer models.Costumer
	err = database.QueryRow("SELECT id_costumer, nama_costumer, password, email, notelp_costumer,id_role, COALESCE(foto_profile, '') FROM costumer WHERE id_costumer=$1", id).
		Scan(&costumer.IDCostumer, &costumer.NamaCostumer, &costumer.Password, &costumer.Email, &costumer.NoTelp, &costumer.IDRole, &costumer.Foto_Profile)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Costumer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data costumer dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(costumer)
}

// CreateCostumerHandler - Menambahkan data costumer baru
func CreateCostumerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var costumer models.Costumer
	err = json.NewDecoder(r.Body).Decode(&costumer)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	costumer.IDRole = "2"

	query := "INSERT INTO costumer (nama_costumer, password, email, notelp_costumer, id_role) VALUES ($1, $2, $3, $4, $5)"
	_, err = database.Exec(query, costumer.NamaCostumer, costumer.Password, costumer.Email, costumer.NoTelp, costumer.IDRole)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Data Costumer berhasil ditambahkan:", costumer)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Data Costumer berhasil ditambahkan"))
}

// UpdateCostumerHandler - Mengupdate data menu berdasarkan ID
func UpdateCostumerHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var costumer models.Costumer
	if err := json.NewDecoder(r.Body).Decode(&costumer); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE costumer SET nama_costumer=$1, password=$2, email=$3, notelp_costumer=$4, id_role=$5, foto_profile=$6 WHERE id_costumer=$7",
		costumer.NamaCostumer, costumer.Password, costumer.Email, costumer.NoTelp, costumer.IDRole, costumer.Foto_Profile, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(costumer)
}

// GetMenuHandler - Mendapatkan semua data guru
func GetMenuHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_menu, nama_menu, harga_menu, kategori, deskripsi FROM menu")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var menu models.Menu
		if err := rows.Scan(&menu.IDMenu, &menu.NamaMenu, &menu.HargaMenu, &menu.Kategori, &menu.Deskripsi); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		menus = append(menus, menu)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menus)
}

// CreateMenuHandler - Menambahkan data menu baru
func CreateMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var menu models.Menu
	err = json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO menu (nama_menu, harga_menu, kategori, deskripsi) VALUES ($1, $2, $3, $4)"
	_, err = database.Exec(query, menu.NamaMenu, menu.HargaMenu, menu.Kategori, menu.Deskripsi)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Menu berhasil ditambahkan:", menu)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Menu berhasil ditambahkan"))
}

// UpdateMenuHandler - Mengupdate data menu berdasarkan ID
func UpdateMenuHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var menu models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE menu SET nama_menu=$1, harga_menu=$2, kategori=$3, deskripsi=$4 WHERE id_menu=$5",
		menu.NamaMenu, menu.HargaMenu, menu.Kategori, menu.Deskripsi, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menu)
}

// DeleteMenuHandler - Menghapus data menu berdasarkan ID
func DeleteMenuHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting menu with ID:", id) // Log ID yang akan dihapus

	// Menghapus menu dari database berdasarkan id_menu
	result, err := database.Exec("DELETE FROM menu WHERE id_menu=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No menu found with ID:", id)
		http.Error(w, "Menu not found", http.StatusNotFound)
		return
	}

	log.Println("Menu successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Menu berhasil dihapus"))
}

// GetMenuByIDHandler - Mendapatkan data menu berdasarkan ID
func GetMenuByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data menu berdasarkan ID
	var menu models.Menu
	err = database.QueryRow("SELECT id_menu, nama_menu, harga_menu, kategori, deskripsi FROM menu WHERE id_menu=$1", id).
		Scan(&menu.IDMenu, &menu.NamaMenu, &menu.HargaMenu, &menu.Kategori, &menu.Deskripsi)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data menu dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Login attempt:", creds.Username, creds.Password)

	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var user models.User
	var query string
	var userID int

	// Cek tabel admin
	query = `SELECT id_admin, id_role, email, password FROM admin WHERE email = $1`
	err = conn.QueryRow(query, creds.Username).Scan(&userID, &user.RoleID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		// Cek tabel customer jika tidak ditemukan di admin
		query = `SELECT id_costumer, id_role, email, password FROM costumer WHERE email = $1`
		err = conn.QueryRow(query, creds.Username).Scan(&userID, &user.RoleID, &user.Username, &user.Password)

		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			fmt.Println("Customer query error:", err)
			http.Error(w, "Database error (customer)", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Database error (admin)", http.StatusInternalServerError)
		return
	}

	// Verifikasi password
	if creds.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Jika password kamu pakai hash bcrypt:
	// if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
	//     http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	//     return
	// }

	response := map[string]interface{}{
		"id":      userID,
		"id_role": user.RoleID,
		"token":   "dummy-token",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler
func GetReservasiByIDcustomer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idCustomerStr := vars["id_costumer"]
    log.Println("id_customer dari URL:", idCustomerStr) // debug log

    idCustomer, err := strconv.Atoi(idCustomerStr)
    if err != nil {
        http.Error(w, "id_customer harus berupa angka", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query("SELECT id_reservasi, id_costumer, tanggal_reservasi, waktu_reservasi, keterangan, status FROM reservasi WHERE id_costumer = $1", idCustomer)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var reservasiList []models.Reservasi
    for rows.Next() {
        var k models.Reservasi
        if err := rows.Scan(&k.IDReservasi, &k.IDCustomer, &k.TanggalReservasi, &k.WaktuReservasi, &k.Keterangan, &k.Status); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        reservasiList = append(reservasiList, k)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservasiList)
}

// CreateReservasiHandler - Menambahkan data guru baru
func CreateReservasiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var reservasi models.Reservasi
	err = json.NewDecoder(r.Body).Decode(&reservasi)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reservasi.Status = "Pending"

	query := "INSERT INTO reservasi (id_costumer, tanggal_reservasi, waktu_reservasi, keterangan, status) VALUES ($1, $2, $3, $4, $5)"
	_, err = database.Exec(query, reservasi.IDCustomer, reservasi.TanggalReservasi, reservasi.WaktuReservasi, reservasi.Keterangan, reservasi.Status)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Reservasi berhasil ditambahkan:", reservasi)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Reservasi berhasil ditambahkan"))
}

// GetReservasiHandler - Mendapatkan semua data reservasi
func GetReservasiHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_reservasi, id_costumer, tanggal_reservasi, waktu_reservasi, keterangan, status FROM reservasi")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reservasis []models.Reservasi
	for rows.Next() {
		var reservasi models.Reservasi
		if err := rows.Scan(&reservasi.IDReservasi, &reservasi.IDCustomer, &reservasi.TanggalReservasi, &reservasi.WaktuReservasi, &reservasi.Keterangan, &reservasi.Status); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		reservasis = append(reservasis, reservasi)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservasis)
}

// GetReservasiByIDHandler - Mendapatkan data reservasi berdasarkan ID
func GetReservasiByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data menu berdasarkan ID
	var reservasi models.Reservasi
	err = database.QueryRow("SELECT id_reservasi, id_costumer, tanggal_reservasi, waktu_reservasi, keterangan, status FROM reservasi WHERE id_reservasi=$1", id).
		Scan(&reservasi.IDReservasi, &reservasi.IDCustomer, &reservasi.TanggalReservasi, &reservasi.WaktuReservasi, &reservasi.Keterangan, &reservasi.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Reservasi not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data reservasi dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservasi)
}

// UpdateReservasiStatusHandler - Mengubah status reservasi berdasarkan ID
func UpdateReservasiStatusHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal terhubung ke database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var input models.Reservasi
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	// Update status reservasi berdasarkan ID
	query := "UPDATE reservasi SET status=$1 WHERE id_reservasi=$2"
	_, err = database.Exec(query, input.Status, id)
	if err != nil {
		http.Error(w, "Gagal memperbarui status reservasi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status reservasi berhasil diperbarui",
	})
}



func BuatPemesananHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Insert into pemesanan
	var idPemesanan int
	queryPemesanan := `
		INSERT INTO pemesanan (id_costumer, id_meja, tanggal_pemesanan, total_harga, status)
		VALUES ($1, $2, CURRENT_DATE, $3, 'pending') RETURNING id_pemesanan
	`
	err = dbConn.QueryRow(queryPemesanan, req.IDCustomer, req.IDMeja, req.TotalHarga).Scan(&idPemesanan)
	if err != nil {
		log.Println("Insert pemesanan error:", err)
		http.Error(w, "Failed to create pemesanan", http.StatusInternalServerError)
		return
	}

	// 2. Insert each item into detailpemesanan
	for _, item := range req.Items {
		queryDetail := `
			INSERT INTO detailpemesanan (id_pemesanan, id_menu, jumlah, sub_total)
			VALUES ($1, $2, $3, $4)
		`
		_, err := dbConn.Exec(queryDetail, idPemesanan, item.IDMenu, item.Jumlah, item.SubTotal)
		if err != nil {
			log.Println("Insert detailpemesanan error:", err)
			http.Error(w, "Failed to insert detail item", http.StatusInternalServerError)
			return
		}
	}

	// 3. Insert into pembayaran
	queryBayar := `
		INSERT INTO pembayaran (id_pemesanan, metode_pembayaran, status)
		VALUES ($1, $2, 'belum dibayar')
	`
	_, err = dbConn.Exec(queryBayar, idPemesanan, req.MetodePembayaran)
	if err != nil {
		log.Println("Insert pembayaran error:", err)
		http.Error(w, "Failed to create pembayaran", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Pemesanan berhasil dibuat"))
}

func GetPemesananByCustomerID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idCustomerStr := vars["id_costumer"]
    log.Println("id_costumer dari URL:", idCustomerStr)

    idCustomer, err := strconv.Atoi(idCustomerStr)
    if err != nil {
        http.Error(w, "id_costumer harus berupa angka", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    query := `
        SELECT 
            p.id_pemesanan, p.id_costumer, p.id_meja, p.total_harga, p.tanggal_pemesanan, p.status,
            d.id_menu, m.nama_menu, d.jumlah, d.sub_total,
            bayar.metode_pembayaran
        FROM pemesanan p
        JOIN detailpemesanan d ON p.id_pemesanan = d.id_pemesanan
        JOIN menu m ON d.id_menu = m.id_menu
        LEFT JOIN pembayaran bayar ON bayar.id_pemesanan = p.id_pemesanan
        WHERE p.id_costumer = $1
        ORDER BY p.id_pemesanan DESC;
    `
    rows, err := dbConn.Query(query, idCustomer)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    orderMap := make(map[int]*models.Order)

    for rows.Next() {
        var (
            idPemesanan, idCustomer, idMeja int
            totalHarga, subTotal            float64
            tanggal, namaMenu, status       string
            idMenu, jumlah                  int
            metodePembayaran                *string // pointer agar bisa null
        )

        err := rows.Scan(&idPemesanan, &idCustomer, &idMeja, &totalHarga, &tanggal, &status,
            &idMenu, &namaMenu, &jumlah, &subTotal, &metodePembayaran)
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }

        if _, exists := orderMap[idPemesanan]; !exists {
            orderMap[idPemesanan] = &models.Order{
                IDPemesanan:       idPemesanan,
                IDCustomer:        idCustomer,
                IDMeja:            idMeja,
                TotalHarga:        totalHarga,
                TanggalPemesanan:  tanggal,
                Status:            status,
                MetodePembayaran:  "", // default
                Items:             []models.Item{},
            }

            if metodePembayaran != nil {
                orderMap[idPemesanan].MetodePembayaran = *metodePembayaran
            }
        }

        item := models.Item{
            IDMenu:   idMenu,
            NamaMenu: namaMenu,
            Jumlah:   jumlah,
            SubTotal: subTotal,
        }

        orderMap[idPemesanan].Items = append(orderMap[idPemesanan].Items, item)
    }

    var result []models.Order
    for _, order := range orderMap {
        result = append(result, *order)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}


func GetAllPemesananHandler(w http.ResponseWriter, r *http.Request) {
    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    query := `
        SELECT 
            p.id_pemesanan, p.id_costumer, p.id_meja, p.total_harga, p.tanggal_pemesanan, p.status,
            d.id_menu, m.nama_menu, d.jumlah, d.sub_total
        FROM pemesanan p
        JOIN detailpemesanan d ON p.id_pemesanan = d.id_pemesanan
        JOIN menu m ON d.id_menu = m.id_menu
        ORDER BY p.id_pemesanan DESC;
    `

    rows, err := dbConn.Query(query)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    orderMap := make(map[int]*models.Order)

    for rows.Next() {
        var (
            idPemesanan, idCustomer, idMeja int
            totalHarga, subTotal           float64
            tanggal, namaMenu, status      string
            idMenu, jumlah                 int
        )

        err := rows.Scan(&idPemesanan, &idCustomer, &idMeja, &totalHarga, &tanggal, &status,
            &idMenu, &namaMenu, &jumlah, &subTotal)
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }

        if _, exists := orderMap[idPemesanan]; !exists {
            orderMap[idPemesanan] = &models.Order{
                IDPemesanan:       idPemesanan,
                IDCustomer:        idCustomer,
                IDMeja:            idMeja,
                TotalHarga:        totalHarga,
                TanggalPemesanan:  tanggal,
                Status:            status,
                Items:             []models.Item{},
            }
        }

        item := models.Item{
            IDMenu:   idMenu,
            NamaMenu: namaMenu,
            Jumlah:   jumlah,
            SubTotal: subTotal,
        }

        orderMap[idPemesanan].Items = append(orderMap[idPemesanan].Items, item)
    }

    var result []models.Order
    for _, order := range orderMap {
        result = append(result, *order)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func UpdateStatusPemesanan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPemesananStr := vars["id_pemesanan"]

	idPemesanan, err := strconv.Atoi(idPemesananStr)
	if err != nil {
		http.Error(w, "id_pemesanan harus berupa angka", http.StatusBadRequest)
		return
	}

	type StatusUpdate struct {
		Status string `json:"status"`
	}

	var input StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Gagal membaca input", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal koneksi database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	query := `UPDATE pemesanan SET status = $1 WHERE id_pemesanan = $2`
	_, err = dbConn.Exec(query, input.Status, idPemesanan)
	if err != nil {
		log.Println("Gagal update status:", err)
		http.Error(w, "Gagal update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status berhasil diubah"))
}


const (
	bucketName  = "kupliqcafe-profile"
	s3Folder    = "profile/"
)

func UploadFotoProfileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Start upload handler")
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("foto")
	if err != nil {
		log.Printf("FormFile error: %v", err)
		http.Error(w, "File not found in request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := r.FormValue("id_costumer")
	if id == "" {
		log.Println("id_costumer is missing")
		http.Error(w, "id_costumer is required", http.StatusBadRequest)
		return
	}
	log.Println("File uploaded:", handler.Filename)

	fileExt := filepath.Ext(handler.Filename)
	fileName := fmt.Sprintf("profile_%s_%d%s", id, time.Now().Unix(), fileExt)
	s3Key := s3Folder + fileName

	// Baca sebagian isi file untuk deteksi tipe
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Printf("Read buffer error: %v", err)
		http.Error(w, "Failed to read file buffer", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)
	// Reset posisi baca file
	file.Seek(0, io.SeekStart)

	log.Println("Detected content type:", contentType)
	log.Println("Uploading to S3:", s3Key)

	// Upload ke S3
	_, err = config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		http.Error(w, "Failed to upload to S3", http.StatusInternalServerError)
		return
	}

	log.Println("S3 upload success")

	s3URL := fmt.Sprintf("https://%s.s3.ap-southeast-3.amazonaws.com/%s", bucketName, s3Key)

	// Update database
	database, err := db.ConnectToDB()
	if err != nil {
		log.Printf("DB connect error: %v", err)
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	_, err = database.Exec("UPDATE costumer SET foto_profile = $1 WHERE id_costumer = $2", s3URL, id)
	if err != nil {
		log.Printf("DB update error: %v", err)
		http.Error(w, "Failed to update database", http.StatusInternalServerError)
		return
	}
	log.Println("DB update success")


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Foto berhasil diunggah",
		"url":     s3URL,
	})
	log.Println("Response sent to client")

}




