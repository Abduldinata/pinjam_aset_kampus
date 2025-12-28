package controllers // <--- PASTIKAN INI "controllers", JANGAN "main" ATAU "models"

import (
	"net/http"
	"pinjam_aset_kampus/config"
	"pinjam_aset_kampus/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 1. LIHAT SEMUA BARANG (IndexItems)
func IndexItems(c *gin.Context) {
	var items []models.Item
	config.DB.Find(&items)

	c.HTML(http.StatusOK, "admin/items.html", gin.H{
		"Items": items,
	})
}

// 2. TAMPILKAN FORM TAMBAH (CreateItem)
func CreateItem(c *gin.Context) {
	// Kirim data kosong agar template tidak error saat cek .Item.ID
	c.HTML(http.StatusOK, "admin/item_form.html", gin.H{
		"Item": models.Item{}, 
	})
}

// 3. PROSES SIMPAN BARANG BARU (StoreItem)
func StoreItem(c *gin.Context) {
	name := c.PostForm("name")
	category := c.PostForm("category")
	stockStr := c.PostForm("stock")
	location := c.PostForm("location")
	description := c.PostForm("description")

	stock, _ := strconv.Atoi(stockStr)

	newItem := models.Item{
		Name:        name,
		Category:    category,
		Stock:       stock,
		Location:    location,
		Description: description,
	}

	if err := config.DB.Create(&newItem).Error; err != nil {
		c.String(http.StatusInternalServerError, "Gagal menyimpan barang")
		return
	}

	c.Redirect(http.StatusFound, "/admin/items")
}

// 4. TAMPILKAN FORM EDIT (EditItem)
func EditItem(c *gin.Context) {
	id := c.Query("id") // Ambil dari ?id=...

	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	c.HTML(http.StatusOK, "admin/item_form.html", gin.H{
		"Item": item, // Kirim data lama untuk di-edit
	})
}

// 5. PROSES UPDATE (UpdateItem)
func UpdateItem(c *gin.Context) {
	id := c.PostForm("id")
	
	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	// Update data
	item.Name = c.PostForm("name")
	item.Category = c.PostForm("category")
	item.Location = c.PostForm("location")
	item.Description = c.PostForm("description")
	
	stockStr := c.PostForm("stock")
	item.Stock, _ = strconv.Atoi(stockStr)

	config.DB.Save(&item)

	c.Redirect(http.StatusFound, "/admin/items")
}

// 6. PROSES HAPUS (DeleteItem)
func DeleteItem(c *gin.Context) {
	id := c.PostForm("id")
	config.DB.Delete(&models.Item{}, id)
	c.Redirect(http.StatusFound, "/admin/items")
}