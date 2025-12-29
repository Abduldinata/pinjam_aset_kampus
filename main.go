package main

import (
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/controllers"
	"pinjam_aset_kampus/middleware" // <-- Penting: Import middleware
	"pinjam_aset_kampus/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDatabase()

	// 1b. Auto Migrate (Sinkronisasi Tabel Otomatis)
	config.DB.AutoMigrate(&models.User{}, &models.Item{}, &models.Loan{}, &models.Notification{})

	// 2. Init Router Gin
	r := gin.Default()

	// 3. Load semua file HTML dari folder views
	r.LoadHTMLGlob("views/**/*")
	r.Static("/static", "./public")   // <-- BARU: Akses CSS/JS statis
	r.Static("/uploads", "./uploads") // <-- Akses foto bukti transfer

	// --- A. ROUTE PUBLIK (Tanpa Login) ---

	// Redirect root ke login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})

	// Halaman Auth
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)

	r.GET("/register", controllers.ShowRegisterPage)
	r.POST("/register", controllers.Register)

	r.GET("/logout", controllers.Logout)

	// [TESTING ONLY] Route buat admin & user cepat (Hapus kalau sudah production)
	r.GET("/buat-akun-test", func(c *gin.Context) {
		// Bikin Admin
		hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		config.DB.Create(&models.User{Name: "Super Admin", Email: "super@admin.com", Password: string(hashAdmin), Role: "admin"})

		// Bikin Mahasiswa
		hashMhs, _ := bcrypt.GenerateFromPassword([]byte("mhs123"), bcrypt.DefaultCost)
		config.DB.Create(&models.User{Name: "Budi Santoso", Email: "budi@mhs.ac.id", Password: string(hashMhs), Role: "user"})

		c.String(200, "Akun Admin (super@admin.com / admin123) & User (budi@mhs.ac.id / mhs123) berhasil dibuat!")
	})

	// --- B. ROUTE ADMIN ---
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin"))
	{
		admin.GET("/dashboard", controllers.DashboardAdmin)

		admin.GET("/items", controllers.IndexItems)
		admin.GET("/items/add", controllers.CreateItem)
		admin.POST("/items", controllers.StoreItem)

		// TAMBAHKAN INI:
		admin.GET("/loans", controllers.IndexLoans)

		admin.GET("/reports", controllers.IndexReports)
		admin.POST("/return", controllers.ReturnLoan)
		admin.POST("/remind", controllers.SendReminder)
		admin.POST("/verify-payment", controllers.VerifyPayment) // <-- BARU: Verifikasi Denda

		admin.GET("/items/edit", controllers.EditItem)      // Form Edit
		admin.POST("/items/update", controllers.UpdateItem) // Proses Update
		admin.POST("/items/delete", controllers.DeleteItem) // Proses Hapus// <--- Halaman Laporan
	}

	// --- C. ROUTE USER (Wajib Login sebagai User) ---
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware("user"))
	{
		user.GET("/dashboard", controllers.DashboardUser)

		// TAMBAHKAN INI:
		user.GET("/pinjam", controllers.CreateLoan)                // Form Pinjam
		user.POST("/pinjam", controllers.StoreLoan)                // Proses Pinjam
		user.GET("/history", controllers.HistoryLoan)              // Riwayat
		user.POST("/pay", controllers.ConfirmPayment)              // <-- BARU: Bayar Denda
		user.POST("/notif/read", controllers.MarkNotificationRead) // <-- BARU: Baca Notif
	}
	// 4. Jalankan Server
	r.Run(":8080")
}
