package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/controllers"
	"pinjam_aset_kampus/middleware"
	"pinjam_aset_kampus/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDatabase()

	// 1b. Auto Migrate (Sinkronisasi Tabel Otomatis)
	config.DB.AutoMigrate(&models.User{}, &models.Item{}, &models.Loan{}, &models.Notification{})

	// 2. Init Router Gin
	r := gin.Default()

	// 3. Load semua file HTML dari folder views (rekursif)
	var htmlFiles []string
	filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".html") {
			htmlFiles = append(htmlFiles, path)
		}
		return nil
	})
	tmpl := template.Must(template.ParseFiles(htmlFiles...))
	r.SetHTMLTemplate(tmpl)
	r.Static("/css", "./static/css")  // Static CSS files
	r.Static("/js", "./static/js")    // Static JS files
	r.Static("/uploads", "./uploads") // User uploaded files (payment proofs, etc)

	// --- A. ROUTE PUBLIK (Tanpa Login) ---

	// Landing page (public)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "landing.html", gin.H{
			"ctaLogin":    "/login",
			"ctaRegister": "/register",
			"fine":        "Rp 20.000 per minggu mulai hari ke-4 keterlambatan",
			"gracePeriod": "3 hari pertama bebas denda",
			"blockNotice": "Pinjaman dikunci jika telat > 3 hari dengan status masih dipinjam",
		})
	})

	// Halaman Auth
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)

	r.GET("/register", controllers.ShowRegisterPage)
	r.POST("/register", controllers.Register)

	r.GET("/logout", controllers.Logout)

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
