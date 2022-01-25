# Crowdfunding-go (WIP)
Crowdfunding API

# Register User
Endpoint : 
**/api/v1/users**
<br>
Request : 
```JSON
    {
      "name" : "Muhammad Ilham Kusumawardhana", 
      "occupation" : "Mahasiswa",
      "email" : "Kawekaweha00@gmail.com", 
      "password" : "Password123"
    }
```
Response :
```JSON
     "meta": {
        "message": "Berhasil daftar akun",
        "code": 200,
        "status": "sukses"
    },
    "data": {
        "id": 15,
        "name": "Muhammad Ilham Kusumawardhana",
        "occupation": "Mahasiswa",
        "email": "Kawekaweha00@gmail.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNX0.kFZqgyl1J5dln_PR90B1c-9JL-eTv3HQnqHz3O1hiZ8"
    }

```
