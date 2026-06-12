package mapper

import (
	dmodel "tictactoe/internal/domain/model"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

func ToUserStorage(user *dmodel.User) *dsmodel.UserRecord {
	return &dsmodel.UserRecord{
		UUID:     user.UUID,
		Login:    user.Login,
		Password: user.Password,
	}
}

func ToUserDomain(record *dsmodel.UserRecord) *dmodel.User {
	return &dmodel.User{
		UUID:     record.UUID,
		Login:    record.Login,
		Password: record.Password,
	}
}
