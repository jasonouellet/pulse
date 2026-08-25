export interface SportLexicon {
  sportId: string;
  terms: {
    pool: string; // ex: Bassin d'âge / Catégorie
    roster: string; // ex: Équipe éphémère / Roster
    match: string; // ex: Match / Partie
    field: string; // ex: Terrain
    subField: string; // ex: Demi-terrain / Surface réduites
    coach: string; // ex: Entraîneur / Coach
    referee: string; // ex: Arbitre
  };
}

export const SOCCER_LEXICON: SportLexicon = {
  sportId: "soccer",
  terms: {
    pool: "Bassin d'âge",
    roster: "Équipe éphémère",
    match: "Match",
    field: "Terrain",
    subField: "Demi-terrain (7v7 / 9v9)",
    coach: "Éducateur / Entraîneur",
    referee: "Arbitre",
  },
};

export const DEFAULT_LEXICON = SOCCER_LEXICON;
