package common

const (
	ProcessDone     = 1000
	ErrQueryDb      = 5809
	ErrAddEntity    = 5810
	ErrUpdateEntity = 5811
	ErrDelEntity    = 5812

	LoginSuccess          = 2001
	LoginErr              = 2009
	RegisterSuccess       = 2002
	RegisterUnSuccess     = 2003
	ErrCreAuth            = 2004
	VerifyOtpErr          = 4001
	DelOtpErr             = 4002
	EntityNotExists       = 4004
	OtpNotExists          = 4005
	EntityIsExists        = 4009
	ErrAuthTypeIsNotValid = 4010
	UserNotExists         = 4011
	UpdateUserErr         = 4012
)

func GetCodeText(code int) string {
	switch code {
	case ErrQueryDb:
		return "Query database error"
	case ErrAddEntity:
		return "Add entity error"
	case ErrUpdateEntity:
		return "Update entity error"
	case ErrDelEntity:
		return "Delete entity error"
	case EntityIsExists:
		return "Entity is has exits"

	case LoginSuccess:
		return "User login successfully"
	case RegisterSuccess:
		return "User register successfully"
	case RegisterUnSuccess:
		return "User registration failed"
	case VerifyOtpErr:
		return "Verify otp do not success"
	case EntityNotExists:
		return "Entity dose not exits"
	case DelOtpErr:
		return "Delete otp record failed"
	case OtpNotExists:
		return "Otp dose not exits"
	case ErrAuthTypeIsNotValid:
		return "auth type is not valid"
	case UserNotExists:
		return "User does not exist"
	case ErrCreAuth:
		return "Create auth failed"
	default:
		return ""
	}
}
