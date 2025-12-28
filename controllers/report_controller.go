package controllers

import (
	"net/http"
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/models"

	"github.com/gin-gonic/gin"
)

// A. HALAMAN RIWAYAT (MANAJEMEN - Ada Tombol Aksi)
func IndexLoans(c *gin.Context) {
	var loans []models.Loan
	config.DB.Preload("User").Preload("Item").Order("created_at desc").Find(&loans)

	c.HTML(http.StatusOK, "admin/loans.html", gin.H{
		"Loans": loans,
	})
}

// B. HALAMAN LAPORAN (ANALISA - Filter User & Bulan)
func IndexReports(c *gin.Context) {
	var loans []models.Loan
	var users []models.User

	// Query Dasar
	query := config.DB.Preload("User").Preload("Item").Order("created_at desc")

	// 1. Filter: Per Mahasiswa
	userID := c.Query("user_id")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 2. Filter: Per Bulan
	month := c.Query("month")
	if month != "" {
		// PostgreSQL syntax untuk ambil bulan dari tanggal
		query = query.Where("TO_CHAR(borrow_date, 'YYYY-MM') = ?", month)
	}

	query.Find(&loans)

	// Ambil data user untuk dropdown filter
	config.DB.Where("role = ?", "user").Find(&users)

	c.HTML(http.StatusOK, "admin/reports.html", gin.H{
		"Loans":         loans,
		"Users":         users,
		"SelectedUser":  userID,
		"SelectedMonth": month,
	})
}