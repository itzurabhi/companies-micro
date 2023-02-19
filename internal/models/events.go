package models

type CompanyEventType int

const (
	CompanyCreatedEvent CompanyEventType = 1
	CompanyPatchedEvent CompanyEventType = 2
	CompanyDeletedEvent CompanyEventType = 3
)

type CompanyEvent struct {
	// type of event
	Type CompanyEventType

	// on which resource it was done
	Uuid string
}
