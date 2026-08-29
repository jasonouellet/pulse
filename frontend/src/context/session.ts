import { createContext, useContext } from "react";

export interface Club {
  id: string;
  name: string;
}

export type UserRole =
  "CLUB_ADMIN" | "TECHNICAL_DIRECTOR" | "COACH" | "GUARDIAN" | "PLAYER";

export const ROLE_LABELS: Record<UserRole, string> = {
  CLUB_ADMIN: "Gestionnaire du club",
  TECHNICAL_DIRECTOR: "Directeur technique",
  COACH: "Entraîneur",
  GUARDIAN: "Responsable",
  PLAYER: "Sportif",
};

export interface RoleGrant {
  club: Club;
  role: UserRole;
}

export interface SessionContextValue {
  userName: string;
  userInitials: string;
  grants: RoleGrant[];
  clubs: Club[];
  activeGrant: RoleGrant;
  setActiveRole: (role: UserRole) => void;
  setActiveClub: (clubId: string) => void;
}

export const SessionContext = createContext<SessionContextValue | undefined>(
  undefined,
);

export function useSession(): SessionContextValue {
  const ctx = useContext(SessionContext);
  if (!ctx) {
    throw new Error("useSession must be used within a SessionProvider");
  }
  return ctx;
}
