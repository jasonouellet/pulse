import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { HomePage } from "./HomePage";
import { SessionProvider } from "../context/SessionProvider";

describe("HomePage Component", () => {
  it("renders welcome title and quick access links", () => {
    render(
      <SessionProvider>
        <HomePage />
      </SessionProvider>,
    );

    expect(screen.getByText(/bienvenue/i)).toBeInTheDocument();
    expect(screen.getByText("Mes enfants")).toBeInTheDocument();
    expect(screen.getByText("Événements à venir")).toBeInTheDocument();
  });
});
