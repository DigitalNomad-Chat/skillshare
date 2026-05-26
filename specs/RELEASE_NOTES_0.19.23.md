# skillshare v0.19.23 Release Notes

Release date: 2026-05-26

## TL;DR

1. **Nested GitHub-installed skills stop reappearing as updateable after update** — skills installed under a subdirectory now refresh their stored update metadata correctly after `skillshare update` or dashboard updates.
2. **Dashboard Update checks are remembered between visits** — the Update page now restores the last completed check status and shows when it was checked, instead of resetting everything to `Unchecked` on entry.

This is a patch release — bug fixes only, no new commands or flags.

---

## Nested Skills No Longer Stay Marked as Updateable

For GitHub-installed skills stored under a subdirectory, a successful update could still leave the item showing **Update available** afterward. v0.19.23 refreshes the stored metadata after the updated skill is moved into place, so the CLI and dashboard agree that the skill has been updated.

```bash
skillshare update tools/agent-browser
```

After the update finishes, the next check uses the refreshed metadata instead of comparing against the old install record.

## Update Page Remembers the Last Check

The dashboard Update page now keeps the last completed check result in browser storage. When you leave and come back, rows keep their latest status — such as **Update available** or **Up to date** — and show the previous check time.

This keeps the page useful between visits while still letting you refresh with **Check All** or **Check Selected** whenever you want a new result.
