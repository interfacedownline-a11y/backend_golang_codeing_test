คู่มือเบี้องต้น

--------------------------------------

สำหรับ setup docker 
รูปแบบที่1 docker (ที่เดียวจบ)
run => docker compose up --build

--------------------------------------
รูปแบบที่2 แบบธรรมดา (ไม่มี Makefile)
ติดตั้ง dependency
run => go mod tidy

Generate gRPC
protoc \
  --go_out=backend_golang_codeing_test/proto/userpb \
  --go-grpc_out=backend_golang_codeing_test/proto/userpb \
  --proto_path=proto \
  proto/*.proto

รันโปรแกรม
run => go run cmd/main.go

(Optional) รัน Test
go test ./... -v

--------------------------------------
รูปแบบที่3 Makefile กรณี มี make แล้ว
run => make all 

--------------------------------------
postman
อยู่ใน folder \doc\collection-postman
- ให้ import ลงใน postman จะมีตัว postman_collection กับ postman_environment
- แล้วสามารถใช้ได้เลย กรณี import อย่าลืมเปลี่ยน env ที่ใช้ด้วย ทำเป็น environmentไว้ ชื่อ BackendGolangCodingTest นี้ไว้

--------------------------------------
Api

- <POST> /api/auth/register

Ex. requests
{
    "name": "test01",
    "email": "test01@example.com",
    "password": "11111111111"
}

Ex. responses
{
    "status": "success",
    "code": 201,
    "message": "Register successful"
}

***************************************

- <POST> /api/auth/login

Ex. requests
{
    "email": "demo@example.com",
    "password": "123456789"
}

Ex. responses
{
    "status": "success",
    "code": 200,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJnb2xhbmctY29kaW5nLXRlc3QtY2xpZW50IiwiZXhwIjoxNzU0MDY2MDk0LCJpYXQiOjE3NTM5Nzk2OTQsImlzcyI6ImdvbGFuZy1jb2RpbmctdGVzdC1hcGkiLCJqdGkiOiI2NjZiMDEzNy01ZTgyLTRkYTgtYmY3NS04NDVhZjRmZjg2ZGIiLCJuYmYiOjE3NTM5Nzk2OTQsInN1YiI6IjY4OGEwNGNkZGQ4MzhhZDI1NzJjMzVhNiIsInVzZXJfaWQiOiI2ODhhMDRjZGRkODM4YWQyNTcyYzM1YTYiLCJ1c2VybmFtZSI6IksifQ.3y-quw3Jr3ccCj-zX5yQQCLEieYR3Ngo47qMKS7fV7s"
    }
}
หมายเหตุ ตอน login เขียน scripts set token ลง setEnvironmentVariable แล้ว ไม่ต้อง copy มาใส่ในแต่ละที่ใช้ Auth

***************************************

- <GET> /api/users

Ex. responses
{
    "status": "success",
    "code": 200,
    "message": "Fetched users successfully",
    "data": [
        {
            "id": "688a04cddd838ad2572c35a6",
            "name": "K",
            "email": "demo@example.com",
            "createdAt": "0001-01-01T00:00:00Z"
        },
        {
            "id": "688b274e46fab31087a5587a",
            "name": "test01",
            "email": "test02@example.com",
            "createdAt": "0001-01-01T00:00:00Z"
        }
    ]
}

***************************************

- <GET> /api/users/:id  
EX. (/api/users/688a04cddd838ad2572c35a6)

Ex. responses
{
    "status": "success",
    "code": 200,
    "message": "Fetched user successfully",
    "data": {
        "id": "688a04cddd838ad2572c35a6",
        "name": "K",
        "email": "demo@example.com",
        "createdAt": "0001-01-01T00:00:00Z"
    }
}

***************************************

- <PUT> /api/users

Ex. requests
{
    "id": "688a04cddd838ad2572c35a6",
    "email": "demo@example.com",
    "name": "test"
}

Ex. responses
{
    "status": "success",
    "code": 200,
    "message": "User updated successfully"
}

***************************************

- <DELETE> /api/users/:id  

Ex. responses
{
    "status": "success",
    "code": 200,
    "message": "User deleted successfully"
}


/// Any assumptions or decisions made

***************************************

สรุป
- ใช้ Clean Architecture คิดว่า maintenance ง่ายกว่า ที่มีความคล้ายกับที่โจทย์ให้มา(Hexagonal Architecture)
- เน้น ความยืดหยุ่น ทดสอบง่าย