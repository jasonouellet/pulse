import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { RegistrationPage } from "./RegistrationPage";

describe("RegistrationPage", () => {
  it("navigates to new child form when clicking button", () => {
    render(<RegistrationPage />);

    expect(
      screen.getByText("Un de vos enfants déjà inscrit ?"),
    ).toBeInTheDocument();

    const newChildBtn = screen.getByRole("button", {
      name: /inscrire un nouvel enfant/i,
    });
    fireEvent.click(newChildBtn);

    expect(screen.getByPlaceholderText("Léo")).toBeInTheDocument();
  });
});
