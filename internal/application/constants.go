package applicationpackage application

import "time"

const (
	PROFILE_DEFAULT    string = "default"
	PROFILE_OPTION     string = "profile"
	PROPERTY_PATH      string = "res"
	PROPERTY_FILE_TYPE string = "env"

	ROTATE_INTERVAL       time.Duration = 24 * time.Hour
	ROTATE_MAX_AGE        time.Duration = ROTATE_INTERVAL * 7
	ROTATE_LINK_EXTENSION string        = ".log"
	ROTATE_FILE_EXTENSION string        = ".%Y-%m-%d.log"
)

