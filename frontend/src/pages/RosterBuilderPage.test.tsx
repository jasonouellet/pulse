import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { describe, it, expect } from "vitest";
import { RosterBuilderPage } from "./RosterBuilderPage";

describe("RosterBuilderPage", () => {
  it("renders pool players and moves player between lists", () => {
    render(
      <MemoryRouter>
        <RosterBuilderPage />
      </MemoryRouter>,
    );

    // Vérifie le titre principal
    expect(screen.getByText("Formation des rosters")).toBeInTheDocument();

    // Trouve le bouton d'ajout du premier joueur disponible
    const addButtons = screen.getAllByRole("button", {
      name: /ajouter/i,
    });
    expect(addButtons.length).toBeGreaterThan(0);

    // Ajoute un joueur au roster
    fireEvent.click(addButtons[0]);

    // Vérifie que le bouton de retrait apparaît
    const removeButtons = screen.getAllByRole("button", {
      name: /retirer/i,
    });
    expect(removeButtons.length).toBeGreaterThan(0);
  });
});
