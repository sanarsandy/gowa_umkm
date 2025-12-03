<template>
  <div class="min-h-screen font-sans selection:bg-blue-100 selection:text-blue-900">
    <!-- Navigation -->
    <nav class="fixed w-full z-50 bg-white/80 backdrop-blur-xl border-b border-slate-100/50 transition-all duration-300">
      <div class="container mx-auto px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-br from-blue-600 to-sky-500 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/20">
              <span class="text-white font-bold text-xl">G</span>
            </div>
            <span class="text-xl font-bold text-slate-800 tracking-tight">
              Gowa UMKM
            </span>
          </div>
          <div class="hidden md:flex items-center space-x-8">
            <a href="#features" class="text-sm font-medium text-slate-600 hover:text-blue-600 transition">Fitur</a>
            <a href="#pricing" class="text-sm font-medium text-slate-600 hover:text-blue-600 transition">Harga</a>
            <a href="#faq" class="text-sm font-medium text-slate-600 hover:text-blue-600 transition">FAQ</a>
          </div>
          <div class="flex items-center space-x-4">
            <NuxtLink to="/login" class="text-sm font-medium text-slate-600 hover:text-blue-600 transition">
              Masuk
            </NuxtLink>
            <NuxtLink to="/register" class="px-6 py-2.5 bg-blue-600 text-white text-sm font-medium rounded-full hover:bg-blue-700 transition shadow-lg shadow-blue-600/20 hover:shadow-blue-600/30 transform hover:-translate-y-0.5">
              Daftar Gratis
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="relative pt-32 pb-20 overflow-hidden bg-slate-50">
      <!-- Soft Background Gradients -->
      <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
        <div class="absolute -top-[20%] -left-[10%] w-[50%] h-[50%] rounded-full bg-blue-100/50 blur-3xl"></div>
        <div class="absolute top-[10%] right-[0%] w-[40%] h-[40%] rounded-full bg-sky-100/50 blur-3xl"></div>
      </div>

      <div class="container mx-auto px-6 relative z-10">
        <div class="max-w-4xl mx-auto text-center">
          <div class="inline-flex items-center space-x-2 px-4 py-2 bg-white border border-blue-100 rounded-full mb-8 shadow-sm">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-sky-500"></span>
            </span>
            <span class="text-xs font-semibold text-slate-600 tracking-wide uppercase">Platform WhatsApp Automation #1</span>
          </div>
          <h1 class="text-5xl md:text-7xl font-extrabold text-slate-900 leading-tight mb-8 tracking-tight">
            Tingkatkan Penjualan dengan <br/>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-sky-500">Asisten AI Cerdas</span>
          </h1>
          <p class="text-xl text-slate-600 leading-relaxed mb-10 max-w-2xl mx-auto">
            Otomatisasi balasan WhatsApp, kelola pelanggan, dan analisis performa bisnis Anda dalam satu dashboard yang <span class="font-semibold text-slate-800">simpel & powerful</span>.
          </p>
          <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
            <NuxtLink to="/register" class="w-full sm:w-auto px-8 py-4 bg-blue-600 text-white rounded-full font-bold hover:bg-blue-700 transition shadow-xl shadow-blue-600/20 hover:shadow-blue-600/30 transform hover:-translate-y-1">
              Mulai Gratis Sekarang
            </NuxtLink>
            <a href="#demo" class="w-full sm:w-auto px-8 py-4 bg-white text-slate-700 border border-slate-200 rounded-full font-bold hover:bg-slate-50 hover:border-slate-300 transition flex items-center justify-center space-x-2 shadow-sm hover:shadow-md">
              <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span>Lihat Demo</span>
            </a>
          </div>
          
          <div class="mt-12 flex flex-wrap items-center justify-center gap-x-8 gap-y-4 text-sm font-medium text-slate-500">
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>Tanpa Kartu Kredit</span>
            </div>
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>Setup 2 Menit</span>
            </div>
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>Batal Kapan Saja</span>
            </div>
          </div>
        </div>

        <!-- Dashboard Preview -->
        <div class="mt-20 relative max-w-5xl mx-auto">
          <!-- Glow Effect -->
          <div class="absolute -inset-1 bg-gradient-to-r from-blue-500 to-sky-500 rounded-3xl blur opacity-20"></div>
          
          <div class="relative bg-white rounded-2xl p-2 md:p-4 border border-slate-200 shadow-2xl">
            <div class="bg-slate-50 rounded-xl overflow-hidden border border-slate-100">
              <!-- Mock Browser Header -->
              <div class="bg-white border-b border-slate-100 px-4 py-3 flex items-center space-x-2">
                <div class="flex space-x-1.5">
                  <div class="w-3 h-3 rounded-full bg-red-400/80"></div>
                  <div class="w-3 h-3 rounded-full bg-amber-400/80"></div>
                  <div class="w-3 h-3 rounded-full bg-green-400/80"></div>
                </div>
                <div class="flex-1 text-center">
                  <div class="inline-block px-4 py-1 bg-slate-50 rounded-lg text-xs text-slate-400 font-medium flex items-center justify-center space-x-2">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
                    <span>dashboard.gowa.id</span>
                  </div>
                </div>
              </div>
              <!-- Mock Content -->
              <div class="p-8 grid grid-cols-12 gap-6 bg-slate-50/50">
                <div class="col-span-3 space-y-4">
                  <div class="h-8 w-3/4 bg-white rounded-lg shadow-sm"></div>
                  <div class="h-4 w-1/2 bg-slate-200/50 rounded-lg"></div>
                  <div class="space-y-2 mt-8">
                    <div class="h-10 w-full bg-blue-600 rounded-lg shadow-lg shadow-blue-600/20"></div>
                    <div class="h-10 w-full bg-white border border-slate-200 rounded-lg"></div>
                    <div class="h-10 w-full bg-white border border-slate-200 rounded-lg"></div>
                  </div>
                </div>
                <div class="col-span-9 space-y-6">
                  <div class="flex space-x-4">
                    <div class="h-32 w-1/3 bg-white rounded-xl border border-slate-100 shadow-sm"></div>
                    <div class="h-32 w-1/3 bg-white rounded-xl border border-slate-100 shadow-sm"></div>
                    <div class="h-32 w-1/3 bg-white rounded-xl border border-slate-100 shadow-sm"></div>
                  </div>
                  <div class="h-64 w-full bg-white rounded-xl border border-slate-100 shadow-sm"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-24 bg-white relative">
      <div class="container mx-auto px-6">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-slate-900 mb-4">Fitur Esensial untuk Pertumbuhan</h2>
          <p class="text-lg text-slate-600">Kami fokus pada fitur yang benar-benar memberikan dampak pada bisnis Anda.</p>
        </div>
        
        <div class="grid md:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="bg-slate-50 p-8 rounded-3xl border border-slate-100 hover:border-blue-100 transition group hover:shadow-xl hover:shadow-blue-100/50 duration-300">
            <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition duration-300 border border-slate-100">
              <svg class="w-7 h-7 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-slate-900 mb-3">Smart Auto-Reply</h3>
            <p class="text-slate-600 leading-relaxed">
              Biarkan AI menangani pertanyaan berulang pelanggan Anda 24/7 dengan bahasa yang natural dan personal.
            </p>
          </div>

          <!-- Feature 2 -->
          <div class="bg-slate-50 p-8 rounded-3xl border border-slate-100 hover:border-blue-100 transition group hover:shadow-xl hover:shadow-blue-100/50 duration-300">
            <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition duration-300 border border-slate-100">
              <svg class="w-7 h-7 text-sky-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-slate-900 mb-3">Lead Detection</h3>
            <p class="text-slate-600 leading-relaxed">
              Sistem otomatis mengidentifikasi dan mengkategorikan pelanggan potensial (Hot/Warm/Cold) dari percakapan.
            </p>
          </div>

          <!-- Feature 3 -->
          <div class="bg-slate-50 p-8 rounded-3xl border border-slate-100 hover:border-blue-100 transition group hover:shadow-xl hover:shadow-blue-100/50 duration-300">
            <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition duration-300 border border-slate-100">
              <svg class="w-7 h-7 text-teal-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-slate-900 mb-3">Customer Insights</h3>
            <p class="text-slate-600 leading-relaxed">
              Dapatkan data mendalam tentang perilaku dan sentimen pelanggan untuk keputusan bisnis yang lebih baik.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Pricing Section -->
    <section id="pricing" class="py-24 bg-slate-50 relative overflow-hidden">
      <div class="container mx-auto px-6 relative z-10">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-slate-900 mb-4">Harga Simpel & Transparan</h2>
          <p class="text-lg text-slate-600">Pilih paket yang sesuai dengan kebutuhan bisnis Anda. Tanpa biaya tersembunyi.</p>
        </div>

        <div class="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <!-- Starter -->
          <div class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-blue-100/50 transition duration-300 relative">
            <h3 class="text-xl font-bold text-slate-900 mb-2">Starter</h3>
            <div class="flex items-baseline mb-6">
              <span class="text-4xl font-extrabold text-slate-900">Gratis</span>
              <span class="text-slate-500 ml-2">/ selamanya</span>
            </div>
            <p class="text-slate-600 mb-8 text-sm">Untuk UMKM yang baru memulai digitalisasi.</p>
            <ul class="space-y-4 mb-8">
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                1 Device WhatsApp
              </li>
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                100 Pesan Otomatis/bln
              </li>
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Basic Analytics
              </li>
            </ul>
            <NuxtLink to="/register" class="block w-full py-3 px-6 text-center rounded-xl border-2 border-slate-100 text-slate-600 font-bold hover:border-blue-600 hover:text-blue-600 transition">
              Mulai Gratis
            </NuxtLink>
          </div>

          <!-- Pro -->
          <div class="bg-white p-8 rounded-3xl border-2 border-blue-500 shadow-xl shadow-blue-500/10 relative transform md:-translate-y-4">
            <div class="absolute top-0 right-0 bg-blue-600 text-white text-xs font-bold px-3 py-1 rounded-bl-xl rounded-tr-2xl">POPULAR</div>
            <h3 class="text-xl font-bold text-slate-900 mb-2">Pro</h3>
            <div class="flex items-baseline mb-6">
              <span class="text-4xl font-extrabold text-slate-900">Rp 99rb</span>
              <span class="text-slate-500 ml-2">/ bulan</span>
            </div>
            <p class="text-slate-600 mb-8 text-sm">Untuk bisnis yang sedang bertumbuh pesat.</p>
            <ul class="space-y-4 mb-8">
              <li class="flex items-center text-slate-900 text-sm font-medium">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                3 Device WhatsApp
              </li>
              <li class="flex items-center text-slate-900 text-sm font-medium">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Unlimited Pesan Otomatis
              </li>
              <li class="flex items-center text-slate-900 text-sm font-medium">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Advanced AI Insights
              </li>
              <li class="flex items-center text-slate-900 text-sm font-medium">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Prioritas Support
              </li>
            </ul>
            <NuxtLink to="/register" class="block w-full py-3 px-6 text-center rounded-xl bg-blue-600 text-white font-bold hover:bg-blue-700 transition shadow-lg shadow-blue-600/20">
              Pilih Paket Pro
            </NuxtLink>
          </div>

          <!-- Business -->
          <div class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-blue-100/50 transition duration-300 relative">
            <h3 class="text-xl font-bold text-slate-900 mb-2">Business</h3>
            <div class="flex items-baseline mb-6">
              <span class="text-4xl font-extrabold text-slate-900">Rp 199rb</span>
              <span class="text-slate-500 ml-2">/ bulan</span>
            </div>
            <p class="text-slate-600 mb-8 text-sm">Solusi lengkap untuk operasional skala besar.</p>
            <ul class="space-y-4 mb-8">
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                10 Device WhatsApp
              </li>
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Custom AI Training
              </li>
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                API Access
              </li>
              <li class="flex items-center text-slate-600 text-sm">
                <svg class="w-5 h-5 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                Dedicated Account Manager
              </li>
            </ul>
            <NuxtLink to="/register" class="block w-full py-3 px-6 text-center rounded-xl border-2 border-slate-100 text-slate-600 font-bold hover:border-blue-600 hover:text-blue-600 transition">
              Hubungi Sales
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ Section -->
    <section id="faq" class="py-24 bg-white">
      <div class="container mx-auto px-6">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-slate-900 mb-4">Pertanyaan Umum</h2>
          <p class="text-lg text-slate-600">Jawaban untuk pertanyaan yang sering diajukan.</p>
        </div>

        <div class="max-w-3xl mx-auto space-y-4">
          <div v-for="(faq, index) in faqs" :key="index" class="border border-slate-200 rounded-2xl overflow-hidden">
            <button @click="toggleFaq(index)" class="w-full px-6 py-4 flex items-center justify-between bg-white hover:bg-slate-50 transition text-left">
              <span class="font-bold text-slate-900">{{ faq.question }}</span>
              <svg class="w-5 h-5 text-slate-400 transform transition-transform duration-200" :class="{'rotate-180': faq.isOpen}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
              </svg>
            </button>
            <div v-show="faq.isOpen" class="px-6 pb-4 bg-white text-slate-600 leading-relaxed">
              {{ faq.answer }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-24 bg-slate-50 border-t border-slate-100">
      <div class="container mx-auto px-6">
        <div class="bg-gradient-to-br from-blue-900 to-slate-900 rounded-[2.5rem] p-12 md:p-20 text-center relative overflow-hidden shadow-2xl shadow-blue-900/20">
          <!-- Decorative circles -->
          <div class="absolute top-0 left-0 w-96 h-96 bg-blue-600 rounded-full mix-blend-overlay filter blur-[100px] opacity-40 -translate-x-1/2 -translate-y-1/2"></div>
          <div class="absolute bottom-0 right-0 w-96 h-96 bg-sky-500 rounded-full mix-blend-overlay filter blur-[100px] opacity-40 translate-x-1/2 translate-y-1/2"></div>
          
          <div class="relative z-10 max-w-2xl mx-auto">
            <h2 class="text-3xl md:text-5xl font-bold text-white mb-6">Siap Mengembangkan Bisnis?</h2>
            <p class="text-lg text-blue-100 mb-10 leading-relaxed">
              Bergabung dengan ribuan UMKM lain yang telah mengotomatisasi bisnis mereka dengan teknologi AI terbaru.
            </p>
            <div class="flex flex-col sm:flex-row justify-center gap-4">
              <NuxtLink to="/register" class="px-8 py-4 bg-white text-blue-900 rounded-full font-bold hover:bg-blue-50 transition shadow-lg transform hover:-translate-y-1">
                Mulai Trial Gratis
              </NuxtLink>
              <NuxtLink to="/login" class="px-8 py-4 bg-transparent border border-blue-400/30 text-white rounded-full font-bold hover:bg-white/10 transition backdrop-blur-sm">
                Hubungi Sales
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-white border-t border-slate-100 py-16">
      <div class="container mx-auto px-6">
        <div class="grid md:grid-cols-4 gap-12 mb-12">
          <div class="col-span-1 md:col-span-2">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">G</span>
              </div>
              <span class="text-xl font-bold text-slate-900">Gowa UMKM</span>
            </div>
            <p class="text-slate-500 text-sm leading-relaxed max-w-sm">
              Platform manajemen WhatsApp berbasis AI untuk membantu UMKM Indonesia tumbuh dan berkembang dengan teknologi yang mudah digunakan.
            </p>
          </div>
          <div>
            <h4 class="font-bold text-slate-900 mb-6">Produk</h4>
            <ul class="space-y-3 text-sm text-slate-600">
              <li><a href="#" class="hover:text-blue-600 transition">Fitur</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Harga</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Integrasi</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Changelog</a></li>
            </ul>
          </div>
          <div>
            <h4 class="font-bold text-slate-900 mb-6">Perusahaan</h4>
            <ul class="space-y-3 text-sm text-slate-600">
              <li><a href="#" class="hover:text-blue-600 transition">Tentang Kami</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Karir</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Blog</a></li>
              <li><a href="#" class="hover:text-blue-600 transition">Kontak</a></li>
            </ul>
          </div>
        </div>
        <div class="pt-8 border-t border-slate-100 flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-sm text-slate-500">&copy; 2025 Gowa UMKM. All rights reserved.</p>
          <div class="flex space-x-6">
            <a href="#" class="text-slate-400 hover:text-blue-600 transition">
              <span class="sr-only">Twitter</span>
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84"/></svg>
            </a>
            <a href="#" class="text-slate-400 hover:text-blue-600 transition">
              <span class="sr-only">GitHub</span>
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"/></svg>
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false
})

const faqs = ref([
  {
    question: 'Apakah saya perlu kartu kredit untuk mendaftar?',
    answer: 'Tidak, Anda bisa mencoba paket Starter secara gratis selamanya tanpa perlu memasukkan kartu kredit.',
    isOpen: false
  },
  {
    question: 'Apakah aman menggunakan Gowa UMKM?',
    answer: 'Sangat aman. Kami menggunakan enkripsi standar industri untuk melindungi data Anda dan pelanggan Anda. Kami tidak pernah menyimpan pesan pribadi Anda.',
    isOpen: false
  },
  {
    question: 'Bisakah saya upgrade paket kapan saja?',
    answer: 'Ya, Anda bisa upgrade atau downgrade paket kapan saja sesuai kebutuhan bisnis Anda. Perubahan akan berlaku efektif segera.',
    isOpen: false
  },
  {
    question: 'Apakah mendukung multi-device?',
    answer: 'Ya, paket Pro dan Business mendukung multi-device, memungkinkan tim CS Anda mengelola satu nomor WhatsApp dari beberapa perangkat sekaligus.',
    isOpen: false
  }
])

const toggleFaq = (index: number) => {
  faqs.value[index].isOpen = !faqs.value[index].isOpen
}
</script>
