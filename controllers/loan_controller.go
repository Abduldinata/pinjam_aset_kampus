package controllers

import (
	"net/http"
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. TAMPILKAN FORM PEMINJAMAN
func CreateLoan(c *gin.Context) {
	var items []models.Item
	// Ambil barang yang stoknya > 0 saja
	config.DB.Where("stock > ?", 0).Find(&items)

	c.HTML(http.StatusOK, "user/loan_form.html", gin.H{
		"Items": items,
	})
}


// 3. TAMPILKAN RIWAYAT PEMINJAMAN USER
func HistoryLoan(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(uint)
	var loans []models.Loan

	// Preload "Item" artinya: "Tolong ambilkan juga data nama barangnya dari tabel items"
	config.DB.Preload("Item").Where("user_id = ?", userID).Find(&loans)

	// 3. Ambil Notifikasi
	var notifs []models.Notification
	// Gunakan error handling biar aman
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifs).Error; err != nil {
		notifs = []models.Notification{} // Jika error, kasih array kosong
	}

	// 4. Hitung Jumlah Belum Dibaca
	unreadCount := 0
	for _, n := range notifs {
		if !n.IsRead {
			unreadCount++
		}
	}

	c.HTML(http.StatusOK, "user/history.html", gin.H{
		"Loans":       loans,
		"Notifs":      notifs,
		"UnreadCount": unreadCount,
	})
}

func StoreLoan(c *gin.Context) {
	// 1. Ambil User ID
	userID, _ := c.MustGet("user_id").(uint)

	// 2. Ambil Input dari Modal
	itemID, _ := strconv.Atoi(c.PostForm("item_id"))
	duration, _ := strconv.Atoi(c.PostForm("duration")) // Input berupa angka (misal: 3 hari)
	notes := c.PostForm("notes")

	// 3. Hitung Tanggal
	borrowDate := time.Now()                 // Tanggal pinjam = HARI INI
	dueDate := borrowDate.AddDate(0, 0, duration) // Tanggal kembali = Hari ini + Durasi

	// 4. Mulai Transaksi (Sama seperti sebelumnya)
	tx := config.DB.Begin()

	var item models.Item
	if err := tx.First(&item, itemID).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusBadRequest, "Barang tidak ditemukan")
		return
	}

	if item.Stock <= 0 {
		tx.Rollback()
		c.String(http.StatusBadRequest, "Stok habis!")
		return
	}

	// Buat Data Pinjam
	loan := models.Loan{
		UserID:     userID,
		ItemID:     uint(itemID),
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		Status:     "dipinjam",
		Notes:      notes,
	}

	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		c.String(500, "Gagal simpan transaksi")
		return
	}

	// Kurangi Stok
	item.Stock = item.Stock - 1
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	// Redirect ke Riwayat
	c.Redirect(http.StatusFound, "/user/history")
}

// ... (kode sebelumnya)

// 4. PROSES PENGEMBALIAN BARANG (ADMIN)
func ReturnLoan(c *gin.Context) {
	// Ambil ID Transaksi Peminjaman dari Form
	loanID := c.PostForm("loan_id")

	// Mulai Transaksi (Wajib, karena kita ubah 2 tabel sekaligus)
	tx := config.DB.Begin()

	// 1. Cari Data Peminjaman
	var loan models.Loan
	if err := tx.First(&loan, loanID).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusNotFound, "Data peminjaman tidak ditemukan")
		return
	}

	// Cek: Jangan sampai barang yang sudah kembali diproses lagi
	if loan.Status == "kembali" {
		tx.Rollback()
		c.String(http.StatusBadRequest, "Barang ini sudah dikembalikan sebelumnya")
		return
	}

	// 2. Update Status Peminjaman & Tanggal Kembali
	now := time.Now()
	loan.Status = "kembali"
	loan.ReturnDate = &now // Pointer ke waktu sekarang
	
	// (Opsional) Cek apakah terlambat?
	if now.After(loan.DueDate) {
		loan.Notes = loan.Notes + " [Dikembalikan Terlambat]"
		// Bisa ubah status jadi 'terlambat' jika mau, tapi 'kembali' lebih logis untuk inventory.
	}

	if err := tx.Save(&loan).Error; err != nil {
		tx.Rollback()
		c.String(500, "Gagal update status peminjaman")
		return
	}

	// 3. Update Stok Barang (BERTAMBAH 1)
	var item models.Item
	if err := tx.First(&item, loan.ItemID).Error; err != nil {
		tx.Rollback()
		c.String(500, "Barang tidak ditemukan")
		return
	}

	item.Stock = item.Stock + 1 // Stok Balik
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.String(500, "Gagal update stok barang")
		return
	}

	// 4. Simpan Perubahan
	tx.Commit()

	// Balik ke halaman laporan
	c.Redirect(http.StatusFound, "/admin/reports")
}

// 5. KIRIM NOTIFIKASI PENGINGAT (ADMIN)
func SendReminder(c *gin.Context) {
	loanID := c.PostForm("loan_id")

	// Cari Data Peminjaman (Kita butuh tahu Siapa user-nya dan Apa barangnya)
	var loan models.Loan
	// Preload User & Item biar kita bisa sebut nama mereka di pesan
	if err := config.DB.Preload("User").Preload("Item").First(&loan, loanID).Error; err != nil {
		c.String(404, "Data tidak ditemukan")
		return
	}

	// Buat Pesan Notifikasi
	pesan := "Halo " + loan.User.Name + ", mohon segera kembalikan aset: " + loan.Item.Name + ". Batas waktu: " + loan.DueDate.Format("02 Jan 2006")

	// Simpan ke Tabel Notifications
	notif := models.Notification{
		UserID:  loan.UserID,
		Message: pesan,
		IsRead:  false, // Belum dibaca
	}

	config.DB.Create(&notif)

	// Balik ke halaman laporan
	c.Redirect(http.StatusFound, "/admin/reports")
}