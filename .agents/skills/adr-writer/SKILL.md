# Skill: ADR Writer (MADR 3.0.0 Standard)

## Role & Description

Expert Architecture Decision Record (ADR) author. Your mission is to analyze, research, and document architectural decisions using the official **MADR 3.0.0 Full Template** in strict English, while ensuring maximum context alignment and adherence to modern industry best practices.

## Triggers

Activate this skill whenever requested to document an architectural choice (e.g., "Create an ADR for...", "Rédige une ADR sur...", "Document this decision...").

## Pre-Execution & Reasoning Protocol

1. **Deep Context Acquisition:**

* Thoroughly inspect the existing repository structure, domain boundary rules (Clean Architecture / Hexagonal), and related ADRs in `docs/adr/`.
* Ensure a complete understanding of the problem space, operational constraints, and technical debt being addressed.

1. **Context-Optimized Reflection:**

* Evaluate how the decision impacts long-term maintainability, developer experience (DX), performance, and domain isolation of Project PULSE.
* Avoid generic arguments; ground every rationale in the project's specific context (e.g., Go modular monolith, PostgreSQL schema isolation per ADR-003, Huma v2 API layer).

1. **Modern Best Practices & Ecosystem Benchmarking:**

* Identify and compare state-of-the-art options, current industry standards, and community best practices for the problem at hand.
* Contrast the chosen solution against viable modern alternatives with realistic, objective trade-offs.

1. **Linting & Formatting Rigor:**

* Strictly follow standard Markdown rules to pass repository linters (`markdownlint` / `markdownlint-cli2`).
* No bare URLs (always use named link syntax `[Title](URL)`).
* No trailing spaces, proper list indentation (2 spaces), and consistent header hierarchy.

## Execution Rules

* **Target Directory:** `docs/refs/adr/`
* **File Naming:** `ADR-00X-short-kebab-case-title.md` (ADR- + 3-digit padded prefix, e.g., `ADR-007-adoption-of-huma-v2-for-openapi-rest-api.md`).
* **Language:** Write the ADR exclusively in **English**.
* **Template Compliance:** Read, parse, and populate every section defined in the template asset without skipping or altering required headers.

## Templates & Assets

When rendering the final ADR content, strictly load and enforce the structure, layout, and section rules defined in:

* [MADR 3.0.0 Full Template](./assets/madr-full-template.md)

### Template Processing Protocol

1. Read the full layout, headers, and section guidelines from `./assets/madr-full-template.md`.
2. Populate all sections using the domain knowledge and technical context acquired during the research phase.
3. Ensure every section title and structural constraint from `./assets/madr-full-template.md` is strictly respected in the generated output file.

## Quality Assurance Checklist

* [ ] Structure strictly matches `./assets/madr-full-template.md`.
* [ ] Context is fully tailored to the current codebase and project goals.
* [ ] Options evaluated reflect current modern engineering best practices.
* [ ] Trade-offs and negative consequences are explicitly detailed.
* [ ] Markdown format passes repository linter checks with zero warnings.
