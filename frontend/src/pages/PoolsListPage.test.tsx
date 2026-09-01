import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { PoolsListPage } from "./PoolsListPage";

describe("PoolsListPage Component", () => {
  it("renders pools table and title", () => {
    render(<PoolsListPage />);

    expect(
      screen.getByRole("heading", { name: "Bassins" }),
    ).toBeInTheDocument();
    expect(screen.getByText("U10_D1")).toBeInTheDocument();
  });
});
