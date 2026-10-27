package constants

type Role string

const (
	SUPER_ADMIN    Role = "SADM"
	ADMIN          Role = "ADM"
	SITE_MANAGER   Role = "SM"
	AREA_MANAGER   Role = "AM"
	SAFETY_OFFICER Role = "SO"
	EMPLOYEE       Role = "EMP"
)
