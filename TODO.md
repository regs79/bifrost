# Bifrost TODO

This file tracks outstanding work to complete Bifrost v1.0 and prepare for distribution.

---

## Core Functionality

- [ ] Load config file from `$HOME/.config/bifrost/bifrost.toml`
- [ ] Use hardcoded defaults if config is missing
- [ ] Auto-create config directory if needed
- [ ] (Optional) Save a default `bifrost.toml` if user chooses
- [ ] Improve error handling on browser launch failures
- [ ] Validate URL input at launch
- [ ] Handle non-http links gracefully (e.g., `mailto:` `tel:`) (optional v2)

---

## Visual / UX Polish

- [ ] Tune canvas background transparency
- [ ] Center window nicely on screen at launch
- [ ] Adjust highlight opacity and color (Nord palette)
- [ ] (Optional) Add smooth highlight animation on hover/selection
- [ ] Support escape (`Esc`), quit (`q`), and number hotkeys for browsers

---

## Configuration / Theming

- [ ] Read browser list from config
- [ ] Read colors (background, text, highlight) from config
- [ ] Support hex codes with optional alpha (8 digits)
- [ ] (Optional) Add simple light/dark mode toggle (future)

---

## Packaging / Distribution

- [ ] Create a public Git repository (GitHub, GitLab)
- [ ] Add README.md (install instructions, screenshots)
- [ ] Add LICENSE file (MIT or Apache 2.0 recommended)
- [ ] Tag releases (`v1.0.0`, etc.) on GitHub
- [ ] Create AUR package (`PKGBUILD` for Arch Linux)
- [ ] (Optional) Create `.desktop` launcher for easy default browser setup

---

## Future Nice-to-Haves (Post v1.0)

- [ ] Real-time installed browser detection (instead of config-only)
- [ ] Per-domain browser rules (open specific domains with specific browsers)
- [ ] Drag window by clicking anywhere (frameless UX)
- [ ] Dynamic config reload (no restart needed)
- [ ] Visual hover effects (scale icons slightly, glow, etc.)
- [ ] Allow custom fonts (optional theme settings)

---
