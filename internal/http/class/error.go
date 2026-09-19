package class

import "errors"

var (
	ErrClassNotFound   = errors.New("class not found")
	ErrStudentNotFound = errors.New("student not found")
	ErrNotClassTeacher = errors.New("not class teacher")
)