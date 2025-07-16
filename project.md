# 🧾 Proje: Esnaf Yönetim Sistemi API (Go + Swagger)

## 🎯 Amaç
Mahalledeki küçük esnafın ürünlerini yönetebileceği ve mobil uygulama üzerinden gelen müşteri siparişlerini takip edebileceği sade ama güçlü bir backend API geliştirilmesi.

---

## ⚙️ Teknoloji Seçimi

- ✅ Go (Golang)
- ✅ Gin framework
- ✅ GORM (ORM için)
- ✅ PostgreSQL veya SQLite
- ✅ Swagger (Swaggo ile)
- ✅ JWT ile kimlik doğrulama
- ✅ Role-Based Access (Admin / Esnaf / Müşteri)

---

## 👥 Kullanıcı Roller

1. **Esnaf**: Ürün ekler, günceller, siparişlerini takip eder  
2. **Müşteri**: Mobil uygulamadan ürünleri listeler, sipariş verir  
3. **Admin (opsiyonel)**: Genel denetim yapar

---

## 🔐 Auth

### Endpointler:
- `POST /auth/register`  
- `POST /auth/login`  
- `GET /auth/me`  

JWT Token döner ve Swagger’da Authorize kısmından test edilebilir olmalı.

---

## 🏪 Esnaf Yönetimi

### Endpointler:
- `GET /shops` → Müşteri tüm esnafları görür  
- `GET /shops/{id}`  
- `POST /shops` → Yeni esnaf kaydı (register ile)  
- `PUT /shops/{id}` → Profil güncelleme  

---

## 📦 Ürün Yönetimi

### Endpointler:
- `GET /products` → Tüm ürünleri göster (müşteri)  
- `GET /shops/{id}/products`  
- `POST /products` → Sadece esnaf ekleyebilir  
- `PUT /products/{id}`  
- `DELETE /products/{id}`  

---

## 🛒 Sipariş Yönetimi

### Endpointler:
- `POST /orders`  
```json
{
  "shopId": 1,
  "items": [
    { "productId": 2, "quantity": 1 },
    { "productId": 5, "quantity": 3 }
  ],
  "note": "Not bırakılabilir"
}
