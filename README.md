# Kasir API

Aplikasi **Point of Sale (POS)** full-stack yang dibangun sebagai hasil pembelajaran **Golang Basic Bootcamp**. Menampilkan backend REST API menggunakan Go (tanpa framework) dan frontend web menggunakan React + TypeScript.

## Tentang Project

Project ini merupakan hasil bootcamp Golang dasar yang berkembang menjadi aplikasi POS yang lebih matang:

- **Sesi 1-2**: Setup project Go, koneksi database PostgreSQL, CRUD dasar, dan clean architecture
- **Sesi 3**: Fitur transaksi checkout dan laporan penjualan harian
- **Sesi 4**: Middleware (CORS, Logger) dan pembuatan frontend React
- **Enhancement**: JWT Authentication, role-based access control, rate limiting, structured logging, graceful shutdown, dan hardened deployment

### Konsep yang Diterapkan

1. **Go Fundamentals** - Struct, interface, error handling, package management
2. **HTTP Server** - Routing, request/response handling tanpa framework (`net/http`)
3. **Clean Architecture** - Separation of concerns dengan layer yang jelas
4. **Dependency Injection** - Constructor-based DI tanpa library
5. **Database** - Koneksi PostgreSQL, raw SQL queries, connection pooling
6. **Middleware Pattern** - Function chaining untuk cross-cutting concerns
7. **Authentication** - JWT token-based auth dengan bcrypt password hashing
8. **Role-Based Access Control** - Admin dan Kasir dengan permission berbeda
9. **REST API Design** - Resource-based URL, HTTP methods, status codes, JSON responses
10. **Configuration Management** - Environment variables dengan Viper
11. **Structured Logging** - `log/slog` dengan log levels dan request tracing
12. **Graceful Shutdown** - Signal handling (SIGINT/SIGTERM)
13. **Rate Limiting** - Token bucket algorithm per IP
14. **Containerization** - Multi-stage Docker build dengan non-root user

## Tech Stack

### Backend

| Teknologi | Versi | Keterangan |
|-----------|-------|------------|
| Go | 1.24 | Bahasa pemrograman utama |
| net/http | stdlib | HTTP server tanpa framework |
| PostgreSQL | - | Database relasional |
| lib/pq | 1.10.9 | PostgreSQL driver |
| Viper | 1.18.2 | Manajemen konfigurasi dari `.env` |
| golang-jwt/jwt | 5.3.1 | JWT token generation & validation |
| golang.org/x/crypto | 0.48.0 | bcrypt password hashing |
| log/slog | stdlib | Structured logging |

### Frontend

| Teknologi | Versi | Keterangan |
|-----------|-------|------------|
| React | 19.2 | UI library |
| TypeScript | 5.9 | Type-safe JavaScript |
| Vite | 7.3 | Build tool dan dev server |
| React Router | 7.13 | Client-side routing |
| Tailwind CSS | 4.1 | Utility-first CSS framework |
| shadcn/ui | - | Komponen UI berbasis Radix UI |
| Lucide React | 0.564 | Icon library |
| Sonner | 2.0 | Toast notifications |
| date-fns | 4.1 | Date utilities |

## Arsitektur

### Backend - Clean Architecture

```
Request → Middleware Stack → Handler → Service → Repository → Database
```

| Layer | Tanggung Jawab |
|-------|---------------|
| **Middleware** | CORS, auth, logging, rate limiting, panic recovery, request ID |
| **Handler** | Menerima HTTP request, validasi input, kirim JSON response |
| **Service** | Business logic dan aturan bisnis |
| **Repository** | Akses data, query SQL, database transactions |
| **Model** | Struktur data / entity |
| **Config** | Konfigurasi aplikasi dari environment |
| **Database** | Koneksi dan connection pool management |

### Middleware Stack

Request diproses melalui middleware stack secara berurutan (outermost → innermost):

| # | Middleware | Fungsi |
|---|-----------|--------|
| 1 | `CORS` | Cross-origin resource sharing, configurable origins |
| 2 | `RequestID` | Generate/propagate `X-Request-ID` header |
| 3 | `Logger` | Structured logging (slog) — method, path, status, duration, IP |
| 4 | `RateLimiter` | Token bucket rate limiting per IP (10 req/s, burst 30) |
| 5 | `Recovery` | Panic recovery dengan stack trace logging |
| 6 | `MaxBody` | Limitasi ukuran request body (1MB) |
| 7 | `JSONErrors` | Konversi `http.Error()` text/plain ke JSON format |

Per-route middleware:
| Middleware | Fungsi |
|-----------|--------|
| `JWTAuth` | Validasi JWT token dari `Authorization: Bearer <token>` |
| `RequireRole` | Enforce role-based access (`admin` / `kasir`) |

### Authentication & Authorization

- **JWT Token** — Token-based authentication, expiry configurable (default 24 jam)
- **bcrypt** — Password di-hash dengan bcrypt sebelum disimpan
- **Role-based** — Dua role: `admin` (full access) dan `kasir` (POS, transaksi, laporan)
- **Register** — Hanya admin yang bisa membuat user baru
- **Admin pertama** — Di-seed langsung ke database

### Frontend Architecture

```
ErrorBoundary → AuthProvider → BrowserRouter → ProtectedRoute → Layout → Pages
```

| Komponen | Fungsi |
|----------|--------|
| `ErrorBoundary` | Menangkap React error, tampilkan fallback UI |
| `AuthContext` | Manajemen state autentikasi (token, user, role) |
| `ProtectedRoute` | Guard route berdasarkan auth status dan role |
| `Layout` | Sidebar navigation, user info, logout, role-based menu |

## Struktur Project

```
kasir-api/
├── backend/
│   ├── config/
│   │   └── config.go              # Konfigurasi (port, DB, JWT, CORS)
│   ├── database/
│   │   └── database.go            # Koneksi DB + connection pool
│   ├── handler/
│   │   ├── auth_handler.go        # Login & Register endpoints
│   │   ├── category_handler.go    # CRUD kategori
│   │   ├── product_handler.go     # CRUD produk
│   │   ├── report_handler.go      # Laporan penjualan
│   │   └── transaction_handler.go # Checkout & riwayat transaksi
│   ├── middleware/
│   │   ├── cors.go                # CORS (configurable origins)
│   │   ├── jsonerrors.go          # Auto-convert errors ke JSON
│   │   ├── jwt.go                 # JWT auth + RequireRole
│   │   ├── logger.go              # Structured logging (slog)
│   │   ├── ratelimit.go           # Token bucket rate limiter
│   │   ├── recovery.go            # Panic recovery + MaxBody
│   │   └── requestid.go           # X-Request-ID generation
│   ├── model/
│   │   ├── category.go
│   │   ├── product.go
│   │   ├── report.go
│   │   ├── transaction.go
│   │   └── user.go                # User, LoginRequest, AuthResponse
│   ├── repository/
│   │   ├── category_repository.go
│   │   ├── product_repository.go
│   │   ├── report_repository.go
│   │   ├── transaction_repository.go  # Batch query (no N+1)
│   │   └── user_repository.go
│   ├── service/
│   │   ├── auth_service.go        # Login, Register, JWT generation
│   │   ├── category_service.go
│   │   ├── product_service.go
│   │   ├── report_service.go
│   │   └── transaction_service.go
│   ├── main.go                    # Entry point, DI, routing, graceful shutdown
│   ├── schema.sql                 # DDL + seed data
│   ├── go.mod
│   └── .env                       # Environment variables (git-ignored)
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ErrorBoundary.tsx   # React Error Boundary
│   │   │   ├── Layout.tsx          # Sidebar, nav, user menu
│   │   │   ├── ProtectedRoute.tsx  # Auth + role guard
│   │   │   └── ui/                 # shadcn/ui components (14 komponen)
│   │   ├── contexts/
│   │   │   └── AuthContext.tsx     # Auth state + useAuth hook
│   │   ├── lib/
│   │   │   ├── api.ts             # API client, JWT inject, timeout, error handling
│   │   │   └── utils.ts           # Helper functions (cn)
│   │   ├── pages/
│   │   │   ├── Login.tsx           # Halaman login
│   │   │   ├── Kasir.tsx           # POS / kasir (product + cart)
│   │   │   ├── Products.tsx        # Manajemen produk (admin)
│   │   │   ├── Categories.tsx      # Manajemen kategori (admin)
│   │   │   ├── Transactions.tsx    # Riwayat transaksi
│   │   │   └── Report.tsx          # Laporan penjualan
│   │   ├── App.tsx                 # Routes + ProtectedRoute wrapper
│   │   ├── main.tsx                # Entry point (ErrorBoundary + AuthProvider)
│   │   └── index.css               # Tailwind CSS + theme tokens
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── Dockerfile                      # Multi-stage build
├── .dockerignore
├── .gitignore
└── README.md
```

## Database Schema

```sql
users
├── id          SERIAL PRIMARY KEY
├── username    VARCHAR(100) NOT NULL UNIQUE
├── password    VARCHAR(255) NOT NULL        -- bcrypt hash
├── role        VARCHAR(20) NOT NULL DEFAULT 'kasir' CHECK (admin|kasir)
└── created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP

categories
├── id          SERIAL PRIMARY KEY
├── name        VARCHAR(255) NOT NULL
└── description TEXT

products
├── id          SERIAL PRIMARY KEY
├── name        VARCHAR(255) NOT NULL
├── price       INTEGER NOT NULL DEFAULT 0 CHECK (>= 0)
├── stock       INTEGER NOT NULL DEFAULT 0 CHECK (>= 0)
└── category_id INTEGER → categories(id) ON DELETE SET NULL

transactions
├── id           SERIAL PRIMARY KEY
├── total_amount INTEGER NOT NULL
└── created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP

transaction_details
├── id             SERIAL PRIMARY KEY
├── transaction_id INTEGER → transactions(id) ON DELETE CASCADE
├── product_id     INTEGER → products(id)
├── quantity       INTEGER NOT NULL CHECK (> 0)
└── subtotal       INTEGER NOT NULL CHECK (>= 0)
```

**Indexes:**

| Index | Tabel | Kolom |
|-------|-------|-------|
| `idx_products_category_id` | products | category_id |
| `idx_transactions_created_at` | transactions | created_at |
| `idx_td_transaction_id` | transaction_details | transaction_id |
| `idx_td_product_id` | transaction_details | product_id |
| `idx_users_username` | users | username |

## API Endpoints

### Public (tanpa autentikasi)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/auth/login` | Login, mendapat JWT token |
| GET | `/health` | Health check + status database |

### Protected (butuh JWT token)

| Method | Endpoint | Role | Deskripsi |
|--------|----------|------|-----------|
| POST | `/api/auth/register` | admin | Registrasi user baru |
| GET | `/api/categories` | semua | Daftar kategori |
| POST | `/api/categories` | admin | Buat kategori |
| GET | `/api/categories/{id}` | semua | Detail kategori |
| PUT | `/api/categories/{id}` | admin | Update kategori |
| DELETE | `/api/categories/{id}` | admin | Hapus kategori |
| GET | `/api/categories/{id}/products` | semua | Produk per kategori |
| GET | `/api/products` | semua | Daftar produk |
| POST | `/api/products` | admin | Buat produk |
| GET | `/api/products/{id}` | semua | Detail produk |
| PUT | `/api/products/{id}` | admin | Update produk |
| DELETE | `/api/products/{id}` | admin | Hapus produk |
| POST | `/api/checkout` | semua | Proses transaksi |
| GET | `/api/transactions` | semua | Riwayat transaksi |
| GET | `/api/report/today` | semua | Laporan penjualan |

### Contoh Request

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin",
    "created_at": "2026-02-16T00:00:00Z"
  }
}
```

**List Products (authenticated):**
```bash
curl http://localhost:8080/api/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Checkout (authenticated):**
```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": 1, "quantity": 2 },
      { "product_id": 2, "quantity": 3 }
    ]
  }'
```

**Register User (admin only):**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"username": "kasir1", "password": "password123", "role": "kasir"}'
```

**Laporan Penjualan (range tanggal):**
```bash
curl "http://localhost:8080/api/report/today?start_date=2026-02-01&end_date=2026-02-15" \
  -H "Authorization: Bearer <token>"
```

**Health Check:**
```bash
curl http://localhost:8080/health
# {"status":"OK","database":"connected"}
```

## Cara Menjalankan

### Prasyarat

- [Go](https://go.dev/dl/) >= 1.24
- [Node.js](https://nodejs.org/) >= 18
- [PostgreSQL](https://www.postgresql.org/) (atau layanan cloud seperti Supabase)

### 1. Clone Repository

```bash
git clone https://github.com/irfan-yulianto/kasir-api.git
cd kasir-api
```

### 2. Setup Database

Jalankan script SQL untuk membuat tabel, indexes, dan seed data:

```bash
psql -U postgres -d your_database -f backend/schema.sql
```

### 3. Seed Admin User

Buat admin pertama (hash password terlebih dahulu):

```bash
# Generate bcrypt hash untuk password admin
# Atau gunakan endpoint bootstrap / script terpisah
```

```sql
INSERT INTO users (username, password, role)
VALUES ('admin', '$2a$10$YOUR_BCRYPT_HASH', 'admin');
```

### 4. Konfigurasi Environment

Buat file `backend/.env`:

```env
PORT=8080
DB_CONN=postgresql://username:password@localhost:5432/your_database
JWT_SECRET=your-secure-random-secret-key
JWT_EXPIRY_HOURS=24
ALLOWED_ORIGINS=*
```

| Variable | Required | Default | Deskripsi |
|----------|----------|---------|-----------|
| `PORT` | - | `8080` | Port HTTP server |
| `DB_CONN` | **Ya** | - | PostgreSQL connection string |
| `JWT_SECRET` | **Ya** | - | Secret key untuk sign JWT (fatal jika kosong) |
| `JWT_EXPIRY_HOURS` | - | `24` | Masa berlaku token dalam jam |
| `ALLOWED_ORIGINS` | - | `*` | CORS allowed origins (comma-separated) |

### 5. Jalankan Backend

```bash
cd backend
go run main.go
```

Server berjalan di `http://localhost:8080`

### 6. Jalankan Frontend (Development)

```bash
cd frontend
npm install
npm run dev
```

Frontend berjalan di `http://localhost:5173`

### 7. (Opsional) Go Workspace

Jika menggunakan VS Code, buat `go.work` di root agar IDE mengenali module:

```bash
go work init ./backend
```

## Fitur Frontend

### Halaman

| Halaman | Route | Role | Deskripsi |
|---------|-------|------|-----------|
| Login | `/login` | Public | Form login username + password |
| Kasir (POS) | `/` | Semua | Daftar produk + keranjang belanja |
| Produk | `/products` | Admin | CRUD manajemen produk |
| Kategori | `/categories` | Admin | CRUD manajemen kategori |
| Transaksi | `/transactions` | Semua | Riwayat transaksi dengan filter tanggal |
| Laporan | `/report` | Semua | Laporan penjualan harian |

### Fitur Utama

- Daftar produk dalam card grid dengan info harga dan stok
- Keranjang belanja dengan atur jumlah item (+/-)
- Cart persistence di localStorage (tidak hilang saat refresh)
- Validasi stok sebelum checkout (revalidasi dari server)
- Proses checkout dengan struk transaksi
- Filter transaksi dan laporan berdasarkan range tanggal (auto-filter)
- Manajemen produk dan kategori via dialog modal
- Format mata uang Rupiah (IDR)
- Toast notification untuk feedback aksi
- Role-based navigation (menu berbeda untuk admin vs kasir)
- User info dan tombol logout di sidebar
- Auto redirect ke login saat token expired / 401
- Error Boundary untuk menangkap crash dan tampilkan fallback

## Docker Deployment

### Build & Run

```bash
# Build image dari root project
docker build -t kasir-api .

# Run container
docker run -p 8080:8080 \
  -e DB_CONN="postgresql://user:pass@host:5432/db" \
  -e JWT_SECRET="your-secret-key" \
  -e ALLOWED_ORIGINS="https://yourdomain.com" \
  kasir-api
```

### Docker Architecture

Multi-stage build menghasilkan image minimal (~30MB):

| Stage | Base Image | Output |
|-------|-----------|--------|
| 1. Frontend Builder | `node:20-alpine` | Static files (`dist/`) |
| 2. Backend Builder | `golang:1.24-alpine` | Binary (`main`) |
| 3. Production | `alpine:3.21` | Binary + static files |

**Security features:**
- Non-root user (`app:1000`)
- Tidak ada build tools/source code di production image
- Health check built-in (`/health` endpoint)
- Pinned Alpine version

```dockerfile
# Health check configuration
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

## Server Configuration

### Connection Pool

```
MaxOpenConns:    25
MaxIdleConns:    5
ConnMaxLifetime: 5 minutes
ConnMaxIdleTime: 2 minutes
```

### HTTP Server Timeouts

```
ReadTimeout:  15 seconds
WriteTimeout: 15 seconds
IdleTimeout:  60 seconds
```

### Rate Limiting

- Algorithm: Token Bucket per IP
- Rate: 10 requests/second
- Burst: 30 requests
- Stale IP cleanup: setiap 1 menit

### Graceful Shutdown

Server menangani `SIGINT` dan `SIGTERM` signals dengan timeout 30 detik, memastikan semua request yang sedang diproses selesai sebelum shutdown.

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
