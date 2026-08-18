export const onboarding = {
  title: "Configuration académique",
  description:
    "Créez chaque niveau que vous enseignez, ajoutez les filières proposées dans ce niveau, puis listez les matières enseignées dans chaque filière.",
  addGrade: "Ajouter un niveau",
  addGradePlaceholder: "Rechercher ou créer un niveau…",
  addMajor: "Ajouter une filière",
  addMajorPlaceholder: "Rechercher ou créer une filière…",
  addSubject: "Ajouter une matière",
  addSubjectPlaceholder: "Rechercher ou créer une matière…",
  continueButton: "Continuer vers le tableau de bord",
  incompleteStructureHint:
    "Chaque niveau doit avoir au moins une filière, et chaque filière au moins une matière, avant de pouvoir continuer.",
  loading: "Chargement de votre configuration académique...",
  suggestions: "Suggestions",
  createNew: "Créer",
  createOption: 'Créer "{name}"',
  alreadyAdded: '"{name}" a déjà été ajouté.',
  allSuggestionsAdded:
    "Toutes les suggestions ont été ajoutées — tapez pour en créer une nouvelle.",
  removeLabel: "Supprimer",
  noMajors:
    "Aucune filière pour le moment — ajoutez les filières enseignées dans ce niveau.",
  noSubjects: "Aucune matière pour le moment.",
  empty: {
    title: "Aucun niveau pour le moment",
    description:
      "Commencez par ajouter un niveau — choisissez parmi les niveaux courants ou créez le vôtre.",
  },
  removeGrade: {
    title: 'Supprimer "{name}" ?',
    description:
      "Cela supprimera également toutes ses filières et matières. Cette action est irréversible.",
    confirm: "Supprimer le niveau",
  },
  gradeAddedSuccess: "Niveau ajouté.",
  gradeAddedSuccessDescription: "Votre structure académique a été mise à jour.",
  gradeRemovedSuccess: "Niveau supprimé.",
  gradeRemovedSuccessDescription:
    "Le niveau et tout ce qu'il contenait ont été supprimés.",
  majorAddedSuccess: "Filière ajoutée.",
  majorAddedSuccessDescription: "Votre structure académique a été mise à jour.",
  majorRemovedSuccess: "Filière supprimée.",
  majorRemovedSuccessDescription:
    "La filière et ses matières ont été supprimées.",
  errors: {
    treeFailed:
      "Impossible de charger votre configuration académique. Veuillez réessayer.",
    addGradeFailed: "Impossible d'ajouter le niveau. Veuillez réessayer.",
    removeGradeFailed: "Impossible de supprimer le niveau. Veuillez réessayer.",
    addMajorFailed: "Impossible d'ajouter la filière. Veuillez réessayer.",
    removeMajorFailed:
      "Impossible de supprimer la filière. Veuillez réessayer.",
    addSubjectFailed: "Impossible d'ajouter la matière. Veuillez réessayer.",
    removeSubjectFailed:
      "Impossible de supprimer la matière. Veuillez réessayer.",
    gradeExists: "Ce niveau existe déjà.",
    majorExists: "Cette filière existe déjà.",
    subjectExists: "Cette matière existe déjà.",
    gradeNotFound: "Niveau introuvable.",
    majorNotFound: "Filière introuvable.",
    subjectNotFound: "Matière introuvable.",
    setupRequired:
      "Terminez la configuration de votre structure académique avant de continuer.",
  },
}
