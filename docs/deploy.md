# Deploy

Primary host, quadlet `deploy/openpak-x360-msgserver.container` (IP pinned in openpak.network), Traefik routes from
`deploy/traefik.yml` merged into the shared dynamic file. Service 2 must be in the KDC's D7 list.

Title names: copy `deploy/titles.json` to `~/openpak/deploy/x360-msgserver-titles.json` on the primary (the quadlet mounts it read-only).
