A ToDo app built with Go.
- Framework: Go Fiber
- Database: Gorm & Postgres
- Cache: Redis

To get this app running,  
1. Run Postres and Redis on Docker. Create a password for the postgres db
```
docker run -it --name todo-postgres -e POSTGRES_USER=todo -e POSTGRES_PASSWORD=your_password -e POSTGRES_DB=todo -p 5432:5432 -d postgres
docker run -it --name todo-redis -p 6379:6379 -d redis
```
2. Set environment variable TODOPASSWORD to your_password
3. Clone this repository and navigate inside
4. Build app
```
go build
```
5. Run air to start the server
```
air
```
