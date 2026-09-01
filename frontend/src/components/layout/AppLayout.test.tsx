import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { describe, it, expect } from "vitest";
import { AppLayout } from "./AppLayout";
import { SessionProvider } from "../../context/SessionProvider";
import "../../i18n/config";

describe("AppLayout Component", () => {
  it("renders brand title and toggles theme mode", () => {
    render(
      <SessionProvider>
        <MemoryRouter>
          <AppLayout />
        </MemoryRouter>
      </SessionProvider>,
    );

    expect(screen.getByText("PULSE")).toBeInTheDocument();

    const themeBtn = screen.getByRole("button", {
      name: /switch to dark mode|switch to light mode/i,
    });
    expect(themeBtn).toBeInTheDocument();

    fireEvent.click(themeBtn);
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });
});
