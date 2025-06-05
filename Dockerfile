# ------------------- Stage 1: Build -------------------
FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🔧 สร้างโฟลเดอร์ไว้ก่อน (แต่จริง ๆ ควรสร้าง runtime ด้วย)
RUN mkdir -p ./uploads
RUN mkdir -p ./images

# ✅ Copy placeholder images (แต่คำสั่งนี้อาจไม่จำเป็นเพราะมัน copy จาก ./images ไป ./images เอง)
# อันนี้เลย **ลบออก** ได้: 
# RUN cp ./images/placeholder-product.png ./images/

RUN go build -o main

# ------------------- Stage 2: Runtime -------------------
FROM alpine:latest

WORKDIR /root/

# ✅ คัดลอก binary
COPY --from=builder /app/main .

# ✅ คัดลอกโฟลเดอร์ uploads และ placeholder ไปด้วย
COPY --from=builder /app/images ./images
COPY --from=builder /app/uploads   ./uploads  

EXPOSE 3000

CMD ["./main"]
