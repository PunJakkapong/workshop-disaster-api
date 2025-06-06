# Workshop Disaster API

API สำหรับจัดการข้อมูลภัยพิบัติ

## Prerequisites

- Go 1.22 หรือสูงกว่า
- Docker
- PostgreSQL
- Redis

## การติดตั้ง

1. Clone repository:

```bash
git clone https://github.com/PunJakkapong/workshop-disaster-api.git
cd workshop-disaster-api
```

2. ติดตั้ง dependencies:

```bash
go mod download
```

3. สร้างไฟล์ .env:

```bash
cp .env.example .env
```

แก้ไขค่าใน .env ตาม environment ของคุณ:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=disaster_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
PORT=8080
HOST=0.0.0.0
```

## การรันแอพพลิเคชัน

### รันแบบ Local

1. รันด้วย Go:

```bash
go run main.go
```

2. หรือ Build และรัน:

```bash
go build -o main
./main
```

### รันด้วย Docker

1. Build Docker image:

```bash
docker build -t workshop-disaster-api-golang-api:v1 .
```

2. รัน Docker container:

```bash
docker run -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  -e REDIS_HOST=your-redis-host \
  workshop-disaster-api-golang-api:v1
```

### รันด้วย Docker Compose

```bash
docker-compose up -d
```

## API Endpoints

- `GET /api/v1/disasters` - ดึงข้อมูลภัยพิบัติทั้งหมด
- `GET /api/v1/disasters/{id}` - ดึงข้อมูลภัยพิบัติตาม ID
- `POST /api/v1/disasters` - เพิ่มข้อมูลภัยพิบัติใหม่
- `PUT /api/v1/disasters/{id}` - อัพเดทข้อมูลภัยพิบัติ
- `DELETE /api/v1/disasters/{id}` - ลบข้อมูลภัยพิบัติ
