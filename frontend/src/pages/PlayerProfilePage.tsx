import React from "react";
import { useParams } from "react-router";
import { Users2, Trophy, ClipboardList, Award } from "lucide-react";
import { useSession } from "../context/session";
import { MOCK_PLAYERS, MOCK_ROSTERS, MOCK_POOLS } from "../data/mock";

const ROSTER_TYPE_LABELS = {
  TRAINING_GROUP: "Groupe d'entraînement",
  SEASON_TEAM: "Équipe de saison",
  EVENT_ROSTER: "Roster de tournoi",
} as const;

const ROSTER_TYPE_ICONS = {
  TRAINING_GROUP: Users2,
  SEASON_TEAM: ClipboardList,
  EVENT_ROSTER: Trophy,
} as const;

export const PlayerProfilePage: React.FC = () => {
  const { playerId } = useParams<{ playerId: string }>();
  const { activeGrant } = useSession();

  const player = MOCK_PLAYERS.find((p) => p.id === playerId) ?? MOCK_PLAYERS[0];
  const pool = MOCK_POOLS.find((p) => p.code === player.poolCode);
  const rosters = MOCK_ROSTERS.filter((r) => r.playerIds.includes(player.id));

  // Score is staff-only (COACH/CLUB_ADMIN/TECHNICAL_DIRECTOR). GUARDIAN and
  // PLAYER never see the numeric value — later this becomes a link to
  // "acquis" (badges/objectives) instead.
  const canSeeScore = ["COACH", "CLUB_ADMIN", "TECHNICAL_DIRECTOR"].includes(
    activeGrant.role,
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-brand-100 text-lg font-medium text-brand-700 dark:bg-brand-900/40 dark:text-brand-300">
          {player.firstName[0]}
          {player.lastName[0]}
        </div>
        <div className="flex-1">
          <h1 className="text-xl font-bold text-slate-900 dark:text-white">
            {player.firstName} {player.lastName}
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            {pool?.name}
          </p>
        </div>
        {canSeeScore ? (
          <div className="rounded-xl bg-slate-50 px-4 py-2 text-center dark:bg-slate-800">
            <p className="text-xs text-slate-400">Score</p>
            <p className="text-xl font-bold text-brand-600 dark:text-brand-400">
              {player.score}
            </p>
          </div>
        ) : (
          <div className="flex items-center gap-1.5 rounded-xl bg-slate-50 px-4 py-2 text-xs font-medium text-slate-500 dark:bg-slate-800 dark:text-slate-400">
            <Award className="h-4 w-4" aria-hidden="true" />
            Acquis
          </div>
        )}
      </div>

      <div>
        <p className="mb-2 text-sm font-medium text-slate-600 dark:text-slate-300">
          Positions préférées
        </p>
        <div className="flex gap-2">
          {player.preferredPositions.map((pos, i) => (
            <span
              key={pos.code}
              className={
                i === 0
                  ? "rounded-md bg-brand-50 px-2.5 py-1 text-xs font-medium text-brand-700 dark:bg-brand-900/40 dark:text-brand-300"
                  : "rounded-md bg-slate-100 px-2.5 py-1 text-xs text-slate-500 dark:bg-slate-800 dark:text-slate-400"
              }
            >
              {pos.name} {i === 0 ? "· 1er choix" : "· 2e choix"}
            </span>
          ))}
        </div>
      </div>

      <div>
        <p className="mb-2 text-sm font-medium text-slate-600 dark:text-slate-300">
          Rosters actifs cette saison
        </p>
        <div className="overflow-hidden rounded-xl border border-slate-200 dark:border-slate-800">
          {rosters.length === 0 && (
            <p className="p-3 text-sm text-slate-400">
              Aucun roster actif pour l'instant
            </p>
          )}
          {rosters.map((r, i) => {
            const Icon = ROSTER_TYPE_ICONS[r.type];
            return (
              <div
                key={r.id}
                className={`flex items-center justify-between px-3 py-2.5 ${
                  i < rosters.length - 1
                    ? "border-b border-slate-200 dark:border-slate-800"
                    : ""
                }`}
              >
                <div>
                  <p className="text-sm text-slate-900 dark:text-white">
                    {r.name}
                  </p>
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    {ROSTER_TYPE_LABELS[r.type]}
                  </p>
                </div>
                <Icon className="h-4 w-4 text-slate-400" aria-hidden="true" />
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
