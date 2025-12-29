# Pinjam Aset Kampus 🏢

Aplikasi Peminjaman Aset Kampus berbasis Web yang dibangun menggunakan **Golang (Gin)** dan **PostgreSQL**. Aplikasi ini dirancang untuk mempermudah mahasiswa dalam meminjam aset kampus serta membantu admin dalam mengelola logistik dan denda secara otomatis.

> [!NOTE]  
> Proyek ini adalah fork dari [Yeflou](https://github.com/Yeflou) dan dikembangkan kembali oleh **Abduldinata (Tim 3)** untuk kebutuhan UAS Cloud Computing.

## ✨ Fitur Utama

- **📦 Manajemen Aset**: Inventaris barang dengan pelacakan stok real-time.
- **🕒 Sistem Denda Otomatis**: Perhitungan denda keterlambatan secara otomatis (Skema mingguan dengan masa tenggang).
- **📝 Form Peminjaman Digital**: Proses peminjaman dengan syarat unggah KTM (Kartu Mahasiswa).
- **💳 Konfirmasi Pembayaran**: Fitur unggah bukti bayar denda bagi mahasiswa dan verifikasi oleh admin.
- **🔔 Sistem Notifikasi**: Pemberitahuan otomatis terkait status pinjaman dan tagihan denda.
- **📊 Laporan Admin**: Rekap peminjaman dan filter laporan yang mendetail.
- **💎 UI Premium**: Interface modern dengan Glassmorphism dan sidebar/navbar yang konsisten.

## 🚀 Setup Project

### Prasyarat
- [Go](https://golang.org/dl/) (Minimal versi 1.25.1)
- [PostgreSQL](https://www.postgresql.org/download/)
- Git

### Langkah Instalasi

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd pinjam_aset_kampus
   ```

2. **Konfigurasi Database**
   - Buat database baru di PostgreSQL (contoh: `pinjam_aset`).
   - Buka **Query Tool** di pgAdmin pada database tersebut, lalu **Open** dan **Execute** file `db_peminjaman_kampus.sql` untuk membuat tabel dan data awal secara otomatis.
   - Buat file `.env` di root folder dan sesuaikan variabel berikut:
     ```env
     DB_HOST=localhost
     DB_USER=postgres
     DB_PASSWORD=yourpassword
     DB_NAME=pinjam_aset
     DB_PORT=5432
     ```
   - *Tip: Pastikan password database Anda sesuai dengan yang ada di properti PostgreSQL.*

3. **Install Dependensi**
   ```bash
   go mod tidy
   ```

4. **Jalankan Aplikasi**
   ```bash
   go run main.go
   ```
   Aplikasi akan berjalan di `http://localhost:8080`.

5. **Akun Testing Default**
   Anda sudah bisa langsung login menggunakan data dari `db_peminjaman_kampus.sql`:
   - **Admin**: `super@admin.com` / `admin123`
   - **User**: `budi@mhs.ac.id` / `mhs123`

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Web Framework**: Gin Gonic
- **ORM**: GORM
- **Database**: PostgreSQL
- **Frontend**: HTML5, CSS3 (Vanilla), Bootstrap 5, FontAwesome 6

## 📄 Lisensi

Proyek ini dilisensikan di bawah **MIT License**. Lihat file [LICENSE](LICENSE) untuk informasi lebih lanjut.

## 👥 Tim Pengembang
- **Yeflou** - Owner Project & Forked from
- **Abduldinata** - Lead Maintenance & Modified
- **Tim 3** - UAS Cloud Computing Project

---
*Dibuat dengan ❤️ untuk sistem kampus yang lebih baik.*
