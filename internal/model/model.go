package model

const (
	Success                 int32 = 100
	EmptyUserNameOrPwdError int32 = -1
	UserNameDupError        int32 = -2
	UserNameOrPasswordError int32 = -3
	UserNotExistError       int32 = -4
	UserLoginError          int32 = -5
	UserUpdateError         int32 = -6
	TokenUnmatchedError     int32 = -7
	NilDataError            int32 = -8
	GetUserProfileError     int32 = -9
)

type User struct {
	Id       int64
	UserName string
	NickName string
	Password string
	Avatar   string
}

func (User) TableName() string {
	return "user"
}
