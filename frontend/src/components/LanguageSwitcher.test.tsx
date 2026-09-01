import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { LanguageSwitcher } from "./LanguageSwitcher";
import i18n from "../i18n/config";

describe("LanguageSwitcher Component", () => {
  it("renders correctly and calls changeLanguage on click", () => {
    const spy = vi.spyOn(i18n, "changeLanguage");

    render(<LanguageSwitcher />);

    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();

    fireEvent.click(button);
    expect(spy).toHaveBeenCalled();
  });
});
