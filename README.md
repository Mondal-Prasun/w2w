# ⚡ w2w: Distributed Video Processing Engine

A lightweight, high-performance asynchronous video transcoding engine built with **Go**, **Redis**, and **FFmpeg**. Designed to handle heavy media tasks like HLS segmentation, video rescaling, and automated thumbnail extraction without blocking your main server.

https://github.com/user-attachments/assets/6c5d43a9-699f-4f9f-a508-6e089953ab94

---

## ✨ Features

* 🚀 **Asynchronous Queue Pipeline:** Offloads heavy encoding jobs to a background queue powered by **Redis BLPop**.
* 🎥 **Multi-Format Processing:**
  * **HLS Segmentation:** Auto-generates adaptive `.m3u8` playlists and `.ts` media segments.
  * **Video Rescaling:** Smart scaling and padding to strict target aspect ratios.
  * **Thumbnail Extraction:** Ultra-fast single-frame extraction.
* 🛡️ **Bounded Worker Pool:** Throttled Go goroutine architecture to prevent CPU thrashing and memory spikes during concurrent batch jobs.
* 📦 **Single-Binary Deployment:** Built with Go's `embed.FS` to bundle static frontend assets and templates directly into one binary executable.
* 📊 **Interactive Dashboard:** Minimalist dark-mode UI with batch staging, real-time health checks, dynamic status polling, and instant `.zip` result downloads.

---

## 🛠️ Tech Stack

* **Backend:** Go (Golang) Standard Library (`net/http`, `os/exec`, `embed`)
* **Task Queue & State:** Redis
* **Media Engine:** FFmpeg
* **Frontend:** Vanilla JavaScript (ES6+), Tailwind CSS

---

## 🏗️ Architecture Overview

```text
[ Client Dashboard ]
         │
         ▼ (HTTP POST /acceptJob)
 [ Go API Server ] ──► [ Redis Pipeline (Push Job ID & State) ]
                               │
                               ▼ (BLPop Worker Queue)
                     [ Bounded Worker Pool ]
                               │
                       (exec.Command)
                               │
                         [ FFmpeg ]
                               │
                               ▼
                    [ ZIP Output Storage ]
```

---

## 🚀 Quick Start

### Prerequisites
* Docker & Docker Compose

### 🐳 Run with Docker Compose (Recommended)

1. **Clone the Repository**
   ```bash
   git clone https://github.com/Mondal-Prasun/w2w.git
   cd w2w
   ```

2. **Start the Application Stack**
   ```bash
   docker-compose up --build -d
   ```

3. **Access the Dashboard**
   Open your browser and navigate to `http://localhost:8080`.

---

## 💻 Manual Setup (Without Docker)

If you prefer to run services individually without Docker:

### Additional Requirements:
* Go 1.22+
* Redis Server
* FFmpeg installed and added to your system PATH

1. **Clone the Repository**
   ```bash
   git clone https://github.com/Mondal-Prasun/w2w.git
   cd w2w
   ```

2. **Set Environment Variables**
   ```bash
   export SERVER_ADDR=":8080"
   export REDIS_ADDR="localhost:6379"
   ```

3. **Run the Server**
   ```bash
   go run main.go
   ```

---

## 🔌 API Reference

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/checkHealth` | `GET` | Health indicator for server status check |
| `/acceptJob` | `POST` | Accepts multipart form with `jobType` and `jobFile` |
| `/getJobDetail/{jobId}` | `GET` | Polls current status (Pending, Done, Failed) |
| `/download/{jobId}` | `GET` | Archives output folder on-the-fly and streams `.zip` |

---

## ⚙️ Configuration Hints

To tweak the maximum parallel FFmpeg executions, update `workerNumber` inside `worker/processJob.go`:

```go
const workerNumber uint8 = 3 // Adjust based on available CPU cores
```

To limit individual FFmpeg encoding threads, modify the `-threads` flag inside `worker/videoProcess.go`.

---

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.



