# Changelog 📜

Riwayat pemeliharaan dan perubahan fitur utama pada aplikasi Pinjam Aset Kampus.

## [1.4.0] - 2026-01-04
### Added
- **Landing Page**: Halaman pendaratan baru dengan hero section, carousel floating cards, proses workflow, dan fitur unggulan dengan desain modern.
- **Interactive Carousel**: Floating cards carousel dengan navigasi prev/next dan dot indicators di landing page untuk menampilkan statistik pengguna.
- **Background Decorations**: Sistem dekorasi background dengan animated blur shapes, floating icons, twinkling particles, dan geometric shapes dengan smooth animations.
- **Green Theme**: Migrasi lengkap dari blue (#3b82f6) ke green (#10b981) untuk konsistensi visual di semua halaman.
- **CSS External Files**: Ekstraksi semua inline CSS ke file eksternal di folder `static/css/` dengan struktur terorganisir per halaman.
- **Static File Serving**: Setup routing `/css`, `/js`, dan `/uploads` di main.go untuk serving aset statis.
- **Auth Toggle Links**: Penambahan link "Belum punya akun? Daftar di sini" di login dan "Sudah punya akun? Login di sini" di register untuk navigasi antar halaman.
- **Auto-hash Password**: Implementasi bcrypt password hashing otomatis di endpoint register (tanpa perlu hash.go manual).
- **Modal Z-index Fix**: Perbaikan modal bayar denda di history page dengan CSS z-index yang tepat agar modal muncul di atas backdrop.
- **Git Ignore Update**: Penambahan `uploads/` dan `public/` ke .gitignore untuk mengecualikan file lokal dari tracking.

### Changed
- **Landing Page Data**: Statistik pengguna diubah dari unrealistic (500+, 1.2K+) ke realistic (120+, 350+, 24/7).
- **Color Scheme**: Migrasi color theme dari blue ke green di seluruh aplikasi:
  - Primary: #3b82f6 → #10b981 (emerald green)
  - Primary Dark: #2563eb → #059669 (darker green)
  - Background Gradients: Blue tones → green tones (#f0fdf4, #dcfce7, #bbf7d0)
  - All shadows, borders, hover states → green rgba(16, 185, 129, x)
- **Main.go**: Pembersihan code:
  - Hapus route `/buat-akun-test` (testing route dengan hardcode password)
  - Hapus import bcrypt (sudah di auth_controller.go)
  - Hapus hardcode data landing page (hanya pass data yang diperlukan)
  - Remove `/static` route duplikat, optimize static routes
- **Hash.go**: File dihapus karena password hashing sekarang langsung di auth_controller.go pada register endpoint.
- **Footer Text Styling**: Teks "Login di sini" dan "Daftar di sini" dibesarkan (0.875rem → 1rem) dan dimudahkan dibaca dengan font-weight 600.
- **Link Colors**: Footer links sekarang menggunakan green theme dengan hover effect ke darker green.

### Improved
- **Code Organization**: CSS sekarang modular dan mudah dimaintain dengan struktur folder `static/css/{admin,auth,user}/`.
- **CSS Performance**: Pengurangan ukuran HTML templates dengan memindahkan styling ke CSS eksternal.
- **Responsive Design**: Landing page responsif untuk mobile, tablet, dan desktop dengan mobile-first approach.
- **Animation Quality**: Carousel dan background animations smooth dengan timing yang tepat (20-35s duration dengan staggered delays).
- **Accessibility**: Link dan buttons memiliki proper color contrast dan hover states untuk better UX.
- **Decorative Elements**: Background decorations tidak mengganggu content readability dengan opacity dan positioning yang tepat.

### Fixed
- **Modal Backdrop Issue**: Modal payment di history page tidak lagi tersembunyi di balik backdrop dengan z-index fix di CSS.
- **CSS Link Path**: Perbaikan path CSS dari `/static/css/` ke `/css/` untuk match dengan routing di main.go.
- **Inline Style Conflict**: Hilangkan inline styles di HTML auth pages dan gunakan CSS classes untuk avoid specificity issues.
- **Duplicate CSS Rules**: Hapus duplicate `.footer-text` rules yang menyebabkan styling conflict.
- **Broken Keyframes**: Perbaikan fragmen CSS yang broken (stray `@keyframes` blocks) di auth CSS files.
- **File Upload Button**: Styling diperbaiki untuk browser compatibility (Chrome, Firefox, Safari) dengan vendor prefixes.

### Removed
- **hash.go**: File yang sudah tidak digunakan (password hashing pindah ke auth_controller.go)
- **Route /buat-akun-test**: Testing route dihapus untuk production readiness
- **Hardcoded Landing Data**: Hapus data statis dari main.go route handler
- **Inline CSS**: Semua CSS dari HTML tags sudah dipindah ke external files

## [1.3.0] - 2025-12-31
### Added
- **Admin Welcome Banner**: Penambahan welcome banner animasi dengan gradient biru di dashboard admin, menampilkan greeting personal dan real-time clock dalam bahasa Indonesia.
- **Login Info Panel**: Implementasi side-by-side layout pada halaman login dengan info panel (kiri) berisi 3 feature cards dan login form (kanan).
- **User Profile Card**: Penambahan user profile card dengan avatar gradient di semua halaman admin (dashboard, items, loans, reports).
- **About Modal Enhancement**: Link eksternal ke CHANGELOG.md di GitHub pada modal "Tentang Project" untuk akses mudah ke riwayat perubahan.
- **Statistics Cards (User History)**: Penambahan 3 cards statistik di halaman riwayat user: Total Riwayat (biru), Dikembalikan (hijau), dan Dipinjam (kuning) dengan icon gradient dan badges.
- **Search & Filter (User History)**: Implementasi search bar dan dropdown filter status dengan real-time filtering di halaman riwayat peminjaman user.
- **Gradient Background (User Pages)**: Penambahan gradient background dengan decorative floating circles (biru & ungu) untuk tampilan lebih modern dan dinamis.
- **Admin Items Deskripsi**: Kolom deskripsi kini tampil di daftar Data Barang agar detail item terlihat tanpa membuka form.

### Changed
- **Light Theme Migration**: Perubahan tema sidebar admin dari dark (#0f172a → #1e293b) ke light gradient (#f0f4f8 → #e8eef5) untuk kenyamanan mata.
- **Login Background**: Update gradient background login dari purple (#667eea → #764ba2) ke blue (#3b82f6 → #2563eb) untuk konsistensi tema.
- **Admin Headers**: Standardisasi header di halaman Data Barang, Riwayat Peminjaman, dan Laporan dengan title + subtitle deskriptif.
- **Modal Design**: Update modal "Tentang Aplikasi" dengan gradient biru, icon boxes, dan versi 1.3.0 di semua halaman admin dan user.
- **User Header Gradient**: Perubahan header user pages dari dark gradient (#1e293b) ke blue gradient (#3b82f6 → #2563eb) yang matching dengan primary theme.
- **Background Enhancement**: Update dari solid color (#ebf4ff) ke 3-color gradient dengan floating decorative shapes untuk visual depth.
- **Admin Item Form Layout**: Form tambah/edit kini terpusat dengan background putih, badge mode putih, dan tombol kembali dipangkas (hanya Batal di footer).
- **Loans/Reports Headers**: Judul kartu Riwayat dan Laporan diringkas untuk tampilan lebih sederhana.

### Improved
- **Dashboard Cards Styling**: Enhancement styling card statistik dengan decorative circles, smooth hover effects, dan gradient colors yang konsisten.
- **Nav Links Colors**: Optimasi warna nav-link dari #94a3b8 ke #64748b dengan hover state rgba(37, 99, 235, 0.1) untuk kontras lebih baik.
- **Responsive Design**: Perbaikan responsivitas login page untuk tablet dan mobile devices.
- **Card Sizes**: Pengurangan ukuran logo (70px → 60px) dan padding card (3rem → 2rem) untuk tampilan lebih compact.
- **Tech Stack Info**: Penambahan informasi teknologi (Go Gin Framework & PostgreSQL) pada modal project information.
- **Hover Consistency**: Standardisasi semua hover scale effects di admin pages dari berbagai nilai (1.05, 1.1) ke 1.02 untuk animasi yang lebih smooth dan profesional.
- **Button Hover Effects**: Unifikasi transform scale pada buttons, badges, action buttons, table rows, KTM links, dan PDF buttons.

### Fixed
- **Logout Button Styling**: Perbaikan warna logout button dari #fca5a5 ke #dc2626 dengan hover effect yang lebih jelas.
- **Border Colors**: Update border dari rgba(255,255,255,0.05) ke rgba(37,99,235,0.1) untuk visibility lebih baik pada light theme.
- **Logout Hover (User)**: Penambahan hover state merah (#fef2f2 background) untuk dropdown logout button di semua halaman user (dashboard, history, loan form).
- **Template Error**: Fix template function error dengan mengganti function "add" menjadi JavaScript client-side counting untuk statistics cards.

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
