package app_errors

type DatabaseConnError struct {
	Message string
}

func (e DatabaseConnError) Error() string {
	return ERROR_STRING + "Failed to connect to database: " + e.Message
}

func NewDatabaseConnError(message string) error {
	return DatabaseConnError{Message: message}
}

type DatabaseInsertionError struct {
	Message string
}

func (e DatabaseInsertionError) Error() string {
	return ERROR_STRING + "Failed to insert items to database: " + e.Message
}

func NewDatabaseInsertionError(message string) error {
	return DatabaseInsertionError{Message: message}
}
