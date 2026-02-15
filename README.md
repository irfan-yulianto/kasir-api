# Kasir API

Aplikasi Point of Sale (POS) / Kasir sederhana yang dibangun sebagai hasil pembelajaran **Golang Basic Bootcamp**. Project ini mencakup backend REST API menggunakan Go dan frontend web menggunakan React.

## Tentang Project

Project ini merupakan hasil dari bootcamp Golang dasar yang membahas:

- **Sesi 1-2**: Setup project Go, koneksi database PostgreSQL, CRUD dasar, dan clean architecture (repository-service-handler pattern)
- **Sesi 3**: Fitur transaksi checkout dan laporan penjualan harian
- **Sesi 4**: Middleware (API Key authentication, CORS, Logger) dan pembuatan frontend

Tujuan utama project ini adalah mempelajari konsep dasar pemrograman backend dengan Go, termasuk:

- Membangun REST API tanpa framework (menggunakan `net/http` standar)
- Menerapkan clean architecture dengan dependency injection
- Menulis raw SQL queries tanpa ORM
- Implementasi middleware pattern
- Integrasi frontend-backend

## Tech Stack

### Backend
| Teknologi | Keterangan |
|-----------|------------|
| Go 1.24 | Bahasa pemrograman utama |
| net/http | HTTP server (standard library, tanpa framework) |
| PostgreSQL | Database relasional |
| lib/pq | PostgreSQL driver untuk Go |
| Viper | Manajemen konfigurasi dari file `.env` |

### Frontend
| Teknologi | Keterangan |
|-----------|------------|
| React 19 | UI library |
| TypeScript | Type-safe JavaScript |
| Vite | Build tool dan dev server |
| Tailwind CSS 4 | Utility-first CSS framework |
| shadcn/ui | Komponen UI berbasis Radix UI |
| Lucide React | Icon library |

## Arsitektur

### Backend - Clean Architecture

```
Request → Middleware → Handler → Service → Repository → Database
```

| Layer | Tanggung Jawab |
|-------|---------------|
| **Middleware** | Cross-cutting concerns (auth, CORS, logging) |
| **Handler** | Menerima HTTP request, validasi input, kirim response |
| **Service** | Business logic dan aturan bisnis |
| **Repository** | Akses data, query SQL |
| **Model** | Struktur data / entity |
| **Config** | Konfigurasi aplikasi dari environment |
| **Database** | Koneksi dan manajemen database |

### Middleware
| Middleware | Fungsi |
|-----------|--------|
| `CORS` | Mengizinkan akses cross-origin dari frontend |
| `Logger` | Mencatat setiap request masuk dan durasi eksekusi |
| `APIKey` | Autentikasi sederhana via header `X-API-Key` |

### Struktur Project

```
kasir-api/
├── backend/
│   ├── config/             # Konfigurasi aplikasi (Viper)
│   │   └── config.go
│   ├── database/           # Koneksi database PostgreSQL
│   │   └── database.go
│   ├── handler/            # HTTP request handler
│   │   ├── category_handler.go
│   │   ├── product_handler.go
│   │   ├── report_handler.go
│   │   └── transaction_handler.go
│   ├── middleware/          # HTTP middleware
│   │   ├── api_key.go      # Autentikasi API Key
│   │   ├── cors.go         # Cross-Origin Resource Sharing
│   │   └── logger.go       # Request logging
│   ├── model/              # Data models / entities
│   │   ├── category.go
│   │   ├── product.go
│   │   ├── report.go
│   │   └── transaction.go
│   ├── repository/         # Data access layer (SQL queries)
│   │   ├── category_repository.go
│   │   ├── product_repository.go
│   │   ├── report_repository.go
│   │   └── transaction_repository.go
│   ├── service/            # Business logic layer
│   │   ├── category_service.go
│   │   ├── product_service.go
│   │   ├── report_service.go
│   │   └── transaction_service.go
│   ├── main.go             # Entry point, dependency injection, routing
│   ├── schema.sql          # Database schema dan seed data
│   ├── Dockerfile          # Container image
│   ├── go.mod
│   └── .env                # Environment variables (tidak di-track git)
├── frontend/
│   ├── src/
│   │   ├── components/ui/  # shadcn/ui components
│   │   ├── lib/
│   │   │   ├── api.ts      # API client (fetch wrapper)
│   │   │   └── utils.ts    # Helper functions
│   │   ├── App.tsx         # Aplikasi utama (split-view layout)
│   │   ├── main.tsx        # Entry point React
│   │   └── index.css       # Tailwind CSS + theme
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── .gitignore
└── README.md
```

## Database Schema

```sql
categories
├── id          SERIAL PRIMARY KEY
├── name        VARCHAR(255) NOT NULL
└── description TEXT

products
├── id          SERIAL PRIMARY KEY
├── name        VARCHAR(255) NOT NULL
├── price       INTEGER NOT NULL DEFAULT 0
├── stock       INTEGER NOT NULL DEFAULT 0
└── category_id INTEGER → categories(id) ON DELETE SET NULL

transactions
├── id           SERIAL PRIMARY KEY
├── total_amount INTEGER NOT NULL
└── created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP

transaction_details
├── id             SERIAL PRIMARY KEY
├── transaction_id INTEGER → transactions(id) ON DELETE CASCADE
├── product_id     INTEGER → products(id)
├── quantity       INTEGER NOT NULL
└── subtotal       INTEGER NOT NULL
```

## API Endpoints

### Public (tanpa API Key)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/categories` | Daftar semua kategori |
| POST | `/api/categories` | Buat kategori baru |
| GET | `/api/categories/{id}` | Detail kategori |
| PUT | `/api/categories/{id}` | Update kategori |
| DELETE | `/api/categories/{id}` | Hapus kategori |
| GET | `/api/categories/{id}/products` | Produk berdasarkan kategori |
| GET | `/api/products` | Daftar semua produk |
| GET | `/api/report/today` | Laporan penjualan hari ini |
| GET | `/health` | Health check |

### Protected (butuh header `X-API-Key`)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/products/{id}` | Detail produk |
| PUT | `/api/products/{id}` | Update produk |
| DELETE | `/api/products/{id}` | Hapus produk |
| POST | `/api/checkout` | Proses transaksi checkout |

### Contoh Request

**Get All Products**
```bash
curl http://localhost:8080/api/products
```

**Get Product Detail (protected)**
```bash
curl http://localhost:8080/api/products/1 \
  -H "X-API-Key: your-secret-api-key-here"
```

**Checkout (protected)**
```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "X-API-Key: your-secret-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": 1, "quantity": 2 },
      { "product_id": 2, "quantity": 3 }
    ]
  }'
```

**Laporan Penjualan**
```bash
# Hari ini
curl http://localhost:8080/api/report/today

# Range tanggal
curl "http://localhost:8080/api/report/today?start_date=2026-02-01&end_date=2026-02-15"
```

## Cara Menjalankan

### Prasyarat

- [Go](https://go.dev/dl/) >= 1.24
- [Node.js](https://nodejs.org/) >= 18
- [PostgreSQL](https://www.postgresql.org/) (atau gunakan layanan cloud seperti Supabase)

### 1. Clone Repository

```bash
git clone https://github.com/irfan-yulianto/kasir-api.git
cd kasir-api
```

### 2. Setup Database

Jalankan script SQL untuk membuat tabel dan seed data:

```bash
psql -U postgres -d your_database -f backend/schema.sql
```

### 3. Konfigurasi Environment

Buat file `backend/.env`:

```env
PORT=8080
DB_CONN=postgresql://username:password@localhost:5432/your_database
API_KEY=your-secret-api-key-here
```

### 4. Jalankan Backend

```bash
cd backend
go run main.go
```

Server berjalan di `http://localhost:8080`

### 5. Jalankan Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend berjalan di `http://localhost:5173`

### 6. (Opsional) Go Workspace

Jika menggunakan VS Code, buat `go.work` di root agar IDE mengenali module:

```bash
go work init ./backend
```

## Fitur Frontend

- Daftar produk dalam bentuk card grid
- Detail produk via dialog modal
- Edit produk (nama, harga, stok)
- Hapus produk dengan konfirmasi
- Keranjang belanja (sidebar kanan)
- Atur jumlah item (+/-)
- Proses checkout dengan struk transaksi
- Format mata uang Rupiah (IDR)
- Toast notification untuk feedback aksi
- Responsive layout (split-view: produk | keranjang)

## Docker

Build dan jalankan backend menggunakan Docker:

```bash
cd backend
docker build -t kasir-api .
docker run -p 8080:8080 --env-file .env kasir-api
```

## Yang Dipelajari

Konsep-konsep yang diterapkan dalam project ini:

1. **Go Fundamentals** - Struct, interface, error handling, package management
2. **HTTP Server** - Routing, request/response handling tanpa framework
3. **Clean Architecture** - Separation of concerns dengan layer yang jelas
4. **Dependency Injection** - Constructor-based DI tanpa library
5. **Database** - Koneksi PostgreSQL, raw SQL queries, prepared statements
6. **Middleware Pattern** - Function chaining untuk cross-cutting concerns
7. **REST API Design** - Resource-based URL, HTTP methods, status codes
8. **Configuration Management** - Environment variables dengan Viper
9. **Frontend Integration** - React + TypeScript, API client, state management
10. **Containerization** - Dockerfile multi-stage build

## Lisensi

```
MIT License

Copyright (c) 2026 Irfan Yulianto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
