# Implementation Roadmap: WhatsApp UMKM SaaS

This document outlines the step-by-step plan to build the AI-Powered WhatsApp Agent for UMKM, based on the analysis in `analisis_project.md`.

## Phase 1: Project Initialization & Infrastructure (Week 1)

### 1.1. Project Structure Setup
- [ ] Create a monorepo or separate repositories for `backend` and `frontend`.
- [ ] Initialize Git repository.
- [ ] Create `docker-compose.yml` for local development (PostgreSQL, Redis).

### 1.2. Database Design (PostgreSQL) - **Multi-tenant Core**
- [ ] Design schema for `tenants` (users/businesses). **Critical:** All subsequent tables must reference `tenant_id`.
- [ ] Design schema for `whatsapp_sessions` (linked to `tenant_id`).
- [ ] Design schema for `messages` (partitioned or indexed by `tenant_id`).
- [ ] Design schema for `ai_configs` (stores prompts, business info, and active modes per tenant).
- [ ] Create migration scripts.

### 1.3. Backend Foundation (Golang)
- [ ] Initialize Go module (e.g., `go mod init gowa-saas`).
- [ ] Install core dependencies:
    - `github.com/gofiber/fiber/v2` (Web Framework)
    - `go.mau.fi/whatsmeow` (WhatsApp Library)
    - `gorm.io/gorm` or `github.com/jackc/pgx` (Database)
- [ ] Setup Database Connection (PostgreSQL).
- [ ] **Implement Tenant Middleware:** Ensure every API request is scoped to the authenticated tenant.
- [ ] **CRITICAL:** Implement Custom `DeviceStore` for `whatsmeow` backed by PostgreSQL.
    - *Challenge:* Map `whatsmeow`'s internal JID to our `tenant_id` to ensure the correct session is loaded.

## Phase 2: Core Backend Features (Week 2)

### 2.1. Session Management API
- [ ] Implement `POST /api/v1/device/connect`:
    - Create a new `whatsmeow` client instance.
    - Generate QR Code.
    - Return QR Code string/image to client.
- [ ] Implement `GET /api/v1/device/status`:
    - Check connection status of the specific user's client.
- [ ] Implement `POST /api/v1/device/disconnect`:
    - Logout and cleanup session.

### 2.2. WebSocket/SSE for Real-time QR
- [ ] Implement a WebSocket or Server-Sent Events (SSE) endpoint to stream QR code updates to the frontend in real-time.
- [ ] Handle session events (Connected, Timeout, Logged Out).

## Phase 3: Frontend MVP (Nuxt.js) (Week 3)

### 3.1. Setup & UI Framework
- [ ] Initialize Nuxt 3 project.
- [ ] Install and configure TailwindCSS.
- [ ] Setup State Management (Pinia).

### 3.2. Authentication (SaaS Layer)
- [ ] Create Login & Register pages.
- [ ] Implement JWT Authentication logic (connect to Backend Auth API).

### 3.3. Dashboard & Device Connection
- [ ] Create `Dashboard` layout.
- [ ] Build **Device Manager** component:
    - Button to "Connect WhatsApp".
    - Display QR Code (render from WebSocket/API).
    - Show Connection Status (Online/Offline).

## Phase 4: Intelligence Layer (AI & Messaging) (Week 4-5)

### 4.1. Message Listener & Queue
- [ ] Setup Redis.
- [ ] Implement `whatsmeow` event handler for `Message` events.
- [ ] **Tenant Router:** When a message arrives, identify which Tenant it belongs to before pushing to Redis.
- [ ] Push incoming messages to a Redis Queue (to avoid blocking the main thread).

### 4.2. AI Worker Service (The "Brain")
- [ ] Create a separate worker process (or goroutine pool) to consume from Redis.
- [ ] **Dynamic Context Injection:**
    - Fetch `ai_configs` for the specific Tenant (e.g., "You are a barista at Kopi Kenangan...").
    - Inject Business Knowledge (Menu, Hours, FAQs) into the System Prompt.
- [ ] **Mode Selection Logic:**
    - Implement logic to switch modes: *General FAQ* vs *Order Taking* vs *Complaint*.
- [ ] Integrate OpenAI API (or similar).
- [ ] Process message: Input -> LLM -> JSON Output.
- [ ] Save results to `customer_insights` table.

### 4.3. CRM Features on Frontend
- [ ] **AI Configuration UI:** Allow users to upload their "Business Profile" (Menu, Price list) which feeds into the AI Context.
- [ ] Create `Inbox` or `Customer List` view.
- [ ] Display AI Tags (e.g., "Hot Lead", "Complaint") next to chats.
- [ ] Show extracted details (Intent, Product Interest).

## Phase 5: Automation & Production Readiness (Week 6+)

### 5.1. Automated Follow-up (Scheduler)
- [ ] Implement `asynq` (Go library) for task scheduling.
- [ ] Logic: If AI detects "Needs Follow-up", schedule a job for X hours later.
- [ ] Worker executes job: Send WhatsApp message via `whatsmeow`.

### 5.2. Deployment & Scaling
- [ ] Dockerize Backend and Frontend.
- [ ] Setup Nginx as Reverse Proxy.
- [ ] **Load Testing:** Test with 10-50 concurrent sessions to monitor RAM usage.
- [ ] Implement Graceful Shutdown to save sessions correctly.

## Phase 6: Monetization (Optional/Later)
- [ ] Integrate Payment Gateway (Xendit/Midtrans).
- [ ] Implement Subscription Management (limit features based on plan).
