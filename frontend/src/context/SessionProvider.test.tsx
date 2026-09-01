import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { SessionProvider } from "./SessionProvider";
import { useSession } from "./session";

const TestComponent = () => {
  const { activeGrant, setActiveRole, setActiveClub } = useSession();
  return (
    <div>
      <p data-testid="active-role">{activeGrant.role}</p>
      <p data-testid="active-club">{activeGrant.club.id}</p>
      <button onClick={() => setActiveRole("COACH")}>Set Coach</button>
      <button onClick={() => setActiveClub("club-hockey-rimouski")}>
        Set Hockey
      </button>
    </div>
  );
};

describe("SessionProvider", () => {
  it("provides default session context and updates role and club", () => {
    render(
      <SessionProvider>
        <TestComponent />
      </SessionProvider>,
    );

    expect(screen.getByTestId("active-role").textContent).toBe("GUARDIAN");

    fireEvent.click(screen.getByText("Set Coach"));
    expect(screen.getByTestId("active-role").textContent).toBe("COACH");

    fireEvent.click(screen.getByText("Set Hockey"));
    expect(screen.getByTestId("active-club").textContent).toBe(
      "club-hockey-rimouski",
    );
  });

  it("throws error when useSession is used outside SessionProvider", () => {
    const InvalidComponent = () => {
      useSession();
      return null;
    };

    expect(() => render(<InvalidComponent />)).toThrow(
      "useSession must be used within a SessionProvider",
    );
  });
});
