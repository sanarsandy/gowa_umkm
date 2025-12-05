# Setup Google OAuth untuk Gowa UMKM

## 🎯 Overview

Aplikasi ini sudah memiliki implementasi Google OAuth yang lengkap. Anda hanya perlu:
1. Membuat Google Cloud Project
2. Mendapatkan OAuth 2.0 credentials
3. Mengkonfigurasi environment variables

---

## 📋 Langkah-langkah Setup

### Step 1: Buat Google Cloud Project

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Klik **Select a project** → **New Project**
3. Isi nama project (contoh: "Gowa UMKM")
4. Klik **Create**

### Step 2: Enable OAuth Consent Screen

1. Di sidebar, pilih **APIs & Services** → **OAuth consent screen**
2. Pilih **External** (untuk umum) atau **Internal** (hanya org Google Workspace)
3. Klik **Create**
4. Isi informasi:
   - **App name**: Gowa UMKM
   - **User support email**: email Anda
   - **Developer contact email**: email Anda
5. Klik **Save and Continue**
6. Di **Scopes**, klik **Add or Remove Scopes**
   - Centang: `openid`, `email`, `profile`
   - Klik **Update** → **Save and Continue**
7. Di **Test users** (jika External), Anda bisa skip atau tambahkan email tester
8. Klik **Save and Continue** → **Back to Dashboard**

### Step 3: Buat OAuth 2.0 Credentials

1. Di sidebar, pilih **APIs & Services** → **Credentials**
2. Klik **+ Create Credentials** → **OAuth client ID**
3. **Application type**: Web application
4. **Name**: Gowa UMKM Web Client
5. **Authorized JavaScript origins**:
   ```
   http://localhost:3000
   https://app3.anakhebat.web.id
   ```
6. **Authorized redirect URIs**:
   ```
   http://localhost:8080/api/auth/google/callback
   https://app3.anakhebat.web.id/api/auth/google/callback
   ```
7. Klik **Create**
8. **Copy** Client ID dan Client Secret

### Step 4: Configure Environment

Edit file `.env` di local:

```env
# Google OAuth 2.0
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# Frontend URL (untuk redirect setelah login)
FRONTEND_URL=http://localhost:3000
```

Untuk **Production** di VPS:

```env
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=https://app3.anakhebat.web.id/api/auth/google/callback
FRONTEND_URL=https://app3.anakhebat.web.id
```

### Step 5: Restart Application

```bash
# Development
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.dev.yml up -d

# Production (di VPS)
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d
```

---

## ✅ Verifikasi

1. Buka halaman login: `http://localhost:3000/login`
2. Klik tombol **Google**
3. Pilih akun Google
4. Harus redirect ke dashboard setelah berhasil login

---

## 🔧 Troubleshooting

### Error: "Google OAuth is not configured"
- Pastikan `GOOGLE_CLIENT_ID` dan `GOOGLE_CLIENT_SECRET` sudah di-set
- Restart container setelah mengubah `.env`

### Error: "redirect_uri_mismatch"
- Cek **Authorized redirect URIs** di Google Console
- Pastikan URL sama persis (termasuk trailing slash)

### Error: "Access blocked: This app's request is invalid"
- Cek Authorized JavaScript origins
- Pastikan domain/port sesuai

---

## 📁 File Reference

| File | Deskripsi |
|------|-----------|
| `backend/handlers/google_auth.go` | OAuth handler (backend) |
| `backend/main.go` | Route registration (line 140-141) |
| `frontend/pages/login.vue` | Login button |
| `frontend/pages/register.vue` | Register button |
| `frontend/pages/auth/google/callback.vue` | Callback handler |
