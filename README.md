# Crowdfunding-go
Crowdfunding API

## Study Case

The tight competition in the selection of research proposals and the decline in the nominal PKM funding in 2021 have hampered students who want to work because of insufficient funding. In addition, there is no alternative funding media at the Yogyakarta Technological University, resulting in students who lack funding or have not passed the funding selection being hampered from realizing their creative ideas. Another obstacle commonly faced by students in seeking support is the difficulty of means of publication. So we need an application that can be an alternative media for funding and alternative publications. This study aims to produce a crowdfunding application as an alternative media for funding and publication of students' creative ideas at the Yogyakarta Technological University that can reach donors more broadly and make it easier for donors to make donations. In the process of making the system, the data collection method is carried out by interviews and observations regarding the efforts that can be made by students to get funding. To overcome this problem, researchers build web and mobile-based crowdfunding applications using the GIN framework, Flutter framework, Vue framework, and MySQL as database. The research results are in the form of a crowdfunding application that is integrated with a payment gateway service to automatically validate donation payments so that it can become an alternative media for funding and publication, as well as making it easier for donors to make donations.

## Technology 
- Gin Gonic 
- GORM 
- MySQL

## Functional Analysis
| No. | Functional  | Description |
| --- | ------------- | ------------- |
| 1 | Login & Register  | The system can save new user data by registering on the form provided and authenticating the user.  |
| 2 | Open Funding  | The system can store open funding requests entered by the user with the specified service standards (identifiers, and proposals). |
| 3 | Funding | The system can facilitate donors for funding projects using a payment gateway for automatic validation. |
| 4 | Transaction History Report | The system can generate reports from each funding project containing the names of donors and the amount collected. |
| 5 | Prize Claims | The system can sort out users who can claim rewards based on the specified minimum donation amount. | 
| 6 | Documentation | Users must document the effects of the funding projects that have been carried out. And will be displayed in the success story menu. |

## Run Project

```
go run main.go
```

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
{
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
}
```
