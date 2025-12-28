package controllers

import (
	"net/http"
	"pinjam_aset_kampus/config" // SESUAIKAN
	"pinjam_aset_kampus/models" // SESUAIKAN
	"pinjam_aset_kampus/utils"  // SESUAIKAN

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 1. Tampilkan Halaman Login (GET)
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", gin.H{})
}

// 2. Proses Login (POST)
func Login(c *gin.Context) {
	// Ambil input dari form HTML
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User

	// Cek apakah email ada di database?
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{"Error": "Email tidak ditemukan!"})
		return
	}

	// Cek Password (Hash vs Input)
	// PENTING: Karena data dummy di DB passwordnya belum di-hash (masih plain text 'admin123'),
	// maka logic bcrypt ini akan gagal untuk data dummy lama.
	// Kita akan perbaiki datanya nanti.
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// SEMENTARA: Kalau mau login pakai data dummy yang belum di-hash, buka komentar baris bawah ini:
	// if user.Password != password { err = 1 } else { err = nil } // HANYA UNTUK TESTING DUMMY

	if err != nil {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{"Error": "Password salah!"})
		return
	}

	// Generate JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "auth/login.html", gin.H{"Error": "Gagal membuat token"})
		return
	}

	// Simpan Token di Cookie Browser (Supaya user tetap login saat pindah halaman)
	// Nama cookie: "token", expired: 3600*24 detik (1 hari)
	c.SetCookie("token", token, 3600*24, "/", "localhost", false, true)

	// Redirect sesuai Role
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.Redirect(http.StatusFound, "/user/dashboard")
	}
}

// 3. Logout
func Logout(c *gin.Context) {
	// Hapus cookie dengan set expired time ke masa lalu
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// --- TAMBAHAN UNTUK REGISTER ---

// 1. Tampilkan Halaman Register
func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{})
}

// 2. Proses Register User Baru
func Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validasi sederhana
	if name == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{"Error": "Semua kolom wajib diisi!"})
		return
	}

	// Hash Password (PENTING! Jangan simpan polos)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "auth/register.html", gin.H{"Error": "Gagal enkripsi password"})
		return
	}

	// Simpan ke Database
	// Catatan: Role otomatis di-set 'user'. Admin tidak bisa daftar lewat sini (harus lewat database/seeder)
	newUser := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		// Error biasanya karena email sudah terpakai
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{"Error": "Email sudah terdaftar!"})
		return
	}

	// Redirect ke halaman login setelah sukses
	c.HTML(http.StatusOK, "auth/login.html", gin.H{"Success": "Pendaftaran berhasil! Silakan login."})
}
