# Sistem Tiket: Panduan Lengkap

Panduan ini memberikan penjelasan komprehensif tentang API Sistem Tiket, mencakup arsitektur, komponen, dan bagaimana mereka bekerja bersama. Di akhir, Anda akan memahami bagaimana aplikasi Go modern ini distruktur menggunakan prinsip-prinsip arsitektur bersih.

## 1. Gambaran Umum

Sistem Tiket adalah API RESTful yang dibangun dengan Go yang memungkinkan pengguna untuk:

- Mendaftar dan masuk (dengan akses berbasis peran)
- Menjelajahi dan membeli tiket untuk acara
- Mengelola acara (pengguna admin)
- Menghasilkan laporan tentang penjualan tiket dan pendapatan

Sistem ini mengikuti prinsip arsitektur bersih, memisahkan fungsi ke dalam lapisan yang berbeda:

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Controllers   │ ──→  │    Services     │ ──→  │  Repositories   │ ──→  │    Database     │
└─────────────────┘      └─────────────────┘      └─────────────────┘      └─────────────────┘
       ↑                         ↑                        ↑                        ↑
       │                         │                        │                        │
       └───── HTTP Requests      └───── Business Logic    └───── Data Access       └───── Persistence
```

## 2. Teknologi yang Digunakan

- **Go**: Bahasa pemrograman utama
- **Gin**: Framework web untuk menangani permintaan HTTP
- **GORM**: Object-Relational Mapper untuk interaksi database
- **MySQL**: Database untuk penyimpanan data
- **JWT**: Untuk autentikasi dan otorisasi
- **Air**: Untuk hot-reloading selama pengembangan
- **Swagger**: Untuk dokumentasi API

## 3. Struktur Proyek

Proyek ini mengikuti struktur yang bersih dan modular:

```
/
├── main.go               # Titik masuk aplikasi
├── router/               # Definisi rute API
├── controller/           # Penanganan permintaan HTTP
├── service/              # Logika bisnis
├── repository/           # Lapisan akses data
├── entity/               # Model domain
├── middleware/           # Middleware HTTP
├── config/               # Konfigurasi aplikasi
├── docs/                 # Dokumentasi Swagger
└── tests/                # Pengujian unit dan integrasi
```

## 4. Komponen Inti

### 4.1. Entitas (Model)

Entitas mewakili objek domain dan tabel database:

```go
// Contoh: Entitas User
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"password,omitempty"`
    Role      Role      `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Tickets   []Ticket  `json:"tickets,omitempty"`
}
```

Entitas utama:

- **User**: Pengguna akhir sistem (pengguna biasa dan admin)
- **Event**: Acara yang tiketnya bisa dibeli
- **Ticket**: Mewakili tiket yang dibeli untuk acara
- **Audit**: Melacak aktivitas pengguna untuk tujuan audit

### 4.2. Repositories

Repositories menyediakan lapisan abstraksi di atas database:

```go
// Contoh: Interface UserRepository
type UserRepository interface {
    FindAll() ([]entity.User, error)
    FindByID(id uint) (*entity.User, error)
    FindByEmail(email string) (*entity.User, error)
    Save(user *entity.User) error
    Delete(id uint) error
}
```

Setiap repository:

- Mendefinisikan interface untuk operasi akses data
- Memiliki implementasi konkret menggunakan GORM
- Menangani query dan transaksi database
- Mengembalikan entitas domain ke lapisan service

### 4.3. Services

Services berisi logika bisnis:

```go
// Contoh: Interface UserService
type UserService interface {
    Register(user *entity.User) error
    Login(email, password string) (string, error)
    GetUser(id uint) (*entity.User, error)
}
```

Services bertanggung jawab untuk:

- Mengimplementasikan aturan bisnis dan validasi
- Mengkoordinasikan operasi di beberapa repositories
- Menangani transaksi bila diperlukan
- Tidak berurusan langsung dengan permintaan/respons HTTP

### 4.4. Controllers

Controllers menangani permintaan HTTP dan membentuk respons:

```go
// Contoh: Controller registrasi user
func (ctrl *userController) Register(c *gin.Context) {
    var user entity.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := ctrl.userService.Register(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}
```

Controllers:

- Mengurai permintaan HTTP dan memvalidasi input
- Memanggil metode layanan yang sesuai
- Menyusun respons HTTP
- Menangani kode status kesalahan

### 4.5. Router

Router mendefinisikan semua endpoint API dan menghubungkannya ke metode controller:

```go
// Contoh: Menyiapkan rute
func SetupRouter(userController controller.UserController, ...) *gin.Engine {
    router := gin.Default()

    // Rute publik
    router.POST("/register", userController.Register)
    router.POST("/login", userController.Login)

    // Rute terproteksi
    authRoutes := router.Group("/")
    authRoutes.Use(middleware.AuthMiddleware())
    {
        authRoutes.GET("/profile", userController.Profile)
        // ...lebih banyak rute
    }

    return router
}
```

## 5. Alur Autentikasi

1. **Pendaftaran**: Pengguna mendaftar dengan memberikan nama, email, dan kata sandi
2. **Login**: Saat masuk, sistem memverifikasi kredensial dan mengeluarkan token JWT
3. **Otorisasi**: Token disertakan dalam header `Authorization` untuk rute terproteksi
4. **Middleware**: `AuthMiddleware` memverifikasi token dan menetapkan konteks pengguna
5. **RBAC**: Kontrol akses berbasis peran membatasi operasi tertentu hanya untuk pengguna admin

## 6. Implementasi Fitur Utama

### 6.1. Proses Pembelian Tiket

```
1. Pengguna memilih acara (GET /events/:id)
2. Pengguna mengirimkan permintaan pembelian tiket (POST /tickets)
3. Sistem memvalidasi:
   - Acara ada dan aktif
   - Acara belum dimulai
   - Acara masih memiliki kapasitas tersedia
4. Sistem membuat catatan tiket yang ditautkan ke pengguna dan acara
5. Sistem mengembalikan konfirmasi tiket
```

### 6.2. Paginasi dan Filtering

Endpoint API yang mengembalikan daftar (seperti acara dan tiket) mendukung:

- Paginasi berbasis halaman (parameter `page` dan `limit`)
- Filtering berdasarkan berbagai kriteria
- Metadata dalam respons (total item, halaman, dll.)

Contoh respons:

```json
{
  "data": [...items...],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

### 6.3. Jejak Audit

Sistem mencatat tindakan pengguna melalui middleware audit:

- Mencatat siapa melakukan apa dan kapan
- Melacak perubahan pada sumber daya
- Memungkinkan admin untuk meninjau aktivitas sistem

## 7. Skema Database

Database memiliki tabel utama berikut:

- **users**: Menyimpan informasi pengguna dan detail autentikasi
- **events**: Berisi informasi tentang acara yang tersedia
- **tickets**: Mencatat pembelian tiket, menghubungkan pengguna ke acara
- **audit_logs**: Melacak aktivitas pengguna di seluruh sistem

Hubungan utama:

- Satu pengguna dapat memiliki banyak tiket (1:N)
- Satu acara dapat memiliki banyak tiket (1:N)
- Satu pengguna dapat memiliki banyak log audit (1:N)

## 8. Alur Pengembangan

### 8.1. Menyiapkan Proyek

1. Clone repository
2. Buat file `.env` dengan pengaturan database dan JWT
3. Jalankan `go mod download` untuk menginstal dependencies
4. Buat database MySQL
5. Jalankan aplikasi: `go run main.go` atau gunakan Air: `make dev`

### 8.2. Pengembangan dengan Hot Reload

Proyek ini menggunakan Air untuk hot reloading selama pengembangan:

```bash
make dev  # Menjalankan aplikasi dengan Air
```

Ini mengawasi perubahan file dan secara otomatis membangun kembali aplikasi.

### 8.3. Pengujian API

Endpoint API dapat diuji menggunakan:

1. Swagger UI: `http://localhost:8080/swagger/index.html`
2. Postman: Impor koleksi yang disediakan

## 9. Praktik Terbaik yang Digunakan

- **Desain Berbasis Interface**: Layanan dan repositories menggunakan interface untuk dependency injection
- **Pemisahan Fungsi**: Batas yang jelas antar lapisan
- **Tanggung Jawab Tunggal**: Setiap komponen memiliki peran yang jelas
- **Penanganan Kesalahan**: Pendekatan konsisten untuk penanganan dan propagasi kesalahan
- **Validasi**: Validasi input di tingkat controller dan service
- **Dependency Injection**: Komponen menerima dependensi mereka melalui konstruktor

## 10. Kesimpulan

Sistem tiket ini menunjukkan cara membangun API yang kuat menggunakan Go dan prinsip arsitektur bersih. Pemisahan fungsi membuatnya:

- **Dapat Dipelihara**: Perubahan di satu lapisan tidak mempengaruhi yang lain
- **Dapat Diuji**: Komponen dapat diuji secara terpisah
- **Dapat Diperluas**: Fitur baru dapat ditambahkan tanpa refactoring besar

Dengan mengikuti pola-pola ini, Anda dapat membangun aplikasi yang skalabel dan dapat dipelihara yang mudah dipahami dan dikembangkan dari waktu ke waktu.
