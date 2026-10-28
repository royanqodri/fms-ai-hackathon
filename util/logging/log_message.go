package logging

const (
	// general
	HTTP_REQUEST         = "http request"
	ERROR_UNMARSHAL_JSON = "unmarshal attributes error"

	// user notification
	ERROR_GET_USER_NOTIFICATION = "get master user notification error"
	USER_NOTIFICATION_NOT_FOUND = "user notification is not found"
	NOTIF_TOKEN_NOT_FOUND       = "token is empty or not found"

	// push notif
	ERROR_SEND_PUSH_NOTIF = "error when send push notification"

	// telegram
	INIT_TELEGRAM             = "initialization telegram"
	ERROR_SEND_TELEGRAM_NOTIF = "error when send telegram notification"

	// file
	ERROR_REMOVE_FILE = "remove file error"
)
