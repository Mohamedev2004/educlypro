export const subcenters = {
  title: "Sous-centres",
  description: "Gérez les emplacements opérationnels au sein de votre centre.",
  empty: "Aucun sous-centre trouvé.",
  loading: "Chargement...",
  addButton: "Ajouter un sous-centre",
  columns: {
    name: "Nom",
    staffCount: "Personnel",
    createdAt: "Créé le",
    openMenu: "Ouvrir le menu",
    actions: "Actions",
  },
  dialog: {
    addTitle: "Ajouter un sous-centre",
    addDescription:
      "Créer un nouvel emplacement opérationnel au sein de votre centre.",
    editTitle: "Modifier le sous-centre",
    editDescription: "Renommer ce sous-centre.",
    submitAdd: "Créer le sous-centre",
    submitEdit: "Enregistrer les modifications",
  },
  fields: {
    name: "Nom",
  },
  deleteDialog: {
    title: "Supprimer le sous-centre",
    description:
      "Voulez-vous vraiment supprimer {name} ? Cette action est irréversible.",
  },
  createSuccess: "Sous-centre créé.",
  createSuccessDescription:
    "Il est désormais disponible lors de l'affectation du personnel.",
  updateSuccess: "Sous-centre mis à jour.",
  updateSuccessDescription: "Les modifications ont été enregistrées.",
  deleteSuccess: "Sous-centre supprimé.",
  deleteSuccessDescription:
    "Il n'est plus disponible pour l'affectation du personnel.",
  errors: {
    notFound: "Sous-centre introuvable.",
    exists: "Un sous-centre portant ce nom existe déjà.",
    hasStaff:
      "Réaffectez ou retirez son personnel avant de supprimer ce sous-centre.",
    mismatch: "Ce sous-centre n'appartient pas à ce centre.",
    loadFailed: "Impossible de charger les sous-centres. Veuillez réessayer.",
    createFailed: "Impossible de créer le sous-centre. Veuillez réessayer.",
    updateFailed:
      "Impossible de mettre à jour le sous-centre. Veuillez réessayer.",
    deleteFailed: "Impossible de supprimer le sous-centre. Veuillez réessayer.",
    nameRequired: "Le nom est requis.",
    nameMin: "Le nom doit contenir au moins 2 caractères.",
  },
}
