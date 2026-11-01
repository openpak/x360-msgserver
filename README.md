# Xbox 360 Message Server

OpenPak's Xbox 360 message server (Live service 2, `msgserver`).

| | |
| --- | --- |
| Network | Xbox 360 |
| Role | String verification for profile text (`XStringVerify`) |
| Depends on | [`x360-common`](https://github.com/openpak/x360-common) (vendored) |
| Image | `ghcr.io/openpak/x360-msgserver` |
| Used by | XAM through the security gateway |
| Licence | AGPL-3.0-only |

## What it serves

| Path | XRL | What |
| --- | --- | --- |
| `POST /msgserver/vetstring2.ashx` | 0x801 | Vet strings: every string is accepted |
| `GET /healthz` | | Health probe |

Identity headers (`X-X360-Account`, `X-X360-Users`) are trusted only with `X360_GATEWAY_KEY`. Split out of
[`x360-xuacs`](https://github.com/openpak/x360-xuacs) on 2026-11-01, one repository per Live service.

## Run

```sh
GOWORK=off go run -mod=vendor ./cmd/x360-msgserver
```


## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
