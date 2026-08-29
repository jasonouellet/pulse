import { AppLayout } from "./components/layout/AppLayout";
import { SessionProvider } from "./context/SessionProvider";
import { HomePage } from "./pages/HomePage";

export function App() {
  return (
    <SessionProvider>
      <AppLayout>
        <div className="space-y-6">
          <header>
            <h1 className="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white sm:text-4xl">
              Tournament Dashboard
            </h1>
            <p className="mt-2 text-slate-600 dark:text-slate-400">
              Manage multi-sport youth pools, rosters, and dynamic match
              scheduling.
            </p>
          </header>

          <HomePage />
        </div>
      </AppLayout>
    </SessionProvider>
  );
}

export default App;
