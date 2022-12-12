### RBTI Golang Server

This server is used for my thesis project, it uses PostgreSQL and for its main database and uses Elasticsearch for faster query speed and Elastic specific features such as Autocomplete.

I uses Gin for it's router because it's more convenient to use rather than using `gorilla/mux` 
also i never uses Gin for production before so i can try out Gin

To run you have to call
`go run app.go`
