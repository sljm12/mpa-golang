module github.com/sljm12/mpa-webservice

go 1.22.2

replace mpa => ../mpa-go

require (
	github.com/magiconair/properties v1.8.10
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/paulmach/orb v0.12.0
	mpa v0.0.0-00010101000000-000000000000
)

require go.mongodb.org/mongo-driver v1.11.4 // indirect
