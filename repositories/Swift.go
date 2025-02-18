package repositories

import (
	"awesomeProject/models"
	"context"
)

type SwiftRepo interface {
	GetBySwiftCode(context.Context, string) (*models.Swift, error)
	GetBranchesBySwiftCode(context.Context, string) ([]models.SwiftMini, error)
	GetByCountryIso2Code(context.Context, string) ([]models.SwiftMini, error)
	GetCountryNameByIso2Code(context.Context, string) (*string, error)
	AddSwift(context.Context, *models.Swift) error
	DeleteSwift(context.Context, string) error
}
