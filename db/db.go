package db

const (
	DBNAME     = "hotel_reservation"
	TestDBNAME = "hotel_reservation_test"
	DBURI      = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
