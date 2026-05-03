---
name: zmhsign
description: Run and manage the locally installed zmhsign.exe command-line program for a manga/comic website daily sign-in. Use when the user asks to sign in, check in, run the comic site sign-in, start or stop scheduled/daemon sign-in, configure or reuse Windows ZAI_USER/ZAI_PASS credentials, or troubleshoot zmhsign.exe flags such as --user, --pass, --serve, --time, and --no-first. Always check Process/User/Machine environment scopes before asking for credentials.
---

# ZMHSIGN

Use this skill to run the user's locally installed `zmhsign.exe` manga site sign-in tool from PowerShell. The user may already have `ZAI_USER` and `ZAI_PASS` configured in Windows user-level or machine-level environment variables.

## Command

The installed program should usually be discoverable from `PATH`:

```powershell
zmhsign.exe -h
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

## Workflow

Prefer existing Windows environment variables for credentials instead of asking the user. Do not request `ZAI_USER` or `ZAI_PASS` until all three scopes have been checked: `Process`, `User`, and `Machine`.

Refresh credentials into the current PowerShell process before running `zmhsign.exe`. This handles Codex or shell sessions that do not already inherit newly configured Windows environment variables:

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

For a one-time sign-in:

```powershell
zmhsign.exe
```

For a daemon that runs immediately and then every day at 09:00:

```powershell
zmhsign.exe -serve -time 09:00
```

For a daemon that waits until the next scheduled time without an immediate first run:

```powershell
zmhsign.exe -serve -time 09:00 -no-first
```

If the user provides credentials in chat, use them only for the immediate command needed to complete the request. Do not store them in the skill, repository, logs, or documentation. Do not print `ZAI_PASS` in command output or status summaries.

## Checks

Before running the command, verify that `zmhsign.exe` is available:

```powershell
Get-Command zmhsign.exe
```

If it is not on `PATH`, do not assume a local user path. Ask the user for the installed executable path, or ask them to add the directory containing `zmhsign.exe` to `PATH`. If they provide an explicit path, use that path only for the current task and do not write it back into the skill.

```powershell
& "<path-to-zmhsign.exe>" -h
```

Check all Windows environment variable scopes before asking for credentials:

```powershell
[Environment]::GetEnvironmentVariable("ZAI_USER", "Process")
[Environment]::GetEnvironmentVariable("ZAI_USER", "User")
[Environment]::GetEnvironmentVariable("ZAI_USER", "Machine")
[Environment]::GetEnvironmentVariable("ZAI_PASS", "Process") -ne $null
[Environment]::GetEnvironmentVariable("ZAI_PASS", "User") -ne $null
[Environment]::GetEnvironmentVariable("ZAI_PASS", "Machine") -ne $null
```

Only ask the user for credentials if `ZAI_USER` or `ZAI_PASS` is absent from `Process`, `User`, and `Machine`. If values exist in `User` or `Machine` but not `Process`, copy them into `$env:ZAI_USER` and `$env:ZAI_PASS` for the current command instead of asking.

## Troubleshooting

If the program prints usage text, credentials were probably missing or empty in the child process. Re-run the environment refresh block from the Workflow section, then run `zmhsign.exe` again.

If the user asks for daily unattended execution, explain that `-serve` keeps the process running in the current session. For a persistent Windows startup or Task Scheduler setup, create a separate scheduler workflow only after the user asks for it.

If the user asks to stop the daemon, find the running `zmhsign` process and confirm before terminating it unless they explicitly requested the stop action.
