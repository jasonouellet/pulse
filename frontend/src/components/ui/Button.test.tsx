import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { axe } from "vitest-axe";
import { Button } from "./Button";

describe("Button Component", () => {
  it("renders correctly with children", () => {
    render(<Button>Create Roster</Button>);
    // Remplacer toBeInDOM() par toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: /create roster/i }),
    ).toBeInTheDocument();
  });

  it("triggers onClick handler when clicked", () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Submit</Button>);

    fireEvent.click(screen.getByRole("button", { name: /submit/i }));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("is accessible according to WCAG rules (axe)", async () => {
    const { container } = render(<Button>Accessible Action</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});
