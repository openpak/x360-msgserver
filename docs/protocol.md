# Protocol: x360-msgserver

Live service 2 (`msgserver`), XAM 17559.

XStringVerify request: `u32 title, u32 flags, u16 locale length (<= 16), u16 count (<= 10), locale, count x
{u16 length, UTF-8}`. Reply: `u16 count, count x u32 HRESULT`, all S_OK. The dashboard's Edit Profile vets Motto
and Bio with it; a failed request shows "Can't change your motto". 

Title names (0x800, `/msgserver/msgserver.ashx`): request `u16 1, u16 key, u16 n (<= 100), n x u32 title id`; reply
`u16 n, n x {u16 length, UTF-8 name}`, length 0 for an unknown title. Names come from the JSON file in
`X360_MSGSERVER_TITLES` (`title id hex -> name`, seed `deploy/titles.json`; empty map when unset).
String lookup (0x802, `/msgserver/getstring2.ashx`): request `u32 a0, u16 n (<= 100), u16 locale length (<= 16),
locale, n x u32 string id`; reply as 0x800, every string empty (no id table is known). All little-endian.

## Wharf

Behaviour is carried over unchanged from x360-xuacs v0.9.7, where the dashboard Edit Profile step was verified.
