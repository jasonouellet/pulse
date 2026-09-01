import React, { useMemo, useState } from "react";
import { ArrowRight, X, Plus } from "lucide-react";
import { Link } from "react-router";
import { Button } from "../components/ui/Button";
import { MOCK_PLAYERS, MOCK_POOLS, type Player } from "../data/mock";

export const RosterBuilderPage: React.FC = () => {
  const [activePoolCode, setActivePoolCode] = useState(MOCK_POOLS[0].code);
  const [rosterPlayerIds, setRosterPlayerIds] = useState<string[]>([]);

  const pool = MOCK_POOLS.find((p) => p.code === activePoolCode);
  const poolPlayers = useMemo(
    () => MOCK_PLAYERS.filter((p) => p.poolCode === activePoolCode),
    [activePoolCode],
  );
  const available = poolPlayers.filter((p) => !rosterPlayerIds.includes(p.id));
  const rosterPlayers = poolPlayers.filter((p) =>
    rosterPlayerIds.includes(p.id),
  );
  const avgScore = rosterPlayers.length
    ? Math.round(
        rosterPlayers.reduce((sum, p) => sum + p.score, 0) /
          rosterPlayers.length,
      )
    : null;

  const addToRoster = (id: string) =>
    setRosterPlayerIds((prev) => [...prev, id]);
  const removeFromRoster = (id: string) =>
    setRosterPlayerIds((prev) => prev.filter((pid) => pid !== id));

  const renderPlayerRow = (
    p: Player,
    action: { icon: React.ReactNode; onClick: () => void; label: string },
  ) => (
    <div
      key={p.id}
      className="flex items-center justify-between border-b border-slate-200 px-3 py-2 last:border-b-0 dark:border-slate-800"
    >
      <div>
        <p className="text-sm text-slate-900 dark:text-white">
          <Link
            to={`/players/${p.id}`}
            className="hover:underline focus:outline-none focus:ring-2 focus:ring-brand-500 rounded"
          >
            {p.firstName} {p.lastName}
          </Link>
        </p>
        <p className="text-xs text-slate-500 dark:text-slate-400">
          {p.preferredPositions.map((pos) => pos.name).join(", ")}
        </p>
      </div>
      <div className="flex items-center gap-3">
        <span className="rounded-md bg-brand-50 px-2 py-0.5 text-xs font-medium text-brand-700 dark:bg-brand-900/40 dark:text-brand-300">
          {p.score}
        </span>
        <button
          type="button"
          onClick={action.onClick}
          aria-label={action.label}
          className="text-slate-400 hover:text-slate-700 dark:hover:text-slate-200"
        >
          {action.icon}
        </button>
      </div>
    </div>
  );

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
          Formation des rosters
        </h1>
        <div className="mt-2 flex items-center gap-2">
          <select
            value={activePoolCode}
            onChange={(e) => {
              setActivePoolCode(e.target.value);
              setRosterPlayerIds([]);
            }}
            className="min-h-touch rounded-lg border border-slate-300 bg-white px-3 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-white"
          >
            {MOCK_POOLS.map((p) => (
              <option key={p.code} value={p.code}>
                {p.name}
              </option>
            ))}
          </select>
          <span className="rounded-md bg-brand-50 px-2.5 py-1 text-xs font-medium text-brand-700 dark:bg-brand-900/40 dark:text-brand-300">
            {poolPlayers.length} joueurs · Saison {pool?.seasonYear}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <p className="mb-2 text-sm font-medium text-slate-600 dark:text-slate-300">
            Disponibles ({available.length})
          </p>
          <div className="min-h-[80px] rounded-xl border border-slate-200 dark:border-slate-800">
            {available.length === 0 && (
              <p className="p-3 text-sm text-slate-400">Tous assignés</p>
            )}
            {available.map((p) =>
              renderPlayerRow(p, {
                icon: <ArrowRight className="h-4 w-4" aria-hidden="true" />,
                onClick: () => addToRoster(p.id),
                label: `Ajouter ${p.firstName} au roster`,
              }),
            )}
          </div>
        </div>

        <div>
          <div className="mb-2 flex items-center justify-between">
            <p className="text-sm font-medium text-brand-600 dark:text-brand-400">
              Roster ({rosterPlayers.length})
            </p>
            {avgScore !== null && (
              <span className="text-xs text-slate-500 dark:text-slate-400">
                Score moyen : {avgScore}
              </span>
            )}
          </div>
          <div className="min-h-[80px] rounded-xl border border-brand-200 dark:border-brand-800">
            {rosterPlayers.length === 0 && (
              <p className="p-3 text-sm text-slate-400">Aucun joueur assigné</p>
            )}
            {rosterPlayers.map((p) =>
              renderPlayerRow(p, {
                icon: <X className="h-4 w-4" aria-hidden="true" />,
                onClick: () => removeFromRoster(p.id),
                label: `Retirer ${p.firstName} du roster`,
              }),
            )}
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between">
        <Button variant="outline">
          <Plus className="mr-1.5 h-4 w-4" aria-hidden="true" />
          Nouveau roster
        </Button>
        <Button variant="primary">Enregistrer les rosters</Button>
      </div>
    </div>
  );
};
