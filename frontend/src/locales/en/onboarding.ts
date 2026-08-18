export const onboarding = {
  title: "Academic setup",
  description:
    "Create each grade you teach, add the majors offered in that grade, then list the subjects taught in every major.",
  addGrade: "Add grade",
  addGradePlaceholder: "Search or create a grade…",
  addMajor: "Add major",
  addMajorPlaceholder: "Search or create a major…",
  addSubject: "Add subject",
  addSubjectPlaceholder: "Search or create a subject…",
  continueButton: "Continue to dashboard",
  incompleteStructureHint:
    "Every grade needs at least one major, and every major needs at least one subject, before you can continue.",
  loading: "Loading your academic setup...",
  suggestions: "Suggestions",
  createNew: "Create new",
  createOption: 'Create "{name}"',
  alreadyAdded: '"{name}" has already been added.',
  allSuggestionsAdded: "All suggestions were added — type to create a new one.",
  removeLabel: "Remove",
  noMajors: "No majors yet — add the tracks taught in this grade.",
  noSubjects: "No subjects yet.",
  empty: {
    title: "No grades yet",
    description:
      "Start by adding a grade — you can pick from common grades or create your own.",
  },
  removeGrade: {
    title: 'Remove "{name}"?',
    description:
      "This will also delete all of its majors and subjects. This action cannot be undone.",
    confirm: "Remove grade",
  },
  gradeAddedSuccess: "Grade added.",
  gradeAddedSuccessDescription: "Your academic structure has been updated.",
  gradeRemovedSuccess: "Grade removed.",
  gradeRemovedSuccessDescription:
    "The grade and everything under it has been deleted.",
  majorAddedSuccess: "Major added.",
  majorAddedSuccessDescription: "Your academic structure has been updated.",
  majorRemovedSuccess: "Major removed.",
  majorRemovedSuccessDescription:
    "The major and its subjects have been deleted.",
  errors: {
    treeFailed: "Failed to load your academic setup. Please try again.",
    addGradeFailed: "Failed to add grade. Please try again.",
    removeGradeFailed: "Failed to remove grade. Please try again.",
    addMajorFailed: "Failed to add major. Please try again.",
    removeMajorFailed: "Failed to remove major. Please try again.",
    addSubjectFailed: "Failed to add subject. Please try again.",
    removeSubjectFailed: "Failed to remove subject. Please try again.",
    gradeExists: "This grade already exists.",
    majorExists: "This major already exists.",
    subjectExists: "This subject already exists.",
    gradeNotFound: "Grade not found.",
    majorNotFound: "Major not found.",
    subjectNotFound: "Subject not found.",
    setupRequired:
      "Finish setting up your academic structure before continuing.",
  },
}
