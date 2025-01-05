"# go-clean-arsitektur-" 

# GO Clean Architecture Project : Order

---

## Project Structure 🏛

```
...
config
├── database.go             # Koneksi Database
├── redis.go                # Koneksi Redis
domain                      # Domain Layer -> Lapisan ini berisi entitas inti aplikasi, yaitu representasi data dan logika bisnis fundamental.
├── category.go 
├── .........go
handler                     #Interface Adapters (Handler Layer) -> jembatan antara lapisan logika bisnis dan user interface
├── category_handler.go 
├── ........_handler.go
repository                  #Data Access (Repository Layer) -> bertanggung jawab untuk interaksi langsung dengan database
├── category_repository.go 
├── ........_repository.go
routes
├── category_routes.go 
├── ........_routes.go
service                     #Use Cases (Service Layer): logika bisnis aplikasi yang mendasari, seperti manipulasi data dan aturan bisnis yang lebih kompleks.
├── category_service.go 
├── ........_service.go
tests                       # folder untuk Unit test/ End to End test
├── e2e_test.go 
├── ........_test.go
utils                       #utilitas umum yang dapat digunakan di seluruh proyek
├── response.go 
├── .........go
.env
.env.example
.gitignore
db.sql
go.mod
go.sum
main.go
README.md
```
---
