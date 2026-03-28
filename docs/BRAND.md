# Brand Guidelines

## Logo

The Rx logo is a pure typographic wordmark — the letters "Rx" with no icons, symbols, or decorative elements.

### Design Decisions

| Aspect | Choice | Rationale |
|---|---|---|
| Style | Wordmark (text only) | Simplicity; must be recognizable at favicon size (16px) |
| Meaning | Pure brand name | No specific meaning imposed (e.g., prescription, CrossFit "as prescribed"); open to interpretation |
| Typeface | DM Sans Black (900) | Rounded geometric sans-serif; legible at small sizes; warm yet bold |
| Letter spacing | Tight (-4 approximate) | Compact, cohesive unit; the two letters read as one mark |
| Casing | Rx (capital R + lowercase x) | Matches the app name; natural reading form |
| Color | Monochrome | Consistent with the UI color system (pure grayscale, no hue) |
| Decoration | None | No underlines, enclosures, or letterform modifications |

### Usage

- **Favicon**: Dark rounded-square background (`#1a1a1a`) with light text (`#fafafa`). File: `web/app/icon.svg`.
- **In-app (sidebar)**: Plain text "Rx" using system font (`text-xl font-bold tracking-tight`). No SVG logo used in the UI itself.

### Color Variants

| Context | Background | Text |
|---|---|---|
| Favicon (light mode) | `#1a1a1a` | `#fafafa` |
| Favicon (dark mode) | `#1a1a1a` | `#fafafa` |
| On light background | — | `#1a1a1a` |
| On dark background | — | `#fafafa` |

### File Locations

| File | Purpose |
|---|---|
| `web/app/icon.svg` | Favicon (served automatically by Next.js App Router) |

### Future Considerations

- **Font outlining**: The current favicon SVG references Google Fonts externally. For offline resilience, convert text to `<path>` outlines.
- **Dark mode favicon**: Add `prefers-color-scheme` media query to serve a light-background variant in dark browser UI.
- **OGP / social images**: Generate from the same wordmark spec for consistency.
