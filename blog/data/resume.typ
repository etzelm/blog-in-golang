// resume.typ — renders the distributed one-page PDF résumé from resume.json,
// the SAME single source the website's about page reads. Compile with:
//   typst compile blog/data/resume.typ "Mitchell-Etzel's-Resume.pdf" \
//     --font-path blog/data/fonts --ignore-system-fonts
//
// Font: Liberation Sans, bundled in data/fonts/ (OFL — see data/fonts/LICENSE).
// --ignore-system-fonts makes the output byte-identical on every machine,
// independent of whatever fonts happen to be installed.

#let data = json("resume.json")
#let b = data.basics

#set document(title: b.name + " — Résumé", author: b.name)
#set page(paper: "us-letter", margin: (x: 0.5in, top: 0.5in, bottom: 0.45in))
#set text(
  font: "Liberation Sans",
  size: 10pt,
  fill: rgb("#1a1a1a"),
)
#set par(leading: 0.52em, justify: false)
#show link: set text(fill: rgb("#0b5cad"))

// Tight, clean bullet lists.
#set list(marker: text(fill: rgb("#0b5cad"))[•], indent: 2pt, body-indent: 6pt, spacing: 4.5pt)

// Pretty date ranges: "December 2024 - February 2026" -> en dash.
#let pretty(s) = s.replace(" - ", " – ")

// Section header with an underline rule.
#let section(title) = {
  v(7.5pt)
  text(size: 13pt, weight: "bold", fill: rgb("#0b5cad"))[#upper(title)]
  v(1.5pt)
  line(length: 100%, stroke: 0.6pt + rgb("#9aa0a6"))
  v(3.5pt)
}

// ── Header ───────────────────────────────────────────────────────────────
#text(size: 21pt, weight: "bold")[#b.name]
#v(2.5pt)
#text(size: 9pt)[
  Email: #link("mailto:" + b.email)[#b.email]
  #h(6pt) | #h(6pt) Professional Blog: #link("https://" + b.blog)[#b.blog]
  #h(6pt) | #h(6pt) LinkedIn: #link("https://" + b.linkedin)[#b.linkedin]
  #h(6pt) | #h(6pt) GitHub: #link("https://" + b.github)[#b.github]
]

// ── Experience ───────────────────────────────────────────────────────────
#section("Experience")
#for job in data.work {
  if job.pdf.include {
    block(breakable: false, width: 100%)[
      #grid(
        columns: (1fr, auto),
        align: (left, right),
        text(weight: "bold")[#job.title #sym.dash.en #job.company],
        text(size: 9.5pt, fill: rgb("#444444"))[#pretty(job.dates)],
      )
      // Keep slash-joined terms (CI/CD, critical/high) from breaking across
      // lines: append a zero-width word joiner after each "/".
      #list(..job.pdf.bullets.map(b => b.replace("/", "/\u{2060}")))
    ]
    v(5pt)
  }
}

// ── Education ─────────────────────────────────────────────────────────────
#section("Education")
#for ed in data.education {
  grid(
    columns: (1fr, auto),
    align: (left, right),
    [#text(weight: "bold")[#ed.degree] #h(8pt) #text(size: 9.5pt, fill: rgb("#444444"))[#ed.school]],
    text(size: 9.5pt, fill: rgb("#444444"))[#pretty(ed.dates)],
  )
  v(3pt)
}

// ── Computer Skills ───────────────────────────────────────────────────────
#section("Computer Skills")
#for sk in data.skills {
  block(spacing: 3.4pt)[#text(weight: "bold")[#sk.category:] #sk.items]
}
