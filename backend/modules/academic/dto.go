package academic

// SubjectResponse is the API shape of a single subject.
type SubjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ClassResponse is the API shape of a major's one auto-generated class.
type ClassResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// MajorResponse is the API shape of a major, its subjects, and its one
// auto-generated class.
type MajorResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Subjects []SubjectResponse `json:"subjects"`
	Class    *ClassResponse    `json:"class"`
}

// GradeResponse is the API shape of a grade and its full major/subject tree.
type GradeResponse struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Majors []MajorResponse `json:"majors"`
}

// TreeResponse is the center owner's whole academic structure, used to
// render (and re-hydrate) the onboarding wizard.
type TreeResponse struct {
	Grades []GradeResponse `json:"grades"`
}

// AddGradeRequest creates a grade in the caller's own center.
type AddGradeRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// AddMajorRequest creates a major under a grade in the caller's own center.
type AddMajorRequest struct {
	Name string `json:"name" binding:"required,min=1,max=150"`
}

// AddSubjectRequest creates a subject under a major in the caller's own
// center.
type AddSubjectRequest struct {
	Name string `json:"name" binding:"required,min=1,max=150"`
}
