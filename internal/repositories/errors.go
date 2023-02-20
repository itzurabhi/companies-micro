package repositories

type DatabaseError struct {
	userMessage string
}

func (dbErr DatabaseError) Error() string {
	return dbErr.userMessage
}

var ErrorRecordNotFound = DatabaseError{
	userMessage: "record not found",
}

var ErrorRecordAlreadyExist = DatabaseError{
	userMessage: "record already exist",
}
