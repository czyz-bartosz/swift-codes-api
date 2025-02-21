package models

type SwiftValidator interface {
	Struct(s interface{}) error
}
