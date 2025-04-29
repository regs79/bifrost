# Changelog

All notable changes to this project will be documented in this file.

This project follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

- Preparing for future releases
- Ongoing migration to GTK4
- Planning Flathub and AUR packaging

---

## [0.1.0] - 2025-04-27

### Added
- Initial public alpha release of Bifrost
- Basic browser picker UI
- Horizontal layout of browser icons and names
- Keyboard shortcuts: Arrow keys, number keys (1-9), Enter, Escape, Q
- Mouse hover highlights browser entry
- Config file support (`bifrost.toml`) for:
  - Defining available browsers
  - Customizing theme colors (background, text, highlight)
- Nord color theme default
- High-DPI friendly icon scaling
- Basic error handling on browser launch
- Window auto-sizing based on number of browsers
- Static window size (no manual resize)
- Window exits cleanly on Esc, Q, or selection
- Embedded default fallback browsers and theme if no config file exists

### Known Limitations
- Window decorations (titlebar) still visible under Fyne 2.x
- Transparency limited due to Fyne's current Wayland/X11 support
- No domain-specific browser rules yet
- No live config reloading
- No Flathub/AUR packages yet (planned)

---

