# Résumé — single source of truth

`resume.json` is the **one place** to edit résumé content. Two consumers render
from it, so the website and the distributed PDF can never drift:

| Output | Rendered by | Where it shows up |
| --- | --- | --- |
| **Website** "about" page | Go `html/template` — [`templates/about.html`](../templates/about.html) via `handlers.AboutPage` (reads [`models.LoadResume`](../src/models/resume.models.go)) | `https://mitchelletzel.com/` |
| **Distributed PDF** | Typst — [`resume.typ`](resume.typ) | `https://files.mitchelletzel.com/Mitchell-Etzel's-Resume.pdf` |

## To update your résumé

1. Edit `resume.json`.
2. Open a PR. CI rebuilds the site (new about page) and, on the prod (`master`)
   path, recompiles the PDF and publishes it to `files.mitchelletzel.com`
   (S3 + CloudFront invalidation). No Google Docs, no manual upload.

## Data model (per work entry)

The website uses the rich fields; the PDF uses the tight ones — same entry:

- `title`, `company`, `dates` — used by both.
- `meta[]` (`{label, value}`), `bullets[]`, `bulletsHeading`, `bulletsFloat`,
  `intro`, `subEntries[]` — **website** layout (detailed).
- `pdf.include` + `pdf.bullets[]` — **PDF** one-pager. `include: false` drops
  the entry from the PDF (e.g. the early internship).

## Preview locally

```bash
# PDF (needs typst >= 0.15: `brew install typst`)
cd blog && typst compile data/resume.typ "Mitchell-Etzel's-Resume.pdf"

# Website
cd blog && go run .   # then open http://localhost:8080/
```

The PDF uses **Liberation Sans** (an Arial-metric-compatible OFL font). CI
installs `fonts-liberation`; a local macOS render falls back to Arial, which is
metric-identical, so line breaks and the one-page fit are the same either way.
