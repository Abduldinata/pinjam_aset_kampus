# Changelog 📜

Riwayat pemeliharaan dan perubahan fitur utama pada aplikasi Pinjam Aset Kampus.

## [1.3.0] - 2025-12-31
### Added
- **Admin Welcome Banner**: Penambahan welcome banner animasi dengan gradient biru di dashboard admin, menampilkan greeting personal dan real-time clock dalam bahasa Indonesia.
- **Login Info Panel**: Implementasi side-by-side layout pada halaman login dengan info panel (kiri) berisi 3 feature cards dan login form (kanan).
- **User Profile Card**: Penambahan user profile card dengan avatar gradient di semua halaman admin (dashboard, items, loans, reports).
- **About Modal Enhancement**: Link eksternal ke CHANGELOG.md di GitHub pada modal "Tentang Project" untuk akses mudah ke riwayat perubahan.

### Changed
- **Light Theme Migration**: Perubahan tema sidebar admin dari dark (#0f172a → #1e293b) ke light gradient (#f0f4f8 → #e8eef5) untuk kenyamanan mata.
- **Login Background**: Update gradient background login dari purple (#667eea → #764ba2) ke blue (#3b82f6 → #2563eb) untuk konsistensi tema.
- **Admin Headers**: Standardisasi header di halaman Data Barang, Riwayat Peminjaman, dan Laporan dengan title + subtitle deskriptif.
- **Modal Design**: Update modal "Tentang Aplikasi" dengan gradient biru, icon boxes, dan versi 1.3.0 di semua halaman admin.

### Improved
- **Dashboard Cards Styling**: Enhancement styling card statistik dengan decorative circles, smooth hover effects, dan gradient colors yang konsisten.
- **Nav Links Colors**: Optimasi warna nav-link dari #94a3b8 ke #64748b dengan hover state rgba(37, 99, 235, 0.1) untuk kontras lebih baik.
- **Responsive Design**: Perbaikan responsivitas login page untuk tablet dan mobile devices.
- **Card Sizes**: Pengurangan ukuran logo (70px → 60px) dan padding card (3rem → 2rem) untuk tampilan lebih compact.
- **Tech Stack Info**: Penambahan informasi teknologi (Go Gin Framework & PostgreSQL) pada modal project information.

### Fixed
- **Logout Button Styling**: Perbaikan warna logout button dari #fca5a5 ke #dc2626 dengan hover effect yang lebih jelas.
- **Border Colors**: Update border dari rgba(255,255,255,0.05) ke rgba(37,99,235,0.1) untuk visibility lebih baik pada light theme.

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
