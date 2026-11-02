# Protocol: x360-msgserver

Live service 2 (`msgserver`), XAM 17559.

XStringVerify request: `u32 title, u32 flags, u16 locale length (<= 16), u16 count (<= 10), locale, count x
{u16 length, UTF-8}`. Reply: `u16 count, count x u32 HRESULT`, all S_OK. The dashboard's Edit Profile vets Motto
and Bio with it; a failed request shows "Can't change your motto". Requests 0x800 and 0x802 are not served.

## Wharf

Behaviour is carried over unchanged from x360-xuacs v0.9.7, where the dashboard Edit Profile step was verified.
