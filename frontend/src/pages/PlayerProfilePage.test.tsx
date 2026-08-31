import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, it, expect } from "vitest";
import { PlayerProfilePage } from "./PlayerProfilePage";
import { SessionProvider } from "../context/SessionProvider";

describe("PlayerProfilePage Component", () => {
  it("renders player profile information", () => {
    render(
      <SessionProvider>
        <MemoryRouter initialEntries={["/players/player-leo"]}>
          <Routes>
            <Route path="/players/:playerId" element={<PlayerProfilePage />} />
          </Routes>
        </MemoryRouter>
      </SessionProvider>,
    );

    expect(screen.getByText("Léo Tremblay")).toBeInTheDocument();
    expect(screen.getByText("Positions préférées")).toBeInTheDocument();
  });
});
