# Changelog 📜

Riwayat pemeliharaan dan perubahan fitur utama pada aplikasi Pinjam Aset Kampus.

## [1.2.0] - 2025-12-29
### Added
- **Premium Unified History List**: Overhaul tampilan riwayat peminjaman mahasiswa menjadi daftar tunggal mewah dengan efek hover baris dan icon dinamis.
- **Dynamic Icons Logic**: Penambahan logika JavaScript untuk menampilkan icon barang sesuai kategori (Laptop, Proyektor, Kamera, dll) secara otomatis.
- **Standardized Admin Sidebar**: Perbaikan posisi tombol logout yang kini ter-pin di bagian bawah sidebar secara konsisten di seluruh halaman admin.
- **Setup Documentation**: Penambahan file `database.sql`, `README.md`, `CHANGELOG.md`, dan `LICENSE` (MIT).

- **Footer Informasi**: Penambahan informasi footer yang kini ter-pin di bagian bawah halaman secara konsisten di seluruh halaman admin dan user.

### Changed
- **Database Migration**: Penambahan file `db_peminjaman_kampus.sql` untuk migrasi database awal.

### Improved
- **Footer Informasi**: memperbaiki supaya tampilan footer konsisten di seluruh halaman admin dan user.

### Fixed
- Perbaikan sinkronisasi warna tombol logout (tetap merah saat aktif/fokus).
- Penghapusan *horizontal scrolling* pada tabel riwayat mahasiswa untuk UX yang lebih bersih.
- Perbaikan visibilitas icon brand "Aset Kampus" di navbar premium.

## [1.1.0] - 2025-12-28
### Added
- **Sistem Denda Otomatis**: Implementasi denda mingguan dengan masa tenggang 3 hari yang dihitung saat pengembalian barang.
- **Validasi KTM**: Penambahan syarat unggah kartu mahasiswa saat melakukan peminjaman untuk meningkatkan kredibilitas data.
- **Konfirmasi Pembayaran Denda**: Fitur bagi mahasiswa untuk mengunggah bukti bayar (DANA/VA Bank) dan verifikasi manual oleh Admin.
- **Blokir Pinjaman**: Mekanisme penolakan peminjaman otomatis bagi user yang memiliki denda tertunggak atau telat mengembalikan barang.

### Improved
- **UI Premium Migration**: Perubahan skema desain dari gaya standar ke "Premium Glassmorphism" di dashboard mahasiswa dan admin.
- **Inline CSS Stability**: Migrasi kembali dari CSS eksternal ke internal (inline) atas permintaan user untuk kemudahan kustomisasi manual.

## [1.0.0] - 2025-12-20
### Initial Fork from Yeflou
- Struktur dasar aplikasi menggunakan Golang (Gin) dan PostgreSQL.
- Fitur dasar Login, Registrasi, Peminjaman Dasar, dan Kelola Barang.
- Implementasi awal GORM untuk auto-migration tabel.

---
*Maintained by Abduldinata (Tim 3).*
