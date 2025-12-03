# Integration Analysis: Project Template + WhatsApp Multi-Device

## 1. Overview
We will utilize the provided `project-template` as the foundation for the **Gowa SaaS** project. This template provides a clean, containerized setup for Golang (Backend) and Nuxt.js (Frontend). We will enhance this template by integrating Redis and the core WhatsApp logic from the `go-whatsapp-web-multidevice` reference.

## 2. Template Structure Analysis

### Existing Architecture
*   **Infrastructure:** Docker Compose with PostgreSQL, Go API, and Nuxt App.
*   **Backend:** Simple Go HTTP server structure (`handlers`, `models`, `middleware`).
*   **Frontend:** Standard Nuxt 3 setup with TailwindCSS.
*   **Missing Components:**
    *   **Redis:** Required for message queuing and session caching.
    *   **WhatsApp Logic:** The template is generic; it lacks `whatsmeow` integration.
    *   **Multi-tenancy:** The template seems to have basic Auth (Google), but needs specific SaaS tenant logic.

## 3. Integration Strategy

### Step 1: Template Initialization
*   Copy contents of `project-template/template` to the project root.
*   Replace all `{{PROJECT_NAME}}` placeholders with `gowa`.
*   Fix `docker-compose.yml` to include **Redis**.

### Step 2: Backend Evolution (The "Brain" Transplant)
We need to port the logic from `go-whatsapp-web-multidevice` into the template's `backend` folder, but **refactored for SaaS**.

| Feature | Reference Repo (`go-whatsapp...`) | Target Implementation (Template `backend`) |
| :--- | :--- | :--- |
| **Library** | `whatsmeow` | Add `go.mau.fi/whatsmeow` to `go.mod` |
| **Session Store** | `sqlstore` (Global DB) | **Custom `PostgresStore`** in `services/whatsapp/store.go` |
| **Client Mgmt** | Global `cli` variable | `ClientManager` struct in `services/whatsapp/manager.go` |
| **Event Loop** | Direct Handler | `EventHandler` pushing to **Redis Queue** |

### Step 3: Frontend Evolution
*   Update `nuxt.config.ts` to proxy `/api` requests to the Go backend.
*   Install `pinia` (if not present) for state management.
*   Create a `useWhatsapp` composable to handle QR code streaming (SSE/WebSocket).

## 4. Detailed Implementation Steps

### 4.1. Infrastructure Setup
1.  **Move Files:** Copy template to root.
2.  **Update `docker-compose.yml`:**
    ```yaml
    redis:
      image: redis:7-alpine
      container_name: gowa-redis
      ports:
        - "6379:6379"
      volumes:
        - redis_data:/data
    ```
3.  **Update `backend/Dockerfile`:** Ensure `git` and build tools are installed (for `whatsmeow` CGO dependencies if needed, though pure Go is preferred).

### 4.2. Backend Implementation
1.  **Dependencies:** Run `go get go.mau.fi/whatsmeow` and `go get github.com/redis/go-redis/v9`.
2.  **Database Schema:** Create a migration in `backend/migrations` for:
    *   `tenants`
    *   `devices` (stores `whatsmeow` session data)
3.  **Service Layer:** Create `backend/services/whatsapp`:
    *   `manager.go`: Holds the `map[string]*whatsmeow.Client`.
    *   `client.go`: Logic to Connect/Disconnect.
    *   `events.go`: Logic to handle incoming messages.

### 4.3. Frontend Implementation
1.  **Dashboard Page:** Create `pages/dashboard/index.vue`.
2.  **QR Component:** Create `components/WhatsappQR.vue` that listens to an SSE endpoint.

## 5. Conclusion
Using the template is the right choice. It saves setup time for the "boring" parts (Docker, basic HTTP). Our focus will be solely on injecting the **WhatsApp Intelligence** into the backend service.

**Ready to execute Step 1 (Initialize Template)?**
