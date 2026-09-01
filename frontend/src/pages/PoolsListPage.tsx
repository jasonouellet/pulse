import React from "react";
import { Plus, ChevronRight } from "lucide-react";
import { Button } from "../components/ui/Button";
import { MOCK_POOLS } from "../data/mock";

function getGenderLabel(gender: "MASCULINE" | "FEMININE" | "MIXED"): string {
  if (gender === "MIXED") return "Mixte";
  if (gender === "MASCULINE") return "Masculin";
  return "Féminin";
}

export const PoolsListPage: React.FC = () => {
  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
            Bassins
          </h1>
          <p className="mt-1 text-sm text-slate-500 dark:text-slate-400">
            Saison 2026 · Soccer
          </p>
        </div>
        <Button variant="primary">
          <Plus className="mr-1.5 h-4 w-4" aria-hidden="true" />
          Créer un bassin
        </Button>
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-200 dark:border-slate-800">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="p-3">Code</th>
              <th className="p-3">Nom</th>
              <th className="p-3">Âge</th>
              <th className="p-3">Genre</th>
              <th className="p-3">Statut</th>
              <th className="p-3" />
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {MOCK_POOLS.map((pool) => (
              <tr
                key={pool.id}
                className="hover:bg-slate-50 dark:hover:bg-slate-800/50"
              >
                <td className="p-3 font-medium text-slate-900 dark:text-white">
                  {pool.code}
                </td>
                <td className="p-3 text-slate-700 dark:text-slate-300">
                  {pool.name}
                </td>
                <td className="p-3 text-slate-700 dark:text-slate-300">
                  {pool.minAge}-{pool.maxAge}
                </td>
                <td className="p-3 text-slate-700 dark:text-slate-300">
                  {getGenderLabel(pool.gender)}
                </td>
                <td className="p-3">
                  <span
                    className={
                      pool.isActive
                        ? "rounded-md bg-brand-50 px-2 py-0.5 text-xs font-medium text-brand-600 dark:bg-brand-900/40 dark:text-brand-300"
                        : "rounded-md bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500 dark:bg-slate-800 dark:text-slate-400"
                    }
                  >
                    {pool.isActive ? "Actif" : "Inactif"}
                  </span>
                </td>
                <td className="p-3 text-right">
                  <ChevronRight
                    className="ml-auto h-4 w-4 text-slate-400"
                    aria-hidden="true"
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
