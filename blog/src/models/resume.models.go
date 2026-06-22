package models

import (
	"encoding/json"
	"os"
)

// Resume is the single source of truth for both the public "about" page
// (rendered by handlers.AboutPage from blog/templates/about.html) and the
// distributed PDF résumé (rendered by blog/data/resume.typ via Typst). Edit
// blog/data/resume.json and BOTH outputs update — see blog/data/README.md.
type Resume struct {
	Basics    ResumeBasics      `json:"basics"`
	Education []ResumeEducation `json:"education"`
	Skills    []ResumeSkill     `json:"skills"`
	Work      []ResumeWorkEntry `json:"work"`
}

// ResumeBasics holds name + contact details. Used by the PDF header; the web
// page sources name/social links from header.html instead.
type ResumeBasics struct {
	Name     string `json:"name"`
	Label    string `json:"label"`
	Email    string `json:"email"`
	Blog     string `json:"blog"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

// ResumeEducation is one school entry. AttendedShort ("June '15 to March '18")
// feeds the web layout; Dates ("June 2015 - March 2018") feeds the PDF.
type ResumeEducation struct {
	Degree        string `json:"degree"`
	School        string `json:"school"`
	Location      string `json:"location"`
	AttendedShort string `json:"attendedShort"`
	Dates         string `json:"dates"`
}

// ResumeSkill is one "Category: items" row. Rendered on the PDF only.
type ResumeSkill struct {
	Category string `json:"category"`
	Items    string `json:"items"`
}

// ResumeMeta is one right-aligned "Label: Value" row under a work entry on the
// web page (Department, Main Languages, Operating Costs Reduced, …).
type ResumeMeta struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ResumeWorkEntry is one job. The web page uses the detailed fields (Meta +
// Bullets + SubEntries); the PDF uses Title/Company/Dates + PDF.Bullets so the
// one-pager stays tight while the web bio stays rich.
type ResumeWorkEntry struct {
	Company        string            `json:"company"`
	Title          string            `json:"title"`
	Dates          string            `json:"dates"`
	BoldTitle      bool              `json:"boldTitle"`
	HeaderOnly     bool              `json:"headerOnly"`
	Intro          string            `json:"intro"`
	Meta           []ResumeMeta      `json:"meta"`
	BulletsHeading string            `json:"bulletsHeading"`
	BulletsFloat   bool              `json:"bulletsFloat"`
	Bullets        []string          `json:"bullets"`
	SubEntries     []ResumeWorkEntry `json:"subEntries"`
	PDF            *ResumePDF        `json:"pdf"`
}

// ResumePDF carries the tight, one-page variant of a work entry's bullets.
// Include=false drops the entry from the PDF entirely (e.g. the internship).
type ResumePDF struct {
	Include bool     `json:"include"`
	Bullets []string `json:"bullets"`
}

// ResumeDataPath is the default location of the résumé data file, relative to
// the process working directory — the same convention LoadHTMLGlob("templates/*")
// and ./public use (blog/ in dev and tests, /app in the container).
const ResumeDataPath = "data/resume.json"

// LoadResume reads and parses the résumé data file from ResumeDataPath.
func LoadResume() (*Resume, error) {
	return LoadResumeFrom(ResumeDataPath)
}

// LoadResumeFrom reads and parses a résumé data file from an explicit path.
func LoadResumeFrom(path string) (*Resume, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Resume
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
