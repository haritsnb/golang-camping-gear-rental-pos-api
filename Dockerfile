# ==========================================
# STAGE 1: Build / Kompilasi Aplikasi
# ==========================================
FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Menyalin SEMUA file dan folder dari komputer lokal ke kontainer builder
COPY . .

# Melakukan kompilasi aplikasi Go menjadi file biner bernama "main"
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ==========================================
# STAGE 2: Image Produksi Akhir (Otomatis & Ringkas)
# ==========================================
FROM alpine:3.18
WORKDIR /app

# Menyalin SEMUA file dari Stage 1 (termasuk folder docs, test, storages, dll) ke image akhir
COPY --from=builder /app .

# Memastikan folder 'storages' untuk upload KYC tetap ada/dibuat jika belum ada di lokal
RUN mkdir -p storages

EXPOSE 8080

# Menjalankan file biner hasil compile
CMD ["./main"]
