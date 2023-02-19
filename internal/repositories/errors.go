package repositories

type DatabaseError struct {
	internal    error
	userMessage string
}

func (dbErr DatabaseError) Error() string {
	return dbErr.userMessage
}

func (dbErr DatabaseError) Cause() error {
	return dbErr.internal
}
