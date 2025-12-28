package controllers

import (
	"net/http"
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/models"

	"github.com/gin-gonic/gin"
)

// Dashboard Admin
func DashboardAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"Title": "Dashboard Admin",
		"Role":  "Admin",
	})
}

// Dashboard User (Mahasiswa)
func DashboardUser(c *gin.Context) {
	// 1. Ambil User ID & Role
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	role, _ := c.Get("role")

	// 2. Ambil Barang
	var items []models.Item
	config.DB.Where("stock > ?", 0).Find(&items)

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

	// 5. KIRIM DATA KE HTML (INI BAGIAN KRUSIAL)
	c.HTML(http.StatusOK, "user/dashboard.html", gin.H{
		"Title":       "Dashboard Mahasiswa",
		"Role":        role,
		"Items":       items,
		"Notifs":      notifs,
		"UnreadCount": unreadCount, // <--- PASTIIN INI ADA!
	})
}