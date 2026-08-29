import React from "react";
import {
  Users,
  Calendar,
  CreditCard,
  ClipboardList,
  Trophy,
  Settings,
  ChevronRight,
} from "lucide-react";
import { ROLE_LABELS, useSession, type UserRole } from "../context/session";

interface QuickAccessItem {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}

// What each role sees in "Accès rapide" — extend as real screens land.
const QUICK_ACCESS_BY_ROLE: Record<UserRole, QuickAccessItem[]> = {
  GUARDIAN: [
    { label: "Mes enfants", icon: Users },
    { label: "Calendrier", icon: Calendar },
    { label: "Paiements", icon: CreditCard },
  ],
  PLAYER: [
    { label: "Mon horaire", icon: Calendar },
    { label: "Mes objectifs", icon: ClipboardList },
  ],
  COACH: [
    { label: "Mes groupes", icon: Users },
    { label: "Formation de rosters", icon: ClipboardList },
    { label: "Calendrier", icon: Calendar },
  ],
  CLUB_ADMIN: [
    { label: "Bassins", icon: Users },
    { label: "Tournois", icon: Trophy },
    { label: "Paramètres du club", icon: Settings },
  ],
  TECHNICAL_DIRECTOR: [
    { label: "Bassins", icon: Users },
    { label: "Tournois", icon: Trophy },
    { label: "Terrains", icon: Calendar },
  ],
};

// Mock — replace with GET /api/v1/{module}/events?upcoming=true once that
// endpoint exists.
const MOCK_UPCOMING_EVENTS = [
  { title: "Match Léo vs Les Faucons", when: "Dim 6 sept · 11h00" },
  { title: "Pratique U10-M", when: "Sam 5 sept · 9h00" },
];

export const HomePage: React.FC = () => {
  const { activeGrant } = useSession();
  const quickAccess = QUICK_ACCESS_BY_ROLE[activeGrant.role];

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-sm font-medium text-slate-500 dark:text-slate-400">
          Accès rapide · {ROLE_LABELS[activeGrant.role]} @{" "}
          {activeGrant.club.name}
        </h2>
        <div className="mt-3 grid grid-cols-2 gap-3 sm:grid-cols-3">
          {quickAccess.map((item) => {
            const Icon = item.icon;
            return (
              <button
                key={item.label}
                type="button"
                className="flex flex-col items-center gap-2 rounded-xl border border-slate-200 bg-white p-4 text-center shadow-sm transition-colors hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:border-slate-800 dark:bg-slate-800/50 dark:hover:bg-slate-800"
              >
                <Icon className="h-5 w-5 text-brand-500" aria-hidden="true" />
                <span className="text-xs font-medium text-slate-700 dark:text-slate-200">
                  {item.label}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      <div>
        <h2 className="text-sm font-medium text-slate-500 dark:text-slate-400">
          Événements à venir
        </h2>
        <div className="mt-3 overflow-hidden rounded-xl border border-slate-200 dark:border-slate-800">
          {MOCK_UPCOMING_EVENTS.map((event, i) => (
            <div
              key={event.title}
              className={`flex items-center justify-between px-4 py-3 ${
                i < MOCK_UPCOMING_EVENTS.length - 1
                  ? "border-b border-slate-200 dark:border-slate-800"
                  : ""
              }`}
            >
              <div>
                <p className="text-sm text-slate-900 dark:text-white">
                  {event.title}
                </p>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  {event.when}
                </p>
              </div>
              <ChevronRight
                className="h-4 w-4 text-slate-400"
                aria-hidden="true"
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
