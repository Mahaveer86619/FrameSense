# FrameSense

## AI-Aware Contextual Advertising Engine for Video Streaming

FrameSense is an experimental video streaming platform that injects context-aware, non-interruptive advertisements directly into video content.

Instead of traditional mid-roll ads, FrameSense displays subtle cartoon-based in-scene ads that blend naturally into ongoing scenes without pausing playback.  
The system is built using a Go-based streaming backend, Python AI pipelines, and server-side HLS manipulation, with all cloud dependencies simulated locally using LocalStack.

---

## Key Features

- End-to-end video ingestion and streaming pipeline  
- Subtitle-based scene metadata extraction  
- Contextual cartoon character ad overlays  
- Server-Side Ad Insertion (SSAI) using HLS  
- Microservices architecture (Go + Python)  
- AWS-like infrastructure using LocalStack  
- Fully containerized with Docker Compose  

---

## How FrameSense Works

### 1. Content Ingestion
- Videos are uploaded via API  
- Backend converts videos into HLS format  
- Subtitles are extracted and analyzed for scene context  

### 2. Scene Metadata Creation (Subtitle-Based Analysis)
Instead of manual tagging or heavy computer vision models, FrameSense uses subtitle analysis to infer scene metadata.

Extracted metadata includes:
- Scene type (dialogue, calm, intense)
- Emotional tone (neutral, emotional, humorous)
- Time ranges suitable for subtle ads

This approach is:
- Lightweight  
- Deterministic  
- AI-ready (easy to replace with ML later)  

### 3. Ad Decision Engine
A Go service evaluates:
- Scene metadata  
- Ad rules  
- Overlay availability  

It produces an ad insertion plan describing:
- When an ad should appear  
- Which cartoon overlay to use  
- How long it should stay visible  

### 4. AI Ad Generation Pipeline (Python Microservice)
A dedicated Python microservice is responsible for ad asset generation.

Responsibilities:
- Generate or customize cartoon overlay videos  
- Add branding, text, or subtle motion  
- Export transparent overlay assets  

This service is decoupled to:
- Avoid blocking the streaming pipeline  
- Allow future ML or generative AI upgrades  

### 5. Server-Side Ad Injection
- FFmpeg overlays cartoon ads onto selected video segments  
- Modified HLS playlists are served dynamically  
- From the player’s perspective, playback is seamless  

### 6. Custom Streaming Player
- Built on HTML5 Video + hls.js  
- No ad SDKs  
- No client-side ad logic  
- Plays a single continuous stream  

---

## Architecture Overview

+------------------+
|   Web Player     |
| (HTML5 + HLS)    |
+--------+---------+
         |
         v
+------------------+
|  API Gateway     |  (Go)
+--------+---------+
         |
+--------+---------+
| Content Service  |  (Go)
+--------+---------+
         |
+--------+---------+
| Scene Analyzer   |  (Go)
| (Subtitle NLP)   |
+--------+---------+
         |
+--------+---------+        +----------------------+
| Ad Decision      | -----> | AI Ad Generator      |
| Engine (Go)      |        | (Python Microservice)|
+--------+---------+        +----------------------+
         |
+--------+---------+
| Stream Stitcher  |  (Go + FFmpeg)
+--------+---------+
         |
+--------+---------+
|   HLS Output     |
+------------------+


---

## Tech Stack

### Backend
- Go (Gin / Fiber)
- FFmpeg
- HLS (m3u8 manipulation)
- Redis (ad decision cache)
- PostgreSQL / SQLite

### AI & Processing
- Python (FastAPI)
- Subtitle sentiment analysis
- Overlay asset generation

### Frontend
- HTML5 Video API
- hls.js
- Minimal JavaScript controller

### Infrastructure
- Docker & Docker Compose
- LocalStack (S3, SQS, DynamoDB simulation)
- No real AWS account required

---

## LocalStack Usage

All AWS-style services are emulated locally:

- S3 → Video & overlay storage  
- SQS → Async processing jobs  
- DynamoDB → Metadata & rules  

This enables:
- Full cloud-style architecture  
- Zero cloud cost  
- Easy local development  

---

## Running the Project

```bash
docker-compose up --build
