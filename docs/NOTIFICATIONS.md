# Notifications

Score notifications require one-time setup depending on your operating system.

## macOS

Notifications use AppleScript, which requires enabling permissions for Script Editor:

1. Open **Script Editor** (`/Applications/Utilities/Script Editor.app`)
2. Run this to test: `display notification "test" with title "test"`
3. Go to **System Settings → Notifications → Script Editor**
4. Enable notifications and set the alert style to "Banners"

## Linux

Notifications require `libnotify`. Install it if not present:

```bash
# Debian/Ubuntu
sudo apt install libnotify-bin

# Fedora
sudo dnf install libnotify

# Arch
sudo pacman -S libnotify
```

## Windows

Notifications work out of the box on Windows 10 and 11.

## Notification Types

Courtside can notify you about:

- Score changes during live games
- End of quarter / halftime
- Game tips off (optional)

Configure which events trigger notifications in the **Settings** view.