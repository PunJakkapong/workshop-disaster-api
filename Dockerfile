# ใช้ Official Golang Image
FROM golang:1.22.0 AS builder

# ตั้งค่า Work Directory
WORKDIR /app

# คัดลอกไฟล์ทั้งหมดไปยัง Container
COPY . .

# ดึง Dependencies และ Build ให้เป็น Linux Binary
RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o /app/app .

# รันโปรแกรม
CMD ["/app/app"]