# Changelog

## v0.2.1 (2026-11-03)

- Every request logs one stderr line, `http <method> <path> -> <status> <duration>ms` (x360-common v0.5.1 `httplog`): no query, headers, cookies, tokens or bodies; `/healthz` skipped.

## v0.2.0

- 0x800 title names (`X360_MSGSERVER_TITLES` JSON, seed `deploy/titles.json`) and 0x802 string lookup (empty replies).

## v0.1.0

- Live service 2 split out of x360-xuacs v0.9.7, behaviour unchanged.
