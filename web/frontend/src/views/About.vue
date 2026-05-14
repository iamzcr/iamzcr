<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { websiteApi } from '../api'

const cdnUrl = ref('')
const loaded = ref(false)

onMounted(async () => {
  try {
    const res = await websiteApi.get()
    cdnUrl.value = res.data.data?.cdn_url || ''
  } catch { /* ignore */ }
  loaded.value = true
})
</script>

<template>
  <div class="about-wrap">
    <transition name="fade-up">
      <div class="about-card" v-if="loaded">
        <div class="card-accent"></div>

        <header class="about-header">
          <div class="header-prefix">
            <span class="prefix-line"></span>
            <span class="prefix-text">ABOUT</span>
            <span class="prefix-line"></span>
          </div>
          <h1>关于我</h1>
        </header>

        <div class="profile-section">
          <div class="qr-frame" v-if="cdnUrl">
            <div class="qr-ring"></div>
            <img
              :src="cdnUrl.replace(/\/+$/, '') + '/about_me.jpg'"
              alt="微信公众号"
              class="qr-img"
            />
            <span class="qr-caption">关注公众号</span>
          </div>

          <div class="bio-section">
            <p class="bio-line">一个喜欢<span class="highlight">踢足球</span>的游戏开发程序员。</p>

            <div class="links-section">
              <a href="https://github.com/iamzcr" target="_blank" class="link-card">
                <span class="link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                </span>
                <span class="link-label">GitHub</span>
                <span class="link-hint">开源项目</span>
              </a>

              <a href="http://blog.iamzcr.com/" target="_blank" class="link-card">
                <span class="link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
                  </svg>
                </span>
                <span class="link-label">Blog</span>
                <span class="link-hint">博客首页</span>
              </a>
            </div>
          </div>
        </div>

        <footer class="about-footer">
          <div class="tech-stack">
            <span class="tech-badge">Go</span>
            <span class="tech-dot">·</span>
            <span class="tech-badge">Vue</span>
            <span class="tech-dot">·</span>
            <span class="tech-badge">Unity</span>
            <span class="tech-dot">·</span>
            <span class="tech-badge">Cocos</span>
          </div>
        </footer>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.about-wrap {
  display: flex;
  justify-content: center;
  padding: 60px 20px 80px;
}

.about-card {
  max-width: 580px;
  width: 100%;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 56px 48px 48px;
  position: relative;
  overflow: hidden;
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent), #4ecb71, var(--accent));
}

/* ── header ── */
.about-header {
  text-align: center;
  margin-bottom: 44px;
}

.header-prefix {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 12px;
}

.prefix-line {
  width: 24px;
  height: 1px;
  background: var(--border);
}

.prefix-text {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 3px;
  color: var(--text-muted);
  font-family: var(--mono);
}

.about-header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 800;
  color: var(--text-h);
  letter-spacing: 0.5px;
}

/* ── profile ── */
.profile-section {
  display: flex;
  gap: 40px;
  align-items: flex-start;
}

.qr-frame {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  position: relative;
  padding: 4px;
}

.qr-ring {
  position: absolute;
  inset: -6px;
  border-radius: 16px;
  border: 2px dashed var(--border);
  pointer-events: none;
}

.qr-img {
  width: 160px;
  height: 160px;
  border-radius: 12px;
  object-fit: cover;
  display: block;
}

.qr-caption {
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
  font-weight: 600;
}

/* ── bio ── */
.bio-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 28px;
  padding-top: 4px;
}

.bio-line {
  margin: 0;
  font-size: 16px;
  line-height: 1.9;
  color: var(--text);
}

.highlight {
  color: var(--accent);
  font-weight: 700;
}

/* ── links ── */
.links-section {
  display: flex;
  gap: 12px;
}

.link-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 12px;
  border-radius: 12px;
  border: 1px solid var(--border);
  text-decoration: none;
  color: var(--text);
  transition: all 0.2s;
  background: var(--bg);
}

.link-card:hover {
  border-color: var(--accent-border);
  background: var(--accent-bg);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.link-icon {
  color: var(--text);
  transition: color 0.2s;
}

.link-card:hover .link-icon {
  color: var(--accent);
}

.link-label {
  font-weight: 700;
  font-size: 14px;
  color: var(--text-h);
}

.link-hint {
  font-size: 11px;
  color: var(--text-muted);
}

/* ── footer ── */
.about-footer {
  margin-top: 40px;
  padding-top: 28px;
  border-top: 1px solid var(--border);
  text-align: center;
}

.tech-stack {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.tech-badge {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  font-family: var(--mono);
  letter-spacing: 0.5px;
}

.tech-dot {
  color: var(--border);
  font-size: 10px;
  margin: 0 2px;
}

/* ── animation ── */
.fade-up-enter-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.fade-up-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

/* ── responsive ── */
@media (max-width: 600px) {
  .about-card {
    padding: 40px 24px 36px;
  }

  .profile-section {
    flex-direction: column;
    align-items: center;
    gap: 28px;
  }

  .bio-section {
    align-items: center;
    text-align: center;
  }

  .links-section {
    width: 100%;
  }
}
</style>
