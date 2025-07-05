# 🌱 Smart Garden IoT System

Selamat datang di repositori **Smart Garden IoT**! 🚀

Solusi lengkap untuk **otomatisasi pemantauan dan penyiraman tanaman/taman** menggunakan **ESP32, sensor, fuzzy logic, API Go, dan Dashboard Python (CustomTkinter)** secara **real-time dan cerdas**.

---

## ✨ Fitur Utama

* ✅ Pemantauan Real-time suhu, kelembaban udara, tanah, dan level air.
* ✅ Fuzzy Logic Cerdas untuk keputusan penyiraman adaptif.
* ✅ Kontrol Pompa Otomatis dengan cooldown protection.
* ✅ Deteksi Level Air agar pompa aman saat air habis.
* ✅ Dashboard Interaktif real-time & kontrol manual/otomatis.
* ✅ Manajemen User & Device oleh admin.
* ✅ API RESTful stabil untuk komunikasi ESP32, Go, dan Dashboard.

---

## 🛠️ Arsitektur Sistem

1. **ESP32 + Sensor & Pompa**: Membaca data & fuzzy logic.
2. **Backend API (Go)**: Menyimpan data, autentikasi, & pengelolaan perangkat.
3. **Dashboard (Python/CustomTkinter)**: Visualisasi data sensor & kontrol pompa.

---

## 🔄 Alur Kerja

1. ESP32 terhubung Wi-Fi & autentikasi ke API backend.
2. Membaca data sensor secara berkala.
3. Data diproses fuzzy logic untuk keputusan penyiraman.
4. Pompa aktif/nonaktif sesuai hasil fuzzy & kondisi air.
5. Data dikirim ke API backend.
6. Dashboard Python menampilkan data real-time.
7. User dapat memantau & kontrol pompa via dashboard.
8. Admin mengelola user & device via dashboard.

---

## 🗄️ Struktur Database

* 📂 `users`: data user (ID, username, email, password, role)
* 📂 `devices`: data perangkat (ID, nama, lokasi, status, mode, IP, user\_id)
* 📂 `sensor_readings`: data sensor (ID, device\_id, suhu, kelembaban, tanah, air, status pompa)

---

## 🏗️ Struktur Proyek (API Backend - Go)

```
.
├── config/
│   └── database.go
├── controllers/
│   ├── auth_controller.go
│   ├── device_controller.go
│   └── sensor_controller.go
├── middleware/
│   └── auth_middleware.go
├── models/
│   ├── device.go
│   ├── sensor_data.go
│   └── user.go
├── utils/
│   └── jwt.go
├── .env
├── go.mod
├── go.sum
├── main.go
```

---

## ⚙️ Instalasi & Penggunaan

### 1️⃣ Kloning Repo

```bash
git clone https://github.com/your-username/smart-garden-iot.git
cd smart-garden-iot
```

### 2️⃣ Jalankan Backend (Go)

* Masuk ke `backend/`
* Buat `.env` dari `.env.example`
* Import `smart_garden.sql` ke database MySQL
* Install dependensi:

```bash
go mod tidy
```

* Jalankan server:

```bash
go run main.go
```

### 3️⃣ Upload Firmware ESP32

* Buka `fuzzy_logic_smart_garden.ino`
* Sesuaikan `ssid`, `password`, dan `api_base_url`
* Upload ke ESP32

### 4️⃣ Jalankan Dashboard (Python)

* Masuk ke `dashboard/`
* Install dependensi:

```bash
pip install -r requirements.txt
```

* Jalankan:

```bash
python main.py
```

✨ **Selamat berkebun cerdas & otomatis!** 🌿
