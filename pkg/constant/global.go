package constant

type JwtKey int

const (
	UserIDKey JwtKey = iota
	RoleKey
	OutletIDKey

	FILE_UPLOAD_MAX_SIZE = 1024 * 1024 * 2    // 2MB
	FILE_UPLOAD_MAX_AGE  = 365 * 24 * 60 * 60 // 1 year
)
