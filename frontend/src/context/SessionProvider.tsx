import React, { useMemo, useState } from "react";
import {
  SessionContext,
  type Club,
  type RoleGrant,
  type SessionContextValue,
  type UserRole,
} from "./session";

// Mock data for this first draft. Replace with a real session fetched from
// something like GET /api/v1/core/me, which should return the user plus
// their full core.user_roles grant list (club + role pairs).
//
// IMPORTANT: activeGrant here is a client-side display/filtering choice
// only. It must never be trusted as an authorization boundary — every
// backend endpoint has to independently verify the caller holds the
// required role for the target club via the session token, regardless of
// what "active" context the client claims to be in.
const MOCK_CLUBS: Club[] = [
  { id: "club-fc-rimouski", name: "FC Rimouski" },
  { id: "club-hockey-rimouski", name: "Hockey Rimouski Élite" },
];

const MOCK_GRANTS: RoleGrant[] = [
  { club: MOCK_CLUBS[0], role: "GUARDIAN" },
  { club: MOCK_CLUBS[0], role: "COACH" },
  { club: MOCK_CLUBS[1], role: "GUARDIAN" },
];

export const SessionProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [activeGrant, setActiveGrant] = useState<RoleGrant>(MOCK_GRANTS[0]);

  const clubs = useMemo(
    () =>
      Array.from(new Map(MOCK_GRANTS.map((g) => [g.club.id, g.club])).values()),
    [],
  );

  const setActiveRole = (role: UserRole) => {
    const grant = MOCK_GRANTS.find(
      (g) => g.club.id === activeGrant.club.id && g.role === role,
    );
    if (grant) setActiveGrant(grant);
  };

  const setActiveClub = (clubId: string) => {
    // When switching clubs, land on that club's first available role
    // rather than trying to preserve the previous role, since the person
    // may not hold the same role there.
    const grant = MOCK_GRANTS.find((g) => g.club.id === clubId);
    if (grant) setActiveGrant(grant);
  };

  const value = useMemo<SessionContextValue>(
    () => ({
      userName: "Marie Tremblay",
      userInitials: "MT",
      grants: MOCK_GRANTS,
      clubs,
      activeGrant,
      setActiveRole,
      setActiveClub,
    }),
    [activeGrant, clubs],
  );

  return (
    <SessionContext.Provider value={value}>{children}</SessionContext.Provider>
  );
};
