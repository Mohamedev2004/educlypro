package academic

import "time"

// Grade is a teaching level a center offers (e.g. "Grade 10"), scoped to a
// single center. Deleting a grade cascades to its majors and, transitively,
// their subjects.
type Grade struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	CenterID uint   `json:"center_id" gorm:"not null;uniqueIndex:idx_grade_center_name"`
	Name     string `json:"name" gorm:"size:100;not null;uniqueIndex:idx_grade_center_name"`

	Majors []Major `json:"majors,omitempty" gorm:"foreignKey:GradeID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Major is a track offered within a grade (e.g. "Science"). Deleting a major
// cascades to its subjects and its class.
type Major struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	GradeID uint   `json:"grade_id" gorm:"not null;uniqueIndex:idx_major_grade_name"`
	Name    string `json:"name" gorm:"size:150;not null;uniqueIndex:idx_major_grade_name"`

	Subjects []Subject `json:"subjects,omitempty" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
	Classes  []Class   `json:"classes,omitempty" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Subject is taught within a major (e.g. "Physics").
type Subject struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	MajorID uint   `json:"major_id" gorm:"not null;uniqueIndex:idx_subject_major_name"`
	Name    string `json:"name" gorm:"size:150;not null;uniqueIndex:idx_subject_major_name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Class is the group of students taking a major (e.g. "Science class").
// Exactly one is created automatically alongside its major (see
// Service.AddMajor / ClassNameForMajor) — there is no separate
// create-a-class action.
type Class struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	MajorID uint   `json:"major_id" gorm:"not null;uniqueIndex:idx_class_major_name"`
	Name    string `json:"name" gorm:"size:150;not null;uniqueIndex:idx_class_major_name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
