# Cron Kuma Pusher

**Cron Kuma Pusher** is a lightweight Go utility that runs cron jobs defined in a YAML configuration file and pushes their execution results to your **Uptime Kuma** instance via the Push API.

Ideal for self-hosted environments, it bridges the gap between your scheduled tasks and uptime monitoring — keeping you informed when jobs fail or succeed.

---

## Features

- Parse and execute jobs from a simple `config.yml`
- Schedule jobs using standard cron expressions
- Send execution results directly to **Uptime Kuma** push monitors
- Minimal setup — one binary, one config file
- Written in Go for performance and portability

---

## Installation

Download the latest binary from the [GitHub Releases](https://github.com/Yohannfra/cron-kuma-pusher/releases) page and install it using `curl`:

```bash
curl -L https://github.com/Yohannfra/cron-kuma-pusher/releases/latest/download/cron-kuma-pusher -o /usr/local/bin/cron-kuma-pusher
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
       expression: "*/1 * * * *"
       command: "du -h"
       pushToken: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
   ```

   - `kumaBaseUrl`: your Uptime Kuma API base URL (ends with `/api/push/`).
   - `expression`: cron schedule (standard syntax).
   - `command`: the command to execute.
   - `pushToken`: your Kuma push token.

3. **Run the program**

   ```bash
   cron-kuma-pusher -config config.yml
   ```
---

## 🪪 License

This project is licensed under the **MIT License**, matching the license of [Uptime Kuma](https://github.com/louislam/uptime-kuma).
