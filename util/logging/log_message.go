package logging

const (
	// general
	HTTP_REQUEST         = "http request"
	ERROR_UNMARSHAL_JSON = "unmarshal attributes error"

	// send notif hazard
	ERROR_DELETE_FILE_FROM_FTP_HAZARD                   = "error when delete file from ftp when post hazard in defer"
	ERROR_SEND_NOTIF_SUBMIT_NOT_NEED_FU_SM_SO_HAZARD    = "error when send notification data not need follow up to site manager and safety officer"
	ERROR_SEND_NOTIF_SUBMIT_NEED_FU_AREA_MANAGER_HAZARD = "error when send notification data need follow up to area manager"
	ERROR_SEND_NOTIF_CLOSED_EMPLOYEE_HAZARD             = "error when send notification data closed to employee"
	ERROR_SEND_NOTIF_CLOSED_SM_SO_HAZARD                = "error when send notification data closed to site manager and safety officer"

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
