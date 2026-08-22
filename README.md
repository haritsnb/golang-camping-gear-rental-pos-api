
# Camping Gear Rental Point of Sale (POS) API

Backend API untuk sistem kasir dan manajemen penyewaan perlengkapan camping. Aplikasi ini dirancang untuk membantu pengelolaan data alat, stok, pelanggan, transaksi penyewaan, pengembalian barang, pembayaran, serta manajemen pengguna dan hak akses.

## A. Panduan Setup & Menjalankan Project

### 1. Clone Repository

Clone project dari repository
```bash
git clone https://github.com/haritsnb/golang-camping-gear-rental-pos-api.git
```

### 2. Konfigurasi File `.env`

Cari file bernama **`.env.example`** yang berada tepat di root folder project, lalu ubah menjadi nama file menjadi **`.env`**, lalu jangan lupa sesuaikan isi file.

### 3. Install Dependensi

Jalankan perintah berikut di terminal untuk mengunduh semua package:

```bash
go mod tidy
```

### 4. Jalankan Aplikasi

Jalankan perintah berikut untuk meng-generate dokumentasi Swagger dan menyalakan server backend:

```bash
go run main.go
```

## B. Dokumentasi Swagger UI

### 1. Buka di Browser
Buka browser dan kunjungi URL:
```url
http://localhost:8081/api/docs/index.html
```
_*sesuaikan host dan port_

### 2. Cara Login & Otentikasi di Swagger UI:

1. **Login untuk Mendapatkan Token:**
   - Di Swagger UI, cari grup **`Auth`** --> klik **`POST /auth/login`**.
   - Klik tombol **`Try it out`**.
   - Pada request body, pastikan isinya:
     ```json
     {
       "username": "admin",
       "password": "password123"
     }
     ```
   - Klik tombol biru **`Execute`**.
   - Pada bagian *Server response*, salin string nilai **`token`** yang muncul (tanpa tanda kutip `"`).

2. **Memasang Token (Authorize):**
   - Gulir ke bagian paling atas halaman Swagger UI.
   - Klik tombol hijau **`Authorize 🔓`** di sebelah kanan atas.
   - Pada kolom input **Value**, ketik `Bearer ` lalu paste token Anda:
     ```text
     Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
     ```
   - Klik tombol **Authorize** --> klik **Close**.
   - Icon gembok akan berubah menjadi terkunci **🔒**.