package model

type UserRecord struct {
	UUID     string
	Login    string
	Password string
}

func (ur *UserRecord) Clone() *UserRecord {
	return &UserRecord{
		UUID:     ur.UUID,
		Login:    ur.Login,
		Password: ur.Password,
	}
}
