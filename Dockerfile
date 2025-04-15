# Stage 1: Build
# Gunakan image resmi Go versi 1.21 sebagai builder
FROM golang:1.24-alpine as builder

# Set direktori kerja di dalam container ke /app
WORKDIR /app

# Salin semua file dari direktori proyek lokal ke dalam container (ke /app)
COPY . .

# Mengunduh semua dependency Go yang didefinisikan di go.mod
RUN go mod download

# Compile file Go menjadi binary dengan nama "server"
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go


# Stage 2: Run
# Menggunakan image alpine untuk menjalankan aplikasi
FROM alpine:latest

# Set direktori kerja di container runtime ke /root/
WORKDIR /app

# Tambahkan user non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Salin binary hasil build (server) dari stage builder ke stage ini
COPY --from=builder /app/server .

# Pastikan binary bisa dieksekusi oleh appuser
RUN chmod 755 /app/server

# Ganti user ke appuser
USER appuser

# Mendeklarasikan bahwa aplikasi berjalan pada port 8080
EXPOSE 8080

# menjalan file binary tersebut
CMD ["./server"]
