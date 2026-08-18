package academic

import "gorm.io/gorm"

type subjectSeed = string

type majorSeed struct {
	name     string
	subjects []subjectSeed
}

type gradeSeed struct {
	name   string
	majors []majorSeed
}

// moroccanTree is a representative (not exhaustive) slice of the Moroccan
// lycée curriculum: two grades, each with two filières and their matières.
var moroccanTree = []gradeSeed{
	{
		name: "1ère Année Bac",
		majors: []majorSeed{
			{
				name:     "Sciences Expérimentales",
				subjects: []subjectSeed{"Mathématiques", "Physique-Chimie", "Sciences de la Vie et de la Terre", "Français", "Anglais"},
			},
			{
				name:     "Lettres et Sciences Humaines",
				subjects: []subjectSeed{"Arabe", "Français", "Philosophie", "Histoire-Géographie"},
			},
		},
	},
	{
		name: "2ème Année Bac",
		majors: []majorSeed{
			{
				name:     "Sciences Mathématiques A",
				subjects: []subjectSeed{"Mathématiques", "Physique-Chimie", "Sciences de la Vie et de la Terre", "Philosophie"},
			},
			{
				name:     "Sciences Physiques",
				subjects: []subjectSeed{"Mathématiques", "Physique-Chimie", "Sciences de la Vie et de la Terre", "Philosophie"},
			},
		},
	},
}

// SeedResult exposes the ids SeedAcademicStructure created, so other seeders
// (e.g. teachers) can attach real subject/class ids without re-deriving the
// tree.
type SeedResult struct {
	// SubjectIDs is keyed "Major > Subject".
	SubjectIDs map[string]uint
	// ClassIDs is keyed by major name — each major has exactly one
	// auto-generated class (see ClassNameForMajor).
	ClassIDs map[string]uint
}

// SeedAcademicStructure idempotently creates moroccanTree for the given
// center — grades, majors, subjects, and (via ClassNameForMajor, the same
// rule Service.AddMajor uses) one class per major.
func SeedAcademicStructure(db *gorm.DB, centerID uint) (SeedResult, error) {
	result := SeedResult{
		SubjectIDs: make(map[string]uint),
		ClassIDs:   make(map[string]uint),
	}

	for _, g := range moroccanTree {
		var grade Grade
		if err := db.Where("center_id = ? AND name = ?", centerID, g.name).
			FirstOrCreate(&grade, Grade{CenterID: centerID, Name: g.name}).Error; err != nil {
			return SeedResult{}, err
		}

		for _, m := range g.majors {
			var major Major
			if err := db.Where("grade_id = ? AND name = ?", grade.ID, m.name).
				FirstOrCreate(&major, Major{GradeID: grade.ID, Name: m.name}).Error; err != nil {
				return SeedResult{}, err
			}

			className := ClassNameForMajor(major.Name)
			var class Class
			if err := db.Where("major_id = ? AND name = ?", major.ID, className).
				FirstOrCreate(&class, Class{MajorID: major.ID, Name: className}).Error; err != nil {
				return SeedResult{}, err
			}
			result.ClassIDs[m.name] = class.ID

			for _, subName := range m.subjects {
				var subject Subject
				if err := db.Where("major_id = ? AND name = ?", major.ID, subName).
					FirstOrCreate(&subject, Subject{MajorID: major.ID, Name: subName}).Error; err != nil {
					return SeedResult{}, err
				}
				result.SubjectIDs[m.name+" > "+subName] = subject.ID
			}
		}
	}

	return result, nil
}
