package api_errors

const ERROR_STRING = "ERROR - "

type QuotaExceedError struct {
	Message string
}

func (e QuotaExceedError) Error() string {
	return ERROR_STRING + "The request cannot be completed because you have exceeded your quota " + e.Message
}

func NewQuotaExceedError(message string) error {
	return QuotaExceedError{Message: message}
}

type YoutubeAPIError struct {
	Message string
}

func (e YoutubeAPIError) Error() string {
	return ERROR_STRING + "Failed to fetch videos from Youtube" + e.Message
}

func NewYoutubeAPIError(message string) error {
	return YoutubeAPIError{Message: message}
}
