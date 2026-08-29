import type { UserRole } from "../context/session";

export interface Pool {
  id: string;
  code: string;
  name: string;
  minAge: number;
  maxAge: number;
  gender: "MASCULINE" | "FEMININE" | "MIXED";
  division: string;
  seasonYear: number;
  isActive: boolean;
}

export interface Position {
  code: string;
  name: string;
}

export interface Player {
  id: string;
  firstName: string;
  lastName: string;
  birthDate: string; // ISO date — identity data, reused every season
  poolCode: string;
  score: number; // staff-only — never expose to GUARDIAN/PLAYER roles
  preferredPositions: Position[];
  registrationStatus: "PENDING" | "ASSIGNED";
}

export interface Roster {
  id: string;
  name: string;
  type: "TRAINING_GROUP" | "SEASON_TEAM" | "EVENT_ROSTER";
  playerIds: string[];
}

export const MOCK_POSITIONS: Record<string, Position> = {
  ATT: { code: "ATT", name: "Attaquant" },
  DEF: { code: "DEF", name: "Défense" },
  MIL: { code: "MIL", name: "Demi" },
  GAR: { code: "GAR", name: "Gardien" },
};

export const MOCK_POOLS: Pool[] = [
  {
    id: "pool-u10-d1",
    code: "U10_D1",
    name: "U9-U10 Division 1",
    minAge: 9,
    maxAge: 10,
    gender: "MIXED",
    division: "Division 1",
    seasonYear: 2026,
    isActive: true,
  },
  {
    id: "pool-u12m-loc",
    code: "U12M_LOC",
    name: "U12 Masculin Local",
    minAge: 11,
    maxAge: 12,
    gender: "MASCULINE",
    division: "Local",
    seasonYear: 2026,
    isActive: true,
  },
  {
    id: "pool-u14f-d1",
    code: "U14F_D1",
    name: "U14 Féminin Division 1",
    minAge: 13,
    maxAge: 14,
    gender: "FEMININE",
    division: "Division 1",
    seasonYear: 2026,
    isActive: true,
  },
  {
    id: "pool-u16-d2",
    code: "U16_D2",
    name: "U16 Mixte Division 2",
    minAge: 15,
    maxAge: 16,
    gender: "MIXED",
    division: "Division 2",
    seasonYear: 2026,
    isActive: false,
  },
];

export const MOCK_PLAYERS: Player[] = [
  {
    id: "player-leo",
    firstName: "Léo",
    lastName: "Tremblay",
    birthDate: "2016-03-14",
    poolCode: "U10_D1",
    score: 78,
    preferredPositions: [MOCK_POSITIONS.ATT, MOCK_POSITIONS.MIL],
    registrationStatus: "ASSIGNED",
  },
  {
    id: "player-zoe",
    firstName: "Zoé",
    lastName: "Tremblay",
    birthDate: "2012-07-02",
    poolCode: "U14F_D1",
    score: 65,
    preferredPositions: [MOCK_POSITIONS.DEF],
    registrationStatus: "PENDING",
  },
  {
    id: "player-nora",
    firstName: "Nora",
    lastName: "Bissonnette",
    birthDate: "2017-01-20",
    poolCode: "U10_D1",
    score: 65,
    preferredPositions: [MOCK_POSITIONS.GAR],
    registrationStatus: "ASSIGNED",
  },
  {
    id: "player-emile",
    firstName: "Émile",
    lastName: "Roy",
    birthDate: "2016-09-11",
    poolCode: "U10_D1",
    score: 71,
    preferredPositions: [MOCK_POSITIONS.MIL],
    registrationStatus: "ASSIGNED",
  },
  {
    id: "player-ana",
    firstName: "Ana",
    lastName: "Dubé",
    birthDate: "2017-05-30",
    poolCode: "U10_D1",
    score: 58,
    preferredPositions: [MOCK_POSITIONS.DEF],
    registrationStatus: "PENDING",
  },
  {
    id: "player-sam",
    firstName: "Sam",
    lastName: "Bouchard",
    birthDate: "2016-11-08",
    poolCode: "U10_D1",
    score: 70,
    preferredPositions: [MOCK_POSITIONS.ATT],
    registrationStatus: "ASSIGNED",
  },
];

export const MOCK_ROSTERS: Roster[] = [
  {
    id: "roster-u10-training",
    name: "U10-M Division 1",
    type: "TRAINING_GROUP",
    playerIds: ["player-leo", "player-nora", "player-emile"],
  },
  {
    id: "roster-faucons-season",
    name: "Faucons U9-U10",
    type: "SEASON_TEAM",
    playerIds: ["player-leo", "player-emile", "player-sam"],
  },
  {
    id: "roster-tournoi-a",
    name: "Roster A · Tournoi Rimouski #3",
    type: "EVENT_ROSTER",
    playerIds: ["player-leo"],
  },
];

// In the real app this would come from core.parents_children for the
// logged-in guardian, not a hardcoded list.
export const MY_CHILDREN_IDS = ["player-leo", "player-zoe"];

// Nav items differ by active role — extend as real screens land per role.
export interface NavItem {
  label: string;
  to: string;
}

export const NAV_BY_ROLE: Record<UserRole, NavItem[]> = {
  GUARDIAN: [
    { label: "Accueil", to: "/" },
    { label: "Mes enfants", to: "/children" },
    { label: "Inscription", to: "/register" },
  ],
  PLAYER: [{ label: "Accueil", to: "/" }],
  COACH: [
    { label: "Accueil", to: "/" },
    { label: "Bassins", to: "/pools" },
    { label: "Rosters", to: "/rosters" },
  ],
  CLUB_ADMIN: [
    { label: "Accueil", to: "/" },
    { label: "Bassins", to: "/pools" },
    { label: "Rosters", to: "/rosters" },
  ],
  TECHNICAL_DIRECTOR: [
    { label: "Accueil", to: "/" },
    { label: "Bassins", to: "/pools" },
    { label: "Rosters", to: "/rosters" },
  ],
};
