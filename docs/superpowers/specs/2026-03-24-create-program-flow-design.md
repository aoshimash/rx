# Create Program Flow Design

**Date:** 2026-03-24
**Status:** Approved

## Overview

When the user clicks "Create Program" on `/programs`, present a 2-option choice dialog instead of jumping directly to a form. The two creation methods are:

1. **Import** — import a Program from a JSON file or pasted JSON (shared by another user or exported from this app)
2. **From Scratch** — the current behavior: enter name, notes, and session names manually

---

## Entry Point

The "Create Program" button (both the empty-state button and the outline button at the bottom of the list) opens a **Choice Dialog** with two option cards.

```
┌─────────────────────────────────────────┐
│  Create Program                          │
│  How would you like to create it?        │
│                                          │
│  ┌────────────────┐  ┌────────────────┐  │
│  │  Import        │  │  Scratch       │  │
│  │                │  │                │  │
│  │ Import from    │  │ Create from    │  │
│  │ JSON           │  │ scratch        │  │
│  └────────────────┘  └────────────────┘  │
│                                           │
│                        [Cancel]           │
└─────────────────────────────────────────┘
```

Selecting an option **transitions within the same dialog** (no new dialog opens). A back button returns to the choice screen.

---

## Flow 1: Import

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

## Flow 2: From Scratch

The existing form behavior (name + notes + session names) is preserved as-is. The form logic is extracted into `ScratchStep.tsx` as a pure refactor with no behavior change. Existing behavior should be verified with a manual smoke test after the refactor.

---

## Component structure (proposed)

```
components/programs/
  CreateProgramDialog.tsx        — root dialog, manages which step is shown
  create-program/
    ChoiceStep.tsx               — 2-option card selection
    ImportStep.tsx               — drag/drop + paste + Zod validation
    ScratchStep.tsx              — current form, extracted into a component
```

---

## Export

An **Export** option is added to the action menu on the Program detail page. Two actions are provided:

- **Copy to clipboard** — copies the JSON and shows a brief toast ("Copied to clipboard")
- **Download .json** — triggers a file download named `{program-name}.json`

The exported JSON uses the format described in [Flow 1: Import](#flow-1-import).

---

## API changes required

- **Program name uniqueness**: Add a uniqueness constraint on `programs.name` (DB migration + domain validation + `409 Conflict` response). Program names serve as natural identifiers for API/CLI/AI consumers.

---

## Out of scope

- Saving gym profile (default increments per exercise) — future enhancement
- Sharing via link or short code
