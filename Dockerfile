# Gunakan image resmi PostgreSQL 16 sebagai basis
# Menggunakan versi spesifik adalah praktik yang baik untuk lingkungan produksi
FROM postgres:16

# Tidak ada lagi ENV di sini. Konfigurasi akan diinjeksikan oleh Docker Compose saat runtime.

# (Opsional) Anda bisa menyalin skrip inisialisasi jika diperlukan
# COPY ./init.sql /docker-entrypoint-initdb.d/init.sql

EXPOSE 3002