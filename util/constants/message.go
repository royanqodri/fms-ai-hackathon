package constants

const (
	// general
	ERROR_UNMARSHAL_JSON = "unmarshal attributes error: %v"

	// login
	USER_NOT_FOUND    = "user not found"
	DEVICE_NOT_FOUND  = "device not found"
	PASSWORD_IS_WRONG = "password is wrong"

	// user
	MSG_MAX_ACTIVE_REACHED = "cannot add user because maximum user active has reached limit : 500 users. Current total active user: %d"

	// firebase
	ERROR_CREATE_FIREBASE_APP = "failed to create firebase app"
	ERROR_GET_MESSAGE_CLIENT  = "error getting messaging client"
	ERROR_SEND_MESSAGE        = "error sending message"

	// hazard
	AREA_MANAGER_NOT_FOUND = "area manager in location %s not found"
	AREA_MANAGER_NOT_EXIST = "no area manager exist in this location"

	// user notification
	USER_NOTIFICATION_NOT_FOUND = "user notification is not found"
	NOTIF_TOKEN_NOT_FOUND       = "token is empty or not found"

	// push notif
	ERROR_SEND_PUSH_NOTIF = "error when send push notification: %v"

	// telegram
	BOT_TOKEN_EMPTY_TELEGRAM        = "bot token is empty"
	ERROR_GET_MASTER_NOTIF_TELEGRAM = "error when get master notif telegram: %v"
	ERROR_INIT_TELEGRAM             = "error when initialize telegram: %v"
	ERROR_SEND_TELEGRAM             = "error when send message to telegram: %v"
)
