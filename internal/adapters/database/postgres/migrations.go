package postgres

import "webTemplate/internal/domain/entity"

// Migrations is a list of all gorm migrations for the database.
var Migrations = []interface{}{
	&entity.User{},
	&entity.Token{},
}
