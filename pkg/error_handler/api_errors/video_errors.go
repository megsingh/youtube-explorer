package api_errors

type VideoFetchError struct {
	Message string
}

func (e VideoFetchError) Error() string {
	return ERROR_STRING + "Failed to fetch videos from the database " + e.Message
}

func NewVideoFetchError(message string) error {
	return VideoFetchError{Message: message}
}

type VideoInsertError struct {
	Message string
}

func (e VideoInsertError) Error() string {
	return ERROR_STRING + "Failed to insert videos into the database " + e.Message
}

func NewVideoInsertError(message string) error {
	return VideoInsertError{Message: message}
}
