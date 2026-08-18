/**
 * Academic setup related constants.
 *
 * Responsibility: Static suggestion lists for the grade/major/subject
 * onboarding wizard, based on the Moroccan education system (primaire ->
 * collège -> lycée, with the Bac filières at the end).
 * Layer: Constants
 */

export const GRADE_SUGGESTIONS = [
  "1ère Année Primaire",
  "2ème Année Primaire",
  "3ème Année Primaire",
  "4ème Année Primaire",
  "5ème Année Primaire",
  "6ème Année Primaire",
  "1ère Année Collège",
  "2ème Année Collège",
  "3ème Année Collège",
  "Tronc Commun",
  "1ère Année Bac",
  "2ème Année Bac",
] as const

export const MAJOR_SUGGESTIONS = [
  "Enseignement Général",
  "Tronc Commun Sciences",
  "Tronc Commun Lettres et Sciences Humaines",
  "Sciences Mathématiques A",
  "Sciences Mathématiques B",
  "Sciences Physiques",
  "Sciences de la Vie et de la Terre",
  "Sciences Expérimentales",
  "Sciences Économiques",
  "Sciences de Gestion Comptable",
  "Lettres et Sciences Humaines",
  "Arts Appliqués",
] as const

export const SUBJECT_SUGGESTIONS = [
  "Mathématiques",
  "Physique-Chimie",
  "Sciences de la Vie et de la Terre",
  "Arabe",
  "Français",
  "Anglais",
  "Espagnol",
  "Philosophie",
  "Histoire-Géographie",
  "Éducation Islamique",
  "Éducation Physique",
  "Informatique",
  "Économie et Organisation Administrative des Entreprises",
  "Comptabilité",
] as const
