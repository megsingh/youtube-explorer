package app_errors

const ERROR_STRING = "ERROR - "

type EnvironmentVariableError struct {
	Message string
}

func (e EnvironmentVariableError) Error() string {
	return ERROR_STRING + "Failed to load environment variables"
}

func NewEnvironmentVariableError(message string) error {
	return EnvironmentVariableError{Message: message}
}

type StartupError struct {
	Message string
}

func (e StartupError) Error() string {
	return ERROR_STRING + "Failed to start the application: " + e.Message
}

func NewStartupError(message string) error {
	return StartupError{Message: message}
}

type ServerError struct {
	Message string
}

func (e ServerError) Error() string {
	return ERROR_STRING + "Failed to run the server " + e.Message
}

func NewServerError(message string) error {
	return ServerError{Message: message}
}
