export const staff = {
  title: "Personnel",
  description: "Gérez les scanneurs et réceptionnistes de votre centre.",
  empty: "Aucun membre du personnel trouvé.",
  loading: "Chargement...",
  addButton: "Ajouter un membre",
  columns: {
    username: "Nom d'utilisateur",
    email: "E-mail",
    role: "Rôle",
    subCenter: "Sous-centre",
    createdAt: "Créé le",
    openMenu: "Ouvrir le menu",
    actions: "Actions",
  },
  filters: {
    search: "Rechercher",
    searchPlaceholder: "Rechercher par nom d'utilisateur ou e-mail...",
    role: "Rôle",
    allRoles: "Tous les rôles",
  },
  dialog: {
    addTitle: "Ajouter un membre",
    addDescription:
      "Créer un nouveau scanneur ou réceptionniste pour votre centre.",
    editTitle: "Modifier le membre",
    editDescription:
      "Mettre à jour les informations de ce membre du personnel.",
    submitAdd: "Créer le membre",
    submitEdit: "Enregistrer les modifications",
  },
  fields: {
    username: "Nom d'utilisateur",
    email: "E-mail",
    role: "Rôle",
    subCenter: "Sous-centre",
    subCenterPlaceholder: "Sélectionner un sous-centre...",
    password: "Mot de passe",
    passwordEditHint: "(laisser vide pour conserver l'actuel)",
  },
  deleteDialog: {
    title: "Supprimer le membre du personnel",
    description:
      "Voulez-vous vraiment supprimer {name} ? Cette action est irréversible.",
  },
  createSuccess: "Membre du personnel créé.",
  createSuccessDescription: "Le nouveau membre peut désormais se connecter.",
  updateSuccess: "Membre du personnel mis à jour.",
  updateSuccessDescription: "Les modifications ont été enregistrées.",
  deleteSuccess: "Membre du personnel supprimé.",
  deleteSuccessDescription: "Son accès a été révoqué.",
  errors: {
    centerMismatch:
      "Vous ne pouvez gérer que le personnel de votre propre centre.",
    notFound: "Membre du personnel introuvable.",
    loadFailed: "Impossible de charger le personnel. Veuillez réessayer.",
    createFailed:
      "Impossible de créer le membre du personnel. Veuillez réessayer.",
    updateFailed:
      "Impossible de mettre à jour le membre du personnel. Veuillez réessayer.",
    deleteFailed:
      "Impossible de supprimer le membre du personnel. Veuillez réessayer.",
    usernameRequired: "Le nom d'utilisateur est requis.",
    usernameMin: "Le nom d'utilisateur doit contenir au moins 3 caractères.",
    emailRequired: "L'e-mail est requis.",
    emailInvalid: "Veuillez saisir un e-mail valide.",
    emailTaken: "Cet e-mail est déjà utilisé.",
    passwordRequired: "Le mot de passe est requis.",
    passwordMin: "Le mot de passe doit contenir au moins 6 caractères.",
    roleInvalid: "Veuillez sélectionner un rôle valide.",
    subCenterRequired: "Veuillez sélectionner un sous-centre.",
  },
}
