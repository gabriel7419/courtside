# Courtside — Desktop Notifications Setup

Courtside sends a desktop notification (with a beep) whenever a live game scores.

## macOS

Notifications work out of the box — no extra setup needed.

## Linux

Install `libnotify`:

```bash
# Debian/Ubuntu
sudo apt-get install libnotify-bin

# Arch
sudo pacman -S libnotify

# Fedora
sudo dnf install libnotify
```

## Windows

Notifications are delivered via the Windows toast system. No additional packages required.

---

## Notification Format

```
🏀 Courtside!

J. Tatum  Q3 4:52  [3PT +3 · BOS]
BOS  89 - 79  MIA
```

**Event labels:**

| Event | Label |
|---|---|
| 2-point field goal | `BASKET +2` |
| 3-point field goal | `3PT +3` |
| Free throw made | `FT +1` |

---

## Disabling Notifications

Courtside does not yet have a UI toggle for notifications. To disable them, run the binary without a notification daemon present (Linux) or deny notification permissions (macOS/Windows system settings).