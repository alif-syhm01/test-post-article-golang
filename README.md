# Test Post Article

Technical test: REST API Article dengan Golang.

## Struktur Proyek

- `main.go` - Entry Point, setup database, migration, route, menjalankan server
- `handler/` - HTTP handler (request/respon)
- `model/` - Struct dari Post
- `repositories/` - Query Database
- `router/` - Custom router dengan dynamic route
- `services/` - Business logic & validation
- `create-table-posts.sql` - Task no 1, pembuatan tabel baru di database article dengan nama "posts"
- `Post Article.postman_collection.json` - Postman Collection untuk testing

## Cara Menjalankan

1. Pastikan MySQL berjalan sesuai di url dan port nya masing masing (untuk local, pakai localhost:3306 (default))
2. Buat database `article` terlebih dahulu menggunakan SQL query
3. Jalankan `go mod tidy` pada terminal
4. Jalankan `go run main.go` pada terminal (sudah include migrasi tabel, koneksi database, dan menjalankan server)
5. Kalau berhasil server akan berjalan di `http://localhost:8080` (berlaku untuk local)

## Endpoint

| Method | URL | Deskripsi |
|---|---|---|
| GET | /api/v1/articles | Mengambil semua data artikel |
| GET | /api/v1/articles/{limit}/{offset} | Mengambil semua data artikel berdasarkan limit dan offset (pagination) |
| POST | /api/v1/articles | Membuat artikel baru |
| GET | /api/v1/articles/{id} | Mengambil data artikel berdasarkan Id |
| PUT | /api/v1/articles/{id} | Mengupdate data artikel berdasarkan Id |
| DELETE | /api/v1/articles/{id} | Menghapus data artikel berdasarkan Id |