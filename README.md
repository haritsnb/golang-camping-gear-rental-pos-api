
# Camping Gear Rental Point of Sale (POS) API

Backend API untuk sistem kasir dan manajemen penyewaan perlengkapan camping. Aplikasi ini dirancang untuk membantu pengelolaan data alat, stok, pelanggan, transaksi penyewaan, pengembalian barang, pembayaran, serta manajemen pengguna dan hak akses.


## A. Panduan Setup & Menjalankan Project

### 1. Clone Repository

Clone project dari repository

```bash
git clone https://github.com/haritsnb/golang-camping-gear-rental-pos-api.git
cd golang-camping-gear-rental-pos-api
```

### 2. Konfigurasi File `.env`

Salin file `.env.example` menjadi `.env` di root folder project

```bash
cp .env.example .env
```

Buka file `.env` dan sesuaikan kredensial konfigurasi Anda (Port aplikasi, Database, dll).

## B. Menjalankan Aplikasi
Pilih salah satu metode di bawah ini untuk menjalankan aplikasi:

### Opsi 1: Tanpa Docker

#### 1. Persiapan Database

Pastikan Anda sudah menginstal **PostgreSQL** di komputer Anda, membuat database baru, dan memastikan variabel di file `.env` mengarah ke database lokal Anda

```env
APP_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=5432
DB_SSL=disable
```

#### 2. Install Dependensi

Unduh seluruh package dan library Go yang dibutuhkan proyek

```bash
go mod tidy
```

#### 3. Jalankan Aplikasi

Eksekusi perintah berikut untuk menyalakan server backend

```bash
go run main.go
```

### Opsi 2: Menggunakan Docker Compose

#### 1. Sesuaikan Konfigurasi `.env` untuk Docker

Pastikan nilai `DB_HOST` diarahkan ke nama layanan kontainer database (`postgres-db`), bukan `127.0.0.1`

```env
APP_PORT=8080
DB_HOST=postgres-db
DB_PORT=5432
DB_SSL=disable
```

#### 2. Jalankan Docker Compose

Jalankan perintah berikut di direktori root proyek untuk membangun image dan mengaktifkan semua layanan di latar belakang (*detached mode*)

```bash
docker compose up -d --build
```

#### 3. Manajemen Kontainer (Opsional)

- Melihat log sistem aplikasi Go: `docker compose logs -f go-app`
- Melihat status kontainer aktif: `docker compose ps`
- Mematikan seluruh sistem kontainer: `docker compose down`


## C. Dokumentasi Swagger UI

### 1. Buka di Browser

Setelah aplikasi berjalan (baik via Opsi 1 maupun Opsi 2), buka browser Anda dan akses dokumentasi API interaktif pada tautan berikut

```url
http://localhost:8080/api/docs/index.html
```

### 2. Cara Login & Otentikasi di Swagger UI

#### 1. Mendapatkan Token Akses:
- Pada halaman Swagger UI, cari grup rute **`Auth`** dan buka **`POST /auth/login`**.
- Klik tombol **`Try it out`**.
- Gunakan request body bawaan (atau sesuaikan dengan akun admin Anda)

```json
{
  "username": "admin",
  "password": "password123"
}
```

- Klik tombol **`Execute`**.
- Salin string nilai token yang muncul di dalam *Server response* (salin teksnya saja tanpa tanda petik **`"`**).

#### 2. Mengaktifkan Akses Global (Authorize):
- Gulir kembali ke bagian paling atas halaman Swagger UI.
- Klik tombol hijau **`Authorize 🔓`** di pojok kanan atas.
- Pada kolom input **Value**, ketik kata `Bearer`  diikuti dengan spasi, lalu tempel (*paste*) token Anda. Contoh formatnya

```text
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```
- Klik tombol **`Authorize`**, kemudian klik **`Close`**.
- Ikon gembok pada rute-rute API terproteksi akan berubah menjadi terkunci (`🔒`), menandakan otentikasi berhasil.

## D. API Tester & Simulator Client-Side

Untuk melakukan simulasi pengujian API langsung dari sisi *client-side*, Anda dapat mengakses dashboard simulator internal di

```url
http://localhost:8080/test
```
