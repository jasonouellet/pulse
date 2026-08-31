import React, { useState } from "react";
import { Lock, Check } from "lucide-react";
import { Button } from "../components/ui/Button";
import { MOCK_PLAYERS, MY_CHILDREN_IDS } from "../data/mock";

const CURRENT_SEASON_YEAR = 2026;

function suggestPoolLabel(birthDate: string): string {
  const birthYear = new Date(birthDate).getFullYear();
  const age = CURRENT_SEASON_YEAR - birthYear;
  return `U${age}`;
}

export const RegistrationPage: React.FC = () => {
  // Identity (name, DOB, medical info) is permanent — it lives on
  // core.player_profiles and shouldn't be re-entered every season.
  // Registration is annual — it's a new core.pool_registrations row each
  // year. So the first choice here is: an existing child (skip straight to
  // confirming the season's pool) or a brand-new one (full identity form).
  const [mode, setMode] = useState<"choose" | "existing" | "new">("choose");
  const [selectedChildId, setSelectedChildId] = useState<string | null>(null);

  const myChildren = MOCK_PLAYERS.filter((p) => MY_CHILDREN_IDS.includes(p.id));
  const selectedChild = myChildren.find((c) => c.id === selectedChildId);

  const [newBirthDate, setNewBirthDate] = useState("2018-03-14");

  if (mode === "choose") {
    return (
      <div className="mx-auto max-w-md space-y-4">
        <div>
          <h1 className="text-xl font-bold text-slate-900 dark:text-white">
            Inscrire un nouvel enfant
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            Saison {CURRENT_SEASON_YEAR} · Soccer
          </p>
        </div>

        <p className="text-sm font-medium text-slate-600 dark:text-slate-300">
          Un de vos enfants déjà inscrit ?
        </p>
        <div className="space-y-2">
          {myChildren.map((child) => (
            <button
              key={child.id}
              type="button"
              onClick={() => {
                setSelectedChildId(child.id);
                setMode("existing");
              }}
              className="flex min-h-[44px] w-full items-center justify-between rounded-lg border border-slate-200 bg-white px-4 py-3 text-left text-sm hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:border-slate-700 dark:bg-slate-800 dark:hover:bg-slate-700"
            >
              <span className="text-slate-900 dark:text-white">
                {child.firstName} {child.lastName}
              </span>
              <span className="text-xs text-slate-400">
                Réinscrire pour {CURRENT_SEASON_YEAR}
              </span>
            </button>
          ))}
        </div>

        <div className="flex items-center gap-3 text-xs text-slate-400">
          <div className="h-px flex-1 bg-slate-200 dark:bg-slate-700" />
          ou
          <div className="h-px flex-1 bg-slate-200 dark:bg-slate-700" />
        </div>

        <Button
          variant="outline"
          className="w-full"
          onClick={() => setMode("new")}
        >
          Inscrire un nouvel enfant
        </Button>
      </div>
    );
  }

  if (mode === "existing" && selectedChild) {
    return (
      <div className="mx-auto max-w-md space-y-4">
        <div>
          <h1 className="text-xl font-bold text-slate-900 dark:text-white">
            Réinscrire {selectedChild.firstName}
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            Saison {CURRENT_SEASON_YEAR} · Soccer
          </p>
        </div>

        <div className="flex items-center gap-2 rounded-lg bg-slate-50 p-3 dark:bg-slate-800">
          <Check className="h-4 w-4 text-brand-500" aria-hidden="true" />
          <p className="text-sm text-slate-600 dark:text-slate-300">
            Identité déjà connue — aucune information à ressaisir.
          </p>
        </div>

        <div className="flex items-start gap-2.5 rounded-lg bg-brand-50 p-3 dark:bg-brand-900/30">
          <Lock
            className="mt-0.5 h-4 w-4 shrink-0 text-brand-600 dark:text-brand-400"
            aria-hidden="true"
          />
          <div>
            <p className="text-sm font-medium text-brand-700 dark:text-brand-300">
              Bassin suggéré : {suggestPoolLabel(selectedChild.birthDate)}
            </p>
            <p className="text-xs text-brand-600 dark:text-brand-400">
              Calculé selon l'année de naissance · confirmé par le staff à
              l'assignation d'équipe
            </p>
          </div>
        </div>

        <div>
          <label
            htmlFor="existing-medical-notes"
            className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
          >
            Notes médicales — à jour ?
          </label>
          <input
            id="existing-medical-notes"
            type="text"
            placeholder="Allergies, conditions particulières..."
            className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
          />
        </div>

        <div className="border-t border-slate-200 pt-3 dark:border-slate-800">
          <div className="flex justify-between text-sm">
            <span className="text-slate-500 dark:text-slate-400">
              Frais d'inscription saison
            </span>
            <span className="text-slate-900 dark:text-white">85 $</span>
          </div>
        </div>

        <Button variant="primary" className="w-full">
          Confirmer la réinscription
        </Button>
        <button
          type="button"
          onClick={() => setMode("choose")}
          className="w-full text-center text-xs text-slate-400 hover:underline"
        >
          Retour
        </button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-md space-y-4">
      <div>
        <h1 className="text-xl font-bold text-slate-900 dark:text-white">
          Inscrire un nouvel enfant
        </h1>
        <p className="text-sm text-slate-500 dark:text-slate-400">
          Saison {CURRENT_SEASON_YEAR} · Soccer
        </p>
      </div>

      <div>
        <label
          htmlFor="new-first-name"
          className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
        >
          Prénom
        </label>
        <input
          id="new-first-name"
          type="text"
          placeholder="Léo"
          className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        />
      </div>

      <div>
        <label
          htmlFor="new-last-name"
          className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
        >
          Nom
        </label>
        <input
          id="new-last-name"
          type="text"
          placeholder="Tremblay"
          className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        />
      </div>

      <div>
        <label
          htmlFor="new-birth-date"
          className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
        >
          Date de naissance
        </label>
        <input
          id="new-birth-date"
          type="date"
          value={newBirthDate}
          onChange={(e) => setNewBirthDate(e.target.value)}
          className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        />
      </div>

      <div className="flex items-start gap-2.5 rounded-lg bg-brand-50 p-3 dark:bg-brand-900/30">
        <Lock
          className="mt-0.5 h-4 w-4 shrink-0 text-brand-600 dark:text-brand-400"
          aria-hidden="true"
        />
        <div>
          <p className="text-sm font-medium text-brand-700 dark:text-brand-300">
            Bassin suggéré : {suggestPoolLabel(newBirthDate)}
          </p>
          <p className="text-xs text-brand-600 dark:text-brand-400">
            Calculé selon l'année de naissance · confirmé par le staff à
            l'assignation d'équipe
          </p>
        </div>
      </div>

      <div>
        <label
          htmlFor="new-emergency-contact"
          className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
        >
          Contact d'urgence
        </label>
        <input
          id="new-emergency-contact"
          type="text"
          placeholder="Nom complet"
          className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        />
      </div>

      <div>
        <label
          htmlFor="new-medical-notes"
          className="mb-1 block text-xs text-slate-500 dark:text-slate-400"
        >
          Notes médicales (optionnel)
        </label>
        <input
          id="new-medical-notes"
          type="text"
          placeholder="Allergies, conditions particulières..."
          className="min-h-[44px] w-full rounded-lg border border-slate-300 px-3 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        />
      </div>

      <div className="border-t border-slate-200 pt-3 dark:border-slate-800">
        <div className="flex justify-between text-sm">
          <span className="text-slate-500 dark:text-slate-400">
            Frais d'inscription saison
          </span>
          <span className="text-slate-900 dark:text-white">85 $</span>
        </div>
      </div>

      <Button variant="primary" className="w-full">
        Confirmer l'inscription
      </Button>
      <button
        type="button"
        onClick={() => setMode("choose")}
        className="w-full text-center text-xs text-slate-400 hover:underline"
      >
        Retour
      </button>
    </div>
  );
};
