# 🧾 Esnaf Yönetim Sistemi API

Mahalledeki küçük esnafın ürünlerini yönetebileceği ve mobil uygulama üzerinden gelen müşteri siparişlerini takip edebileceği sade ama güçlü bir backend API.

## 🚀 Özellikler

- ✅ **JWT Kimlik Doğrulama** - Güvenli giriş sistemi
- ✅ **Rol Bazlı Yetkilendirme** - Admin, Esnaf, Müşteri rolleri
- ✅ **Esnaf Yönetimi** - Dükkan oluşturma ve düzenleme
- ✅ **Ürün Yönetimi** - Ürün ekleme, güncelleme, silme
- ✅ **Sipariş Sistemi** - Müşteri siparişleri ve durum takibi
- ✅ **SQLite Veritabanı** - Hafif ve pratik
- ✅ **Swagger Dokümantasyonu** - Interaktif API dokümantasyonu
- ✅ **CORS Desteği** - Frontend entegrasyonu için hazır

## ⚙️ Teknolojiler

- **Go 1.21+** - Modern ve performanslı backend
- **Gin** - Hızlı web framework
- **GORM** - Güçlü ORM kütüphanesi
- **SQLite** - Hafif veritabanı
- **JWT** - Token tabanlı kimlik doğrulama
- **Swagger** - API dokümantasyonu

## 🏃‍♂️ Hızlı Başlangıç

### 1. Proje Kurulumu

```bash
# Bağımlılıkları yükle
go mod tidy

# Swagger dokümantasyonu oluştur (opsiyonel)
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 2. Sunucuyu Başlat

```bash
go run main.go
```

Sunucu başarıyla başladığında:
- 🌐 **API**: http://localhost:8080
- 📚 **Swagger**: http://localhost:8080/swagger/index.html

## 📡 API Endpoints

### 🔐 Authentication
- `POST /auth/register` - Kullanıcı kaydı
- `POST /auth/login` - Kullanıcı girişi
- `GET /auth/me` - Profil bilgileri (🔒 Auth gerekli)

### 🏪 Esnaf Yönetimi
- `GET /shops` - Tüm esnafları listele
- `GET /shops/{id}` - Esnaf detayı
- `POST /shops` - Yeni esnaf oluştur (🔒 Esnaf rolü)
- `PUT /shops/{id}` - Esnaf bilgilerini güncelle (🔒 Esnaf rolü)
- `GET /shops/{id}/products` - Esnafın ürünleri

### 📦 Ürün Yönetimi
- `GET /products` - Tüm ürünleri listele
- `GET /products/{id}` - Ürün detayı
- `POST /products` - Yeni ürün ekle (🔒 Esnaf rolü)
- `PUT /products/{id}` - Ürün güncelle (🔒 Esnaf rolü)
- `DELETE /products/{id}` - Ürün sil (🔒 Esnaf rolü)

### 🛒 Sipariş Yönetimi
- `POST /orders` - Sipariş ver (🔒 Müşteri rolü)
- `GET /orders` - Siparişleri listele (🔒 Auth gerekli)
- `GET /orders/{id}` - Sipariş detayı (🔒 Auth gerekli)
- `PUT /orders/{id}/status` - Sipariş durumu güncelle (🔒 Esnaf rolü)

## 👥 Kullanıcı Rolleri

### 🛒 **Customer (Müşteri)**
- Esnafları ve ürünleri görüntüleyebilir
- Sipariş verebilir
- Kendi siparişlerini takip edebilir

### 🏪 **Shop (Esnaf)**
- Dükkan oluşturabilir ve yönetebilir
- Ürün ekleyebilir, güncelleyebilir, silebilir
- Gelen siparişleri görüntüleyebilir
- Sipariş durumlarını güncelleyebilir

### 👑 **Admin**
- Tüm verilere erişim
- Sistem geneli kontrol

## 🔒 Authentication

API, JWT token tabanlı kimlik doğrulama kullanır. Swagger arayüzünde "Authorize" butonuna tıklayarak token'ınızı girebilirsiniz.

**Header Format:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 📊 Veritabanı Şeması

### Users (Kullanıcılar)
- `id`, `email`, `password`, `name`, `phone`, `role`, `created_at`, `updated_at`

### Shops (Esnaflar)
- `id`, `user_id`, `name`, `description`, `address`, `phone`, `is_active`, `created_at`, `updated_at`

### Products (Ürünler)
- `id`, `shop_id`, `name`, `description`, `price`, `stock`, `is_active`, `image_url`, `created_at`, `updated_at`

### Orders (Siparişler)
- `id`, `user_id`, `shop_id`, `total_amount`, `status`, `note`, `created_at`, `updated_at`

### Order Items (Sipariş Kalemleri)
- `id`, `order_id`, `product_id`, `quantity`, `price`, `created_at`

## 📋 Sipariş Durumları

- `pending` - Beklemede
- `confirmed` - Onaylandı
- `preparing` - Hazırlanıyor
- `ready` - Hazır
- `delivered` - Teslim edildi
- `cancelled` - İptal edildi

## 🛠️ Geliştirme

### Test Verisi Oluşturma

1. Önce bir esnaf kullanıcısı kaydet:
```json
POST /auth/register
{
  "name": "Ahmet Usta",
  "email": "ahmet@example.com",
  "password": "123456",
  "phone": "0555-123-4567",
  "role": "shop"
}
```

2. Dükkan oluştur:
```json
POST /shops
{
  "name": "Ahmet'in Bakkalı",
  "description": "Mahallenin en iyi bakkalı",
  "address": "Atatürk Mah. Cumhuriyet Cad. No:15",
  "phone": "0555-123-4567"
}
```

3. Ürün ekle:
```json
POST /products
{
  "name": "Ekmek",
  "description": "Taze günlük ekmek",
  "price": 2.50,
  "stock": 100,
  "image_url": "https://example.com/ekmek.jpg"
}
```

## 🤝 Katkıda Bulunma

1. Projeyi fork edin
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

## 📄 Lisans

Bu proje Apache 2.0 lisansı altında lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.

## 📞 İletişim

Herhangi bir sorunuz varsa lütfen issue açın veya email gönderin.

---

**Mutlu kodlamalar! 🚀** 