# monitoring/

Source of truth for blog-specific Grafana dashboards. Each `dashboards/*.json` is
provisioned to the Studio's Grafana through its HTTP API by the `provision-grafana`
job in `.github/workflows/studio-workflow.yml` on every `develop` push that
touches a file here.

The repo doesn't own Grafana itself (that's `studio-configs`) — only the
dashboards that reference this repo's metric and log label names.

## How a dashboard gets to Grafana

1. Commit a JSON dashboard to `monitoring/dashboards/`. `uid` must be stable
   (Grafana uses it for upsert).
2. Push to `develop`. The CI job uploads it as the contents of a Grafana API
   payload (`POST /api/dashboards/db`) using `secrets.GRAFANA_TOKEN`.
3. The dashboard appears in Grafana within seconds.

UI edits in Grafana will be overwritten on the next deploy — the JSON file is
the source of truth. If you want to iterate on a panel visually, do it in
Grafana, export the JSON (`Share → Export → Save JSON`), and commit it back.

## Secret

`GRAFANA_TOKEN` — Grafana Service Account token, role Editor (Admin if a
dashboard ever needs to write Library Panels). Created in Grafana under
*Administration → Service accounts → blog-in-golang-ci → Add token*.

## Layout

```
monitoring/
├── README.md
└── dashboards/
    └── blog.json     # uid blog-http-logs — HTTP metrics + Loki logs
```
