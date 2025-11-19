# Cron Kuma Pusher

**Cron Kuma Pusher** is a lightweight Go utility that runs cron jobs defined in a YAML configuration file and pushes their execution results to your **Uptime Kuma** instance via the Push API.

Ideal for self-hosted environments, it bridges the gap between your scheduled tasks and uptime monitoring — keeping you informed when jobs fail or succeed.

---

## Features

- Parse and execute jobs from a simple config **YAML** file
- Schedule jobs using standard cron expressions
- Send execution results directly to **Uptime Kuma** push monitors
- Minimal setup — one binary, one config file
- Written in Go for performance and portability

---

## Installation

Download the latest binary from the [GitHub Releases](https://github.com/Yohannfra/cron-kuma-pusher/releases) using `curl`:

```bash
# darwin-amd64 and windows-amd64 are also available on the release page

# linux-amd64
curl -L https://github.com/Yohannfra/cron-kuma-pusher/releases/download/v0.0.3/cron-kuma-pusher-linux-amd64 -o /usr/local/bin/cron-kuma-pusher

# darwin-arm64
curl -L https://github.com/Yohannfra/cron-kuma-pusher/releases/download/v0.0.3/cron-kuma-pusher-darwin-arm64 -o /usr/local/bin/cron-kuma-pusher

chmod +x /usr/local/bin/cron-kuma-pusher
```

---

## Usage

1. **Create a Push Monitor in Uptime Kuma**
   - Go to your Uptime Kuma dashboard.
   - Click **Add New Monitor** → choose **Push**.
   - Copy the **Url** and the **Push Token** (the part after `/push/`).

2. **Create your configuration file**

   Save the following as `config.yml`:

   ```yaml
   kumaBaseUrl: 'https://uptime.yourdomain.com/api/push/'

   jobs:
     - name: example job
       expression: "0 * * * * *" # every minute
       command: "du -h"
       pushToken: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
   ```

   - `kumaBaseUrl`: your Uptime Kuma API base URL (ends with `/api/push/`).
   - `expression`: cron schedule (standard syntax).
   - `command`: the command to execute.
   - `pushToken`: your Kuma push token.

3. **Run the program**

   ```bash
   cron-kuma-pusher --config config.yml
   ```
---

## Running as a Systemd Service (Linux)

To keep **Cron Kuma Pusher** running in the background and automatically restart on reboot or crash, you can configure it as a `systemd` service.

1. **Create a systemd unit file**

```bash
$ sudo nano /etc/systemd/system/cron-kuma-pusher.service
```

Paste the following content (adjust it as needed):

```ini
[Unit]
Description=Cron Kuma Pusher Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cron-kuma-pusher --config /etc/cron-kuma-pusher/config.yml
WorkingDirectory=/etc/cron-kuma-pusher
Restart=always
RestartSec=10
User=ubuntu
Group=ubuntu
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

2. Create the config directory

```bash
$ sudo mkdir -p /etc/cron-kuma-pusher
$ sudo cp config.yml /etc/cron-kuma-pusher/
```

3. Reload and enable the service

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable cron-kuma-pusher
$ sudo systemctl start cron-kuma-pusher
```

4. Check the status

```bash
$ systemctl status cron-kuma-pusher
```

5. View logs (optional)

```bash
journalctl -u cron-kuma-pusher -f
```

---

## License

This project is licensed under the **MIT License**, matching the license of [Uptime Kuma](https://github.com/louislam/uptime-kuma).
