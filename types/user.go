package types

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	Firstname string `bson:"firstname" json:"firstname"`
	Lasttname string `bson:"lastname" json:"lastname"`
}
