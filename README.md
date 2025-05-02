# Go-Api-market

ระบบ api ร้านค้าออนไลน์ โดยใช้ Golang Fiber gorm (DB PostgresSQL) 

## framework 
Fiber 	github.com/gofiber/fiber/v2 v2.52.6 
Gorm    gorm.io/gorm v1.25.12 
JWT 	github.com/golang-jwt/jwt/v5 v5.2.2 
DB      gorm.io/driver/postgres v1.5.11 

## Features
User
- Register 
- Login
- Get User 
- Get Profile
- Update Profile
- Update Password
- Delete User

Producet 
- Creact Product
- Get Product
- Update Product 
- Delete Product 

Order
- Creat Order
- Get Order
- Update Order
- Delete Order

## โครงสร้างโปรเจกต์
├── cmd/main.go
├── config/ # การตั้งค่าฐานข้อมูล
├── models/ # โครงสร้างของข้อมูล
├── Hanler/ # จัดการ request/response
├── middleware/ # จัดการ JWT และสิทธิ์การเข้าถึง
├── repository/ # โค้ดจัดการฐานข้อมูล (CRUD logic)
├── routes/ # กำหนด routing
├── serviceฝ # ธุรกิจลอจิก (Business Logic)
├── middleware/ # JWT middleware
└── utils/ # ฟังก์ชันเสริม เช่น hashing


## การติดตั้งและใช้งาน

### 1. Clone Project

git clone https://github.com/birdlax/go-api-market.git

### 2. ติดตั้ง depen

go mod tidy

### 3. ตั้งค่า config .env

### 4. run sever

cd cmd
go run dev 
