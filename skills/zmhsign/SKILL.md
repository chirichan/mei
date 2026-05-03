---
name: zmhsign
description: Run and manage the locally installed zmhsign command-line program for a manga/comic website daily sign-in on Windows, macOS, or Linux. Use when the user asks to sign in, check in, run the comic site sign-in, start or stop scheduled/daemon sign-in, configure or reuse ZAI_USER/ZAI_PASS credentials, or troubleshoot zmhsign flags such as -user, -pass, -serve, -time, and -no-first. On Windows, check Process/User/Machine environment scopes before asking for credentials; on macOS/Linux, check the current process environment before asking.
---

# ZMHSIGN

Use this skill to run the user's locally installed `zmhsign` manga site sign-in CLI.

## Platform Rules

Use the executable name that matches the operating system:

```text
Windows:      zmhsign.exe
macOS/Linux:  zmhsign
```

Do not assume a local installation path, username, home directory, Go workspace, or shell profile path. Prefer `PATH` discovery first. If the executable is not discoverable, ask the user for the installed executable path or ask them to add the install directory to `PATH`.

## Command

Check whether the command is available.

Windows PowerShell:

```powershell
Get-Command zmhsign.exe
```

macOS/Linux shell:

```sh
command -v zmhsign
```

Show help.

Windows PowerShell:

```powershell
zmhsign.exe -h
```

macOS/Linux shell:

```sh
zmhsign -h
```

Supported flags:

```text
-user string      username
-pass string      password
-serve            run as daemon
-time string      daily run time (HH:MM), default 09:00
-no-first         skip first immediate run when serving
```

The program also reads credentials from environment variables:

```text
ZAI_USER
ZAI_PASS
```

## Credential Rules

Prefer existing environment variables for credentials instead of asking the user. Do not request `ZAI_USER` or `ZAI_PASS` until the platform-appropriate checks below have been performed.

If the user provides credentials in chat, use them only for the immediate command needed to complete the request. Do not store them in the skill, repository, logs, shell profile, or documentation unless the user explicitly asks to persist them. Do not print `ZAI_PASS` in command output or status summaries.

## Windows Workflow

On Windows, check `Process`, `User`, and `Machine` environment scopes. Refresh credentials into the current PowerShell process before running `zmhsign.exe`; this handles sessions that do not already inherit newly configured Windows environment variables.

```powershell
$zaiUser = [Environment]::GetEnvironmentVariable("ZAI_USER", "Process")
if ([string]::IsNullOrWhiteSpace($zaiUser)) {
  $zaiUser = [Environment]::GetEnvironmentVariable("ZAI_USER", "User")
}
if ([string]::IsNullOrWhiteSpace($zaiUser)) {
  $zaiUser = [Environment]::GetEnvironmentVariable("ZAI_USER", "Machine")
}

$zaiPass = [Environment]::GetEnvironmentVariable("ZAI_PASS", "Process")
if ([string]::IsNullOrWhiteSpace($zaiPass)) {
  $zaiPass = [Environment]::GetEnvironmentVariable("ZAI_PASS", "User")
}
if ([string]::IsNullOrWhiteSpace($zaiPass)) {
  $zaiPass = [Environment]::GetEnvironmentVariable("ZAI_PASS", "Machine")
}

if ([string]::IsNullOrWhiteSpace($zaiUser) -or [string]::IsNullOrWhiteSpace($zaiPass)) {
  throw "ZAI_USER or ZAI_PASS is missing from Process, User, and Machine environment scopes."
}

$env:ZAI_USER = $zaiUser
$env:ZAI_PASS = $zaiPass
```

Run once:

```powershell
zmhsign.exe
```

Run as a daemon immediately and then every day at 09:00:

```powershell
zmhsign.exe -serve -time 09:00
```

Run as a daemon without an immediate first run:

```powershell
zmhsign.exe -serve -time 09:00 -no-first
```

If the executable is not on `PATH` and the user provides a path, invoke it with PowerShell's call operator:

```powershell
& "<path-to-zmhsign.exe>" -h
```

## macOS/Linux Workflow

On macOS and Linux, use the current process environment. Check for credentials before asking the user:

```sh
test -n "$ZAI_USER" && test -n "$ZAI_PASS"
```

Run once:

```sh
zmhsign
```

Run as a daemon immediately and then every day at 09:00:

```sh
zmhsign -serve -time 09:00
```

Run as a daemon without an immediate first run:

```sh
zmhsign -serve -time 09:00 -no-first
```

For a one-off command with credentials supplied by the user, prefer inline environment variables for the child process instead of writing to a shell profile:

```sh
ZAI_USER="<username>" ZAI_PASS="<password>" zmhsign
```

If the executable is not on `PATH` and the user provides a path, invoke that path directly:

```sh
"<path-to-zmhsign>" -h
```

## Troubleshooting

If the program prints usage text, credentials were probably missing or empty in the child process. Re-run the platform-specific credential checks, then run the command again.

If macOS/Linux reports `permission denied`, ask the user whether they want to mark the binary executable, then run `chmod +x <path-to-zmhsign>` only with their approval.

If macOS/Linux reports `command not found`, verify `command -v zmhsign`; if missing, ask for the install path or ask the user to add the binary directory to `PATH`.

If the user asks for daily unattended execution, explain that `-serve` keeps the process running in the current session. For persistent startup, use a platform-specific supervisor only after the user asks for it: Windows Task Scheduler, macOS launchd, or Linux systemd/cron.

If the user asks to stop the daemon, find the running `zmhsign` process and confirm before terminating it unless they explicitly requested the stop action.
