# Tech Stack Analysis & Project Feasibility Study

## 1. Executive Summary
The proposed stack (**Golang + Nuxt.js + Redis**) is **highly suitable** for a high-performance, scalable WhatsApp SaaS. The existing repository `go-whatsapp-web-multidevice` provides a solid starting point but requires significant architectural refactoring to support multi-tenancy (SaaS).

**Verdict:** **Proceed with the proposed stack.** The combination offers the best balance of performance (Go), developer experience (Nuxt), and scalability (Redis/PostgreSQL).

---

## 2. Tech Stack Component Analysis

### A. Backend: Golang (Fiber + whatsmeow)
*   **Why it works:**
    *   **Concurrency:** Go's goroutines are perfect for handling thousands of concurrent WebSocket connections (one per WhatsApp session) with minimal memory footprint compared to Node.js.
    *   **Performance:** `whatsmeow` is a Go library that reverse-engineers the WhatsApp Web protocol efficiently. It is currently the most stable open-source solution.
    *   **Fiber Framework:** Used in the existing repo, Fiber is extremely fast (Express-like but in Go) and suitable for high-load APIs.
*   **Gap Analysis (Existing Repo):**
    *   The current `init.go` uses global variables (`cli`, `db`) which effectively limits the application to a **single WhatsApp account**.
    *   **Action:** Must refactor global state into a `ClientManager` map (e.g., `map[TenantID]*whatsmeow.Client`) to handle multiple sessions simultaneously.

### B. Frontend: Nuxt.js (Vue 3)
*   **Why it works:**
    *   **SEO & Performance:** Server-Side Rendering (SSR) is crucial if you plan to have public landing pages or SEO-friendly directories for UMKM.
    *   **Developer Experience:** Nuxt 3's auto-imports and module ecosystem (Pinia, Tailwind) speed up development significantly.
    *   **Real-time:** Excellent support for WebSocket/SSE integration for live QR code streaming.
*   **Integration Strategy:**
    *   Nuxt will act as the "Control Panel".
    *   It will consume the Go Backend APIs.
    *   **Recommendation:** Use `nuxt-auth-utils` or `NextAuth` pattern for handling SaaS user sessions (JWT from Go backend).

### C. Database: PostgreSQL + Redis
*   **PostgreSQL (Primary DB):**
    *   **Role:** Store Tenant data, Billing info, and **WhatsApp Session Data** (replacing SQLite).
    *   **Why:** Relational integrity is non-negotiable for a SaaS with billing and multi-user access. JSONB support allows flexible storage for "AI Configs".
*   **Redis (Queue & Cache):**
    *   **Role:**
        1.  **Message Queue:** Buffer incoming WhatsApp messages before processing by AI (prevents bottlenecks).
        2.  **Session Cache:** Store active JWT tokens.
        3.  **Job Queue:** Schedule "Follow-up" messages (using `asynq`).
    *   **Necessity:** **Critical.** Without Redis, a spike in messages from one viral UMKM could crash the entire server for everyone else.

---

## 3. Deep Dive: Existing Repository Analysis (`go-whatsapp-web-multidevice`)

I have analyzed the source code, specifically `src/infrastructure/whatsapp/init.go`.

### **Strengths:**
*   **Clean Structure:** Uses Domain-Driven Design (DDD) principles (`domains`, `usecase`, `infrastructure`).
*   **PostgreSQL Support:** The code already has logic to switch between SQLite and Postgres (`initDatabase` function), which saves us time.
*   **Event Handling:** Good separation of event handlers (`handleMessage`, `handleImageMessage`).

### **Critical Weaknesses (Must Fix for SaaS):**
1.  **Global State Anti-Pattern:**
    ```go
    // src/infrastructure/whatsapp/init.go
    var (
        cli *whatsmeow.Client // <--- SINGLE CLIENT ONLY!
        db  *sqlstore.Container
    )
    ```
    *   **Impact:** The current code can only run **ONE** WhatsApp number at a time per server instance.
    *   **Fix:** Replace global `cli` with a `ClientManager` struct that manages a pool of clients.

2.  **Hardcoded Logic:**
    *   Auto-reply logic is hardcoded in `handleAutoReply`.
    *   **Fix:** Move this logic to a dynamic service that checks the `tenant_id` and their specific configuration in the DB.

3.  **Missing Tenant Context:**
    *   Database queries do not filter by `tenant_id`.
    *   **Fix:** All repository methods must accept `tenant_id` as a parameter.

---

## 4. Proposed Architecture Evolution

| Feature | Current Repo (Single Tenant) | Target SaaS Architecture (Multi-Tenant) |
| :--- | :--- | :--- |
| **Session Storage** | SQLite / Single Postgres Table | Postgres `device_store` table with `tenant_id` column |
| **Client Instance** | Global `var cli *Client` | `Map<TenantID, *Client>` in Memory |
| **Message Handling** | Direct Processing | Message -> **Redis Queue** -> Worker Pool |
| **Auth** | Basic / None | JWT Middleware (Scoped to Tenant) |
| **AI Integration** | None | Event-Driven AI Worker (OpenAI Consumer) |

## 5. Conclusion & Recommendation

**The path is clear.** We will use the existing repo as a **library of patterns**, but we should essentially **rewrite the core infrastructure layer** (`src/infrastructure/whatsapp`) to support multi-tenancy.

**Next Immediate Step:**
Start **Phase 1** of the roadmap: Initialize the new project structure and implement the **Multi-tenant Database Schema**, specifically focusing on how to store `whatsmeow` sessions in Postgres with a `tenant_id`.
