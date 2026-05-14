<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSpin } from 'naive-ui'
import { websiteApi } from '../api'

const cdnUrl = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await websiteApi.get()
    cdnUrl.value = res.data.data?.cdn_url || ''
  } catch { /* ignore */ }
  loading.value = false
})
</script>

<template>
  <div class="page-wrap about-page">
    <n-spin :show="loading">
      <div class="about-content">
        <div class="about-section">
          <h2>关于我</h2>
          <div class="qr-section">
            <img
              v-if="cdnUrl"
              :src="cdnUrl.replace(/\/+$/, '') + '/about_me.jpg'"
              alt="微信公众号"
              class="qr-image"
            />
            <div class="bio-text">
              <p>一个喜欢踢足球的游戏开发程序员。</p>
              <div class="social-links">
                <a href="https://github.com/iamzcr" target="_blank" class="social-link">
                  <svg class="social-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                  GitHub
                </a>
                <a href="http://blog.iamzcr.com/" target="_blank" class="social-link">
                  <svg class="social-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
                  </svg>
                  Blog
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </n-spin>
  </div>
</template>

<style scoped>
.about-page {
  display: flex;
  justify-content: center;
}

.about-content {
  max-width: 600px;
  width: 100%;
  text-align: center;
}

.about-section h2 {
  font-size: 22px;
  margin-bottom: 24px;
  color: var(--text);
}

.qr-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.qr-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid var(--border);
}

.bio-text {
  color: var(--text);
}

.bio-text p {
  font-size: 16px;
  line-height: 1.8;
  margin-bottom: 20px;
}

.social-links {
  display: flex;
  gap: 24px;
  justify-content: center;
}

.social-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text);
  text-decoration: none;
  font-size: 15px;
  padding: 8px 16px;
  border: 1px solid var(--border);
  border-radius: 6px;
  transition: all 0.2s;
}

.social-link:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.social-icon {
  width: 18px;
  height: 18px;
}
</style>
