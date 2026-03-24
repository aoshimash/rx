# Create Program Flow Design

**Date:** 2026-03-24
**Status:** Approved

## Overview

When the user clicks "Create Program" on `/programs`, present a 3-option choice dialog instead of jumping directly to a form. The three creation methods are:

1. **From Template** — generate a Program from an existing ProgramTemplate with concrete loads calculated from user-provided 1RM values
2. **Import** — import a Program from a JSON file or pasted JSON (shared by another user or exported from this app)
3. **From Scratch** — the current behavior: enter name, notes, and session names manually

---

## Entry Point

The "Create Program" button (both the empty-state button and the outline button at the bottom of the list) opens a **Choice Dialog** with three option cards.

```
┌─────────────────────────────────────────┐
│  Create Program                          │
│  How would you like to create it?        │
│                                          │
│  ┌────────────────┐  ┌────────────────┐  │
│  │  Template      │  │  Import        │  │
│  │                │  │                │  │
│  │ Generate from  │  │ Import from    │  │
│  │ a template     │  │ JSON           │  │
│  └────────────────┘  └────────────────┘  │
│  ┌────────────────┐                       │
│  │  Scratch       │                       │
│  │                │                       │
│  │ Create from    │                       │
│  └────────────────┘                       │
│                                           │
│                        [Cancel]           │
└─────────────────────────────────────────┘
```

Selecting an option **transitions within the same dialog** (no new dialog opens). A back button returns to the choice screen.

---

## Flow 1: From Template

### Step 1 — Select a ProgramTemplate

Display a searchable list of ProgramTemplates. The user selects one and clicks Next.

### Step 2 — Enter weight values and load increments

The dialog extracts all unique exercise names from the selected template's entries. The API behavior differs by intensity type:

- **`percent_1rm` entries**: `load_kg = percent_1rm × target_weight` → user enters their **1RM**
- **`rpe`-only entries** (no `percent_1rm`): `load_kg = target_weight` (direct copy) → user enters the **target weight** directly

The UI labels the input field accordingly per exercise. If a template has no entries, the Generate button is disabled with a message ("This template has no exercises").

```
← Back  |  Create from Template

  Program Name
  [ 5/3/1 BBB                  ]

  Exercise settings

  ┌──────────────────────────────────────────────┐
  │ Squat             (% 1RM)                    │
  │  1RM    [ 140 ] kg    Increment [ 2.5 ] kg   │
  ├──────────────────────────────────────────────┤
  │ Bench Press       (% 1RM)                    │
  │  1RM    [ 100 ] kg    Increment [ 2.5 ] kg   │
  ├──────────────────────────────────────────────┤
  │ Pause Squat       (RPE only)                 │
  │  Weight [ 120 ] kg    Increment [ 2.5 ] kg   │
  └──────────────────────────────────────────────┘

  * 1RM is required for all exercises. Increment defaults to 2.5 kg.
  * Negative or zero values are not accepted.

                [Cancel]  [Generate Program]
```

Validation: all weight fields must be positive numbers. The Generate button is disabled until all fields are filled.

### API call

```
POST /program-templates/{id}/generate
{
  "name": "5/3/1 BBB",
  "target_weights": { "Squat": 140, "Bench Press": 100, "Pause Squat": 120 },
  "load_increments": { "Squat": 2.5, "Bench Press": 2.5, "Pause Squat": 2.5 }
}
```

This endpoint already exists and creates a Program with `ProgramSessionEntry.load_kg` values populated. The frontend API client (`web/lib/api/programTemplates.ts`) already has a `generate()` function.

### Template list behavior

- Archived templates (`archived_at` is set) are excluded from the list
- The list is fetched with cursor-based pagination; all pages are loaded client-side for simplicity (templates are not expected to be numerous enough to require server-side search)
- A client-side search filter is applied to the fetched list

### Result

A Program is created with concrete `load_kg` values for every entry, linked to the source template via `program_template_id`.

---

## Flow 2: Import

### Import UI

Three input methods, all feeding into the same Zod validation pipeline:

1. **Drag and drop** — drop a `.json` file onto the drop zone
2. **Browse** — click to open a file picker
3. **Paste** — paste JSON text directly into the text area

```
← Back  |  Import Program

  ┌─────────────────────────────────────┐
  │                                     │
  │   📂 Drop JSON file here            │
  │   or paste JSON below               │
  │                                     │
  └─────────────────────────────────────┘

  ┌─────────────────────────────────────┐
  │ { "name": "5/3/1 BBB", ...          │  ← populated after drop/paste
  └─────────────────────────────────────┘
  [📂 Browse file]

  ❌ Invalid format: "sessions" is required   ← Zod error shown inline

                [Cancel]  [Import]   ← enabled only when validation passes
```

### JSON format

The export/import format is a thin envelope around `ProgramCreate`:

```json
{
  "rx_version": "1",
  "name": "5/3/1 BBB",
  "notes": "optional",
  "sessions": [
    {
      "session_name": "Week1 Day1",
      "order": 0,
      "entries": [
        {
          "exercise_name": "Squat",
          "sets": 5,
          "reps": 5,
          "load_kg": 100,
          "order": 0
        }
      ]
    }
  ]
}
```

`rx_version` is hardcoded to `"1"` on export by the frontend. The Zod schema validates the full structure at import time and provides clear error messages for malformed input. Before calling the API, `rx_version` is stripped from the payload (it is not part of `ProgramCreate`).

### Import API call

After validation passes, the frontend calls:

```
POST /programs
{
  "name": "5/3/1 BBB",
  "notes": "optional",
  "sessions": [ ... ]
}
```

If a program with the same name already exists, the API returns a `409 Conflict` error. The UI surfaces this as an inline error on the name field ("A program with this name already exists").


---

## Flow 3: From Scratch

The existing form behavior (name + notes + session names) is preserved as-is. The form logic is extracted into `ScratchStep.tsx` as a pure refactor with no behavior change. Existing behavior should be verified with a manual smoke test after the refactor.

---

## Component structure (proposed)

```
components/programs/
  CreateProgramDialog.tsx        — root dialog, manages which step is shown
  create-program/
    ChoiceStep.tsx               — 3-option card selection
    TemplateSelectStep.tsx       — searchable template list
    TemplateConfigStep.tsx       — 1RM + increment inputs per exercise
    ImportStep.tsx               — drag/drop + paste + Zod validation
    ScratchStep.tsx              — current form, extracted into a component
```

---

## Export

An **Export** option is added to the action menu on the Program detail page. Two actions are provided:

- **Copy to clipboard** — copies the JSON and shows a brief toast ("Copied to clipboard")
- **Download .json** — triggers a file download named `{program-name}.json`

The exported JSON uses the format described in [Flow 2: Import](#flow-2-import).

---

## API changes required

- **Program name uniqueness**: Add a uniqueness constraint on `programs.name` (DB migration + domain validation + `409 Conflict` response). Program names serve as natural identifiers for API/CLI/AI consumers.

---

## Out of scope

- Saving gym profile (default increments per exercise) — future enhancement
- Importing ProgramTemplates (only Programs are in scope here)
- Sharing via link or short code
