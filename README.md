# The Note taking API
A simple CRUD app using Go; GORM as ORM; JWT authorization 

## To run enter command:
./mynoteapp.exe

### main.go
1. Database setup 
2. Model migrations
3. Handler functions

### note.go
1. The `Note` model
2. `Authorize()` to validate cookies
3. Controllers: `GetNotes()` `CreateNotes()` `DeleteNotes()` `UpdateNotes()`

### user.go
1. The `User` model
2. `Register()` adds new User to the database
3. `Login()` authenticates User; if successful creates cookies, else throws `401`.
4. `GetAllUsers()` fetches all users int the database

