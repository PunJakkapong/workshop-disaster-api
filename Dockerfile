# ใช้ Official Golang Image
FROM golang:latest

# ตั้งค่า Work Directory
WORKDIR /app

# คัดลอกไฟล์ทั้งหมดไปยัง Docker
COPY . .

# ดึง Dependencies และ Build
RUN go mod tidy
RUN go build -o app .

# รันโปรแกรม
CMD ["/app/app"]
