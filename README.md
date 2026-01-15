# 🧾 Tradesman Management System API

A simple yet powerful backend API that allows neighborhood small tradesmen to manage their products and track customer orders from mobile applications.

## 🚀 Features


- ✅ **JWT Authentication** - Secure login system
- ✅ **Role-Based Authorization** - Admin, Shop, Customer roles
- ✅ **Shop Management** - Create and edit shops
- ✅ **Product Management** - Add, update, delete products
- ✅ **Order System** - Customer orders and status tracking
- ✅ **SQLite Database** - Lightweight and practical
- ✅ **Swagger Documentation** - Interactive API documentation
- ✅ **CORS Support** - Ready for frontend integration

## ⚙️ Technologies

- **Go 1.21+** - Modern and performant backend
- **Gin** - Fast web framework
- **GORM** - Powerful ORM library
- **SQLite** - Lightweight database
- **JWT** - Token-based authentication
- **Swagger** - API documentation

## 🏃‍♂️ Quick Start

### 1. Project Setup

```bash
# Install dependencies
go mod tidy

# Generate Swagger documentation (optional)
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 2. Start Server

```bash
go run main.go
```

When server starts successfully:
- 🌐 **API**: http://localhost:8080
- 📚 **Swagger**: http://localhost:8080/swagger/index.html

## 📡 API Endpoints

### 🔐 Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /auth/me` - Profile information (🔒 Auth required)

### 🏪 Shop Management
- `GET /shops` - List all shops
- `GET /shops/{id}` - Shop details
- `POST /shops` - Create new shop (🔒 Shop role)
- `PUT /shops/{id}` - Update shop information (🔒 Shop role)
- `GET /shops/{id}/products` - Shop's products

### 📦 Product Management
- `GET /products` - List all products
- `GET /products/{id}` - Product details
- `POST /products` - Add new product (🔒 Shop role)
- `PUT /products/{id}` - Update product (🔒 Shop role)
- `DELETE /products/{id}` - Delete product (🔒 Shop role)

### 🛒 Order Management
- `POST /orders` - Place order (🔒 Customer role)
- `GET /orders` - List orders (🔒 Auth required)
- `GET /orders/{id}` - Order details (🔒 Auth required)
- `PUT /orders/{id}/status` - Update order status (🔒 Shop role)

## 👥 User Roles

### 🛒 **Customer**
- Can view shops and products
- Can place orders
- Can track their own orders

### 🏪 **Shop (Tradesman)**
- Can create and manage shop
- Can add, update, delete products
- Can view incoming orders
- Can update order statuses

### 👑 **Admin**
- Access to all data
- System-wide control

## 🔒 Authentication

The API uses JWT token-based authentication. You can enter your token by clicking the "Authorize" button in the Swagger interface.

**Header Format:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 📊 Database Schema

### Users
- `id`, `email`, `password`, `name`, `phone`, `role`, `created_at`, `updated_at`

### Shops
- `id`, `user_id`, `name`, `description`, `address`, `phone`, `is_active`, `created_at`, `updated_at`

### Products
- `id`, `shop_id`, `name`, `description`, `price`, `stock`, `is_active`, `image_url`, `created_at`, `updated_at`

### Orders
- `id`, `user_id`, `shop_id`, `total_amount`, `status`, `note`, `created_at`, `updated_at`

### Order Items
- `id`, `order_id`, `product_id`, `quantity`, `price`, `created_at`

## 📋 Order Statuses

- `pending` - Pending
- `confirmed` - Confirmed
- `preparing` - Preparing
- `ready` - Ready
- `delivered` - Delivered
- `cancelled` - Cancelled

## 🛠️ Development

### Creating Test Data

1. First register a shop user:
```json
POST /auth/register
{
  "name": "John Smith",
  "email": "john@example.com",
  "password": "123456",
  "phone": "0555-123-4567",
  "role": "shop"
}
```

2. Create a shop:
```json
POST /shops
{
  "name": "John's Grocery",
  "description": "The best grocery in the neighborhood",
  "address": "Main Street No:15",
  "phone": "0555-123-4567"
}
```

3. Add a product:
```json
POST /products
{
  "name": "Bread",
  "description": "Fresh daily bread",
  "price": 2.50,
  "stock": 100,
  "image_url": "https://example.com/bread.jpg"
}
```

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache 2.0 License. See the `LICENSE` file for details.

## 📞 Contact

If you have any questions, please open an issue or send an email.

---

**Happy coding! 🚀** 