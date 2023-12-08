package main

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type UUID = uuid.UUID

type ContactDto struct {
	Id    UUID
	Name  string
	Email string
}

type ContactsDto = []ContactDto

func newContactDto(name string, email string) *ContactDto {
	return &ContactDto{
		Id:    uuid.New(),
		Name:  name,
		Email: email,
	}
}

type ContactFormModel struct {
	Values map[string]string
	Errors map[string]string
}

type IndexViewModel struct {
	Count          int
	Contacts       ContactsDto
	ContactsErrors []string
	ContactForm    ContactFormModel
	Constants      map[string]string
}

func newIndexViewModel() *IndexViewModel {
	return &IndexViewModel{
		Count:       0,
		ContactForm: *newContactFormModel(),
		Contacts: []ContactDto{
			*newContactDto("John Doe", "john.doe@gmail.com"),
			*newContactDto("Jane Doe", "jane.doe@gmail.com"),
		},
		ContactsErrors: []string{},
		Constants: map[string]string{
			"EmptyUuid": uuid.Nil.String(),
		},
	}
}

func (vm *IndexViewModel) hasEmail(email string) bool {
	_, _, ok := lo.FindIndexOf(vm.Contacts, func(c ContactDto) bool {
		return c.Email == email
	})
	return ok
}

func (vm *IndexViewModel) findContactIndex(contactId UUID) (*ContactDto, int) {
	contact, idx, _ := lo.FindIndexOf(vm.Contacts, func(c ContactDto) bool {
		return bytes.Equal(contactId[:], c.Id[:])
	})
	if idx < 0 {
		return nil, idx
	}
	return &contact, idx
}

func newContactFormModel() *ContactFormModel {
	return &ContactFormModel{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}
