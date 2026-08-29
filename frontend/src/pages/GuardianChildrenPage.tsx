import React, { useState } from "react";
import { Plus, Check } from "lucide-react";
import { MOCK_PLAYERS, MOCK_POOLS, MY_CHILDREN_IDS } from "../data/mock";

export const GuardianChildrenPage: React.FC = () => {
  const children = MOCK_PLAYERS.filter((p) => MY_CHILDREN_IDS.includes(p.id));
  const [selectedIds, setSelectedIds] = useState<string[]>([children[0].id]);

  const toggle = (id: string) =>
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id],
    );

  const selected = children.filter((c) => selectedIds.includes(c.id));

  return (
    <div className="mx-auto max-w-md space-y-5">
      <div>
        <p className="mb-2 text-sm font-medium text-slate-600 dark:text-slate-300">
          Mes enfants
        </p>
        <div className="flex gap-2 overflow-x-auto">
          {children.map((child) => {
            const isSelected = selectedIds.includes(child.id);
            return (
              <button
                key={child.id}
                type="button"
                onClick={() => toggle(child.id)}
                className={`flex min-h-[44px] shrink-0 items-center gap-2 rounded-lg border px-3 py-1.5 text-sm ${
                  isSelected
                    ? "border-2 border-brand-500 bg-white dark:bg-slate-800"
                    : "border-slate-200 bg-white text-slate-500 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-400"
                }`}
              >
                {isSelected && (
                  <Check
                    className="h-3.5 w-3.5 text-brand-500"
                    aria-hidden="true"
                  />
                )}
                {child.firstName} · {child.poolCode.split("_")[0]}
              </button>
            );
          })}
          <button
            type="button"
            className="flex min-h-[44px] shrink-0 items-center justify-center rounded-lg border border-dashed border-slate-300 px-4 dark:border-slate-600"
            aria-label="Inscrire un autre enfant"
          >
            <Plus className="h-4 w-4 text-slate-400" aria-hidden="true" />
          </button>
        </div>
      </div>

      <div className="space-y-2">
        {selected.map((child) => {
          const pool = MOCK_POOLS.find((p) => p.code === child.poolCode);
          return (
            <div
              key={child.id}
              className="flex items-center justify-between rounded-xl border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-800/50"
            >
              <div>
                <p className="text-sm font-medium text-slate-900 dark:text-white">
                  {child.firstName} {child.lastName}
                </p>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  {pool?.name}
                </p>
              </div>
              <span
                className={
                  child.registrationStatus === "ASSIGNED"
                    ? "rounded-md bg-brand-50 px-2.5 py-1 text-xs font-medium text-brand-700 dark:bg-brand-900/40 dark:text-brand-300"
                    : "rounded-md bg-amber-50 px-2.5 py-1 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
                }
              >
                {child.registrationStatus === "ASSIGNED"
                  ? "Assigné"
                  : "En attente"}
              </span>
            </div>
          );
        })}
        {selected.length === 0 && (
          <p className="text-sm text-slate-400">
            Sélectionnez au moins un enfant pour voir ses informations.
          </p>
        )}
        {selected.some((c) => c.registrationStatus === "PENDING") && (
          <p className="text-xs text-slate-400">
            "En attente" : le bassin est confirmé, l'équipe reste à assigner par
            le club.
          </p>
        )}
      </div>
    </div>
  );
};
