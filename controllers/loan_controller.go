package controllers

import (
	"net/http"
	"os"
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. TAMPILKAN FORM PEMINJAMAN
func CreateLoan(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(uint)
	var items []models.Item
	// Ambil barang yang stoknya > 0 saja
	config.DB.Where("stock > ?", 0).Order("id asc").Find(&items)

	// Ambil Notifikasi
	var notifs []models.Notification
	config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifs)

	unreadCount := 0
	for _, n := range notifs {
		if !n.IsRead {
			unreadCount++
		}
	}

	userName, _ := c.Get("user_name")
	c.HTML(http.StatusOK, "user/loan_form.html", gin.H{
		"UserName":    userName,
		"Items":       items,
		"Notifs":      notifs,
		"UnreadCount": unreadCount,
	})
}

// 3. TAMPILKAN RIWAYAT PEMINJAMAN USER
func HistoryLoan(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(uint)

	// --- OTOMATIS CEK DENDA ---
	CheckAndCreateLateNotifications(userID)
	// --------------------------

	var loans []models.Loan

	// Preload "Item" artinya: "Tolong ambilkan juga data nama barangnya dari tabel items"
	config.DB.Preload("Item").Where("user_id = ?", userID).Order("borrow_date desc").Find(&loans)

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

	userName, _ := c.Get("user_name")
	c.HTML(http.StatusOK, "user/history.html", gin.H{
		"UserName":    userName,
		"Loans":       loans,
		"Notifs":      notifs,
		"UnreadCount": unreadCount,
	})
}

func StoreLoan(c *gin.Context) {
	// 1. Ambil User ID
	userID, _ := c.MustGet("user_id").(uint)

	// --- BARU: Validasi Blokir Pinjaman ---
	blocked, reason := IsUserBlocked(userID)
	if blocked {
		c.String(http.StatusForbidden, "Peminjaman ditolak: "+reason)
		return
	}

	// 2. Ambil Input dari Modal
	itemID, _ := strconv.Atoi(c.PostForm("item_id"))
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	notes := c.PostForm("notes")

	// --- BARU: Kelola Folder Per User ---
	userFolder := "uploads/user_" + strconv.FormatUint(uint64(userID), 10)
	os.MkdirAll(userFolder, os.ModePerm)

	// --- BARU: Handle Upload KTM ---
	file, err := c.FormFile("student_id_card")
	var ktmFilename string
	if err == nil {
		// Simpan KTM dengan nama unik di folder user
		filename := "ktm_" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename
		c.SaveUploadedFile(file, userFolder+"/"+filename)
		ktmFilename = "user_" + strconv.FormatUint(uint64(userID), 10) + "/" + filename
	} else {
		// Jika KTM wajib, bisa return error di sini.
		// Untuk sekarang kita asumsikan wajib jika ingin kredibilitas tinggi.
		c.String(http.StatusBadRequest, "Kartu Mahasiswa (KTM) wajib diunggah!")
		return
	}

	// 3. Hitung Tanggal
	borrowDate := time.Now()                      // Tanggal pinjam = HARI INI
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
		UserID:        userID,
		ItemID:        uint(itemID),
		BorrowDate:    borrowDate,
		DueDate:       dueDate,
		Status:        "dipinjam",
		Notes:         notes,
		StudentIDCard: ktmFilename, // Simpan KTM
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

	// (Opsional) Cek apakah terlambat & Hitung Denda Final
	if now.After(loan.DueDate) {
		diff := now.Sub(loan.DueDate)
		daysLate := int(diff.Hours() / 24)
		if daysLate < 1 {
			daysLate = 1
		}

		totalFine := 0
		if daysLate >= 4 {
			weeksLate := ((daysLate - 4) / 7) + 1
			totalFine = weeksLate * 20000
		}
		loan.FineAmount = totalFine
		loan.Notes = loan.Notes + " [Kembali Terlambat: " + strconv.Itoa(daysLate) + " hari]"
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

// 7. MAHASISWA KONFIRMASI PEMBAYARAN
func ConfirmPayment(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(uint)
	loanID := c.PostForm("loan_id")
	method := c.PostForm("payment_method")

	// --- BARU: Handle Upload File ---
	file, err := c.FormFile("payment_proof")
	var filename string
	if err == nil {
		// Tentukan folder user
		userFolder := "uploads/user_" + strconv.FormatUint(uint64(userID), 10)
		os.MkdirAll(userFolder, os.ModePerm)

		// Simpan file dengan nama unik
		saveName := "proof_" + loanID + "_" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename
		c.SaveUploadedFile(file, userFolder+"/"+saveName)
		filename = "user_" + strconv.FormatUint(uint64(userID), 10) + "/" + saveName
	}

	var loan models.Loan
	if err := config.DB.First(&loan, loanID).Error; err != nil {
		c.String(404, "Data tidak ditemukan")
		return
	}

	// SECURITY: Pastikan yang bayar adalah pemilik peminjaman ini
	if loan.UserID != userID {
		c.String(403, "Anda tidak memiliki akses!")
		return
	}

	loan.PaymentMethod = method
	loan.PaymentProof = filename // Simpan nama file saja
	loan.Notes = loan.Notes + " [Bukti Diunggah: " + method + "]"
	config.DB.Save(&loan)

	c.Redirect(http.StatusFound, "/user/history")
}

// 8. ADMIN VERIFIKASI PEMBAYARAN
func VerifyPayment(c *gin.Context) {
	loanID := c.PostForm("loan_id")

	var loan models.Loan
	if err := config.DB.First(&loan, loanID).Error; err != nil {
		c.String(404, "Data tidak ditemukan")
		return
	}

	loan.IsFinePaid = true
	loan.Notes = loan.Notes + " [Denda Lunas]"
	config.DB.Save(&loan)

	c.Redirect(http.StatusFound, "/admin/reports")
}

// 9. FUNGSI PEMBANTU: CEK APAKAH USER DIBLOKIR PINJAM (VERSI RINGAN)
func IsUserBlocked(userID uint) (bool, string) {
	var lateLoan models.Loan
	now := time.Now()

	// Hanya BLOKIR jika barang belum kembali dan sudah lewat 3 hari (masuk hari ke-4)
	errLate := config.DB.Where("user_id = ? AND status = ? AND due_date < ?", userID, "dipinjam", now).First(&lateLoan).Error
	if errLate == nil {
		diff := now.Sub(lateLoan.DueDate)
		daysLate := int(diff.Hours() / 24)
		if daysLate >= 3 { // Masuk hari ke-4 (karena 3 hari pertama masa tenggang)
			return true, "Peminjaman dikunci karena Anda terlambat mengembalikan barang lebih dari 3 hari."
		}
	}

	// Untuk denda yang belum bayar tapi barang sudah kembali -> JANGAN BLOKIR (Biar mahasiswa bisa pinjam barang lain)
	return false, ""
}

// 6. FUNGSI PEMBANTU: CEK & BUAT NOTIFIKASI DENDA OTOMATIS
func CheckAndCreateLateNotifications(userID uint) {
	var lateLoans []models.Loan
	now := time.Now()

	// Cari pinjaman yang: milih user ini, status 'dipinjam', dan sudah melewati DueDate
	config.DB.Preload("Item").
		Where("user_id = ? AND status = ? AND due_date < ?", userID, "dipinjam", now).
		Find(&lateLoans)

	finePerWeek := 20000 // Rp 20.000 per minggu

	for _, loan := range lateLoans {
		// --- BARU: Skip jika denda sudah dibayar ---
		if loan.IsFinePaid {
			continue
		}

		// 1. Hitung Selisih Hari
		diff := now.Sub(loan.DueDate)
		daysLate := int(diff.Hours() / 24)
		if daysLate < 1 {
			daysLate = 1 // Minimal 1 hari jika sudah melewati batas jam
		}

		totalFine := 0
		statusNotif := "PERINGATAN"

		// 2. Logika Denda: Hari 1-3 Gratis, Hari 4 mulai denda mingguan
		if daysLate >= 4 {
			// Perhitungan mingguan: hari ke 4-10 = 1 minggu, 11-17 = 2 minggu, dst.
			weeksLate := ((daysLate - 4) / 7) + 1
			totalFine = weeksLate * finePerWeek
			statusNotif = "DENDA"
		}

		// 3. Update Nilai Denda di Database
		loan.FineAmount = totalFine
		config.DB.Save(&loan)

		// 4. Susun Pesan
		var pesan string
		if statusNotif == "PERINGATAN" {
			pesan = "⚠️ [PERINGATAN] Aset [" + loan.Item.Name + "] belum dikembalikan. Batas waktu: " + loan.DueDate.Format("02 Jan 2006") + ". Hari ke-" + strconv.Itoa(daysLate) + " telat (Masih masa tenggang)."
		} else {
			pesan = "🚨 [DENDA] Aset [" + loan.Item.Name + "] terlambat " + strconv.Itoa(daysLate) + " hari. Denda berjalan: Rp " + strconv.Itoa(totalFine) + " (Dihitung per minggu). Mohon segera kembalikan!"
		}

		// 5. Cek Notifikasi di Database
		var existingNotif models.Notification
		err := config.DB.Where("user_id = ? AND loan_id = ? AND type = ?", userID, loan.ID, "denda").First(&existingNotif).Error

		if err != nil { // Jika belum ada, buat baru
			newNotif := models.Notification{
				UserID:  userID,
				LoanID:  loan.ID,
				Type:    "denda",
				Message: pesan,
				IsRead:  false,
			}
			config.DB.Create(&newNotif)
		} else {
			// Jika sudah ada, update pesannya agar nominal terbaru & status peringatan/denda terupdate
			existingNotif.Message = pesan
			config.DB.Save(&existingNotif)
		}
	}
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

// 10. MAHASISWA TANDAI NOTIFIKASI SUDAH DIBACA
func MarkNotificationRead(c *gin.Context) {
	notifID := c.PostForm("notif_id")
	userID, _ := c.MustGet("user_id").(uint)

	var notif models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notifID, userID).First(&notif).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	notif.IsRead = true
	config.DB.Save(&notif)

	// Kita balikkan JSON sukses saja karena ini biasanya dipanggil via AJAX/JS
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi dibaca"})
}
