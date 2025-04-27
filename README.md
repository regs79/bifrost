# 🌈 Bifrost

**Bifrost** is a lightweight, configurable link picker for Linux desktops.

When set as your default browser, Bifrost intercepts link clicks and displays a selection screen where you can quickly choose **which browser** to open the link with.

Inspired by the Norse "Bifrost bridge" — connecting worlds — Bifrost connects your clicks to the browser you want.

---

## ✨ Features

- Horizontal browser icon picker
- Keyboard shortcuts (1-9, Enter, Esc, Q)
- Mouse hover and selection highlighting
- Dynamic theming via config file
- Built-in defaults (works even without a config)
- Semi-transparent window support (Wayland/X11 friendly)
- Ultra-fast startup

---

## 📷 Screenshot

*(Coming soon)*

---

## ⚙️ Configuration

Bifrost looks for a config file at:

`$HOME/.config/bifrost/bifrost.toml`


If not found, Bifrost uses built-in defaults.

You can create your own `bifrost.toml` to customize:

- Browser list (icons, executable commands)
- Theme colors (background, text, highlight)

Example `bifrost.toml`:

```
[colors]
background = "#2E3440CC"
text = "#ECEFF4FF"
highlight = "#5E81ACAA"

[[browsers]]
name = "Firefox"
exec = "firefox"
icon = "assets/firefox.png"

[[browsers]]
name = "Brave"
exec = "brave"
icon = "assets/brave.png"
```

---

## 🔥 Roadmap

- AUR release
- Optional domain-based browser rules
- Config hot reload (no restart needed)
- Frameless window dragging
- Light/dark mode toggle
- Animated highlight transitions

See [`TODO.md`](TODO.md) for full development roadmap.

---

## 📜 License

MIT License.  
See [`LICENSE`](LICENSE) for full details.

---

## ❤️ Contributions

PRs welcome!  
If you have ideas, feature requests, or improvements, feel free to open an issue or submit a pull request.

---
