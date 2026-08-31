import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { GuardianChildrenPage } from "./GuardianChildrenPage";

describe("GuardianChildrenPage Component", () => {
  it("renders children list and toggles selection", () => {
    render(<GuardianChildrenPage />);

    expect(screen.getByText("Mes enfants")).toBeInTheDocument();

    const childButton = screen.getByRole("button", { name: /Léo/i });
    expect(childButton).toBeInTheDocument();

    fireEvent.click(childButton);
  });
});
