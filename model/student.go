package model

import (
	"net/http"

	"github.com/mholt/binding"
)

type Student struct {
	ID         string `gorethink:"id,omitempty"json:"id"`
	PersonID   string `gorethink:"personId,omitempty"json:"personId"`
	GradeLevel string `gorethink:"gradeLevel,omitempty"json:"gradeLevel"`
	TimeStamp
}

func (s Student) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if s.PersonID == "" {
		errs = append(errs, RequiredErr("personId"))
	}
	if s.GradeLevel == "" {
		errs = append(errs, RequiredErr("gradeLevel"))
	}
	return errs
}

func (s *Student) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&s.ID:         field("id", false),
		&s.PersonID:   field("personId", true),
		&s.GradeLevel: field("gradeLevel", true),
	}
}
