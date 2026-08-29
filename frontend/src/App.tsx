import { Routes, Route } from "react-router";
import { AppLayout } from "./components/layout/AppLayout";
import { SessionProvider } from "./context/SessionProvider";
import { HomePage } from "./pages/HomePage";
import { PoolsListPage } from "./pages/PoolsListPage";
import { RosterBuilderPage } from "./pages/RosterBuilderPage";
import { PlayerProfilePage } from "./pages/PlayerProfilePage";
import { GuardianChildrenPage } from "./pages/GuardianChildrenPage";
import { RegistrationPage } from "./pages/RegistrationPage";

export function App() {
  return (
    <SessionProvider>
      <AppLayout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/pools" element={<PoolsListPage />} />
          <Route path="/rosters" element={<RosterBuilderPage />} />
          <Route path="/players/:playerId" element={<PlayerProfilePage />} />
          <Route path="/children" element={<GuardianChildrenPage />} />
          <Route path="/register" element={<RegistrationPage />} />
        </Routes>
      </AppLayout>
    </SessionProvider>
  );
}

export default App;
