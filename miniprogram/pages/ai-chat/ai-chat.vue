<template>
  <view class="chat-page">
    <scroll-view class="messages" scroll-y :scroll-top="scrollTop">
      <!-- 欢迎提示 -->
      <view v-if="messages.length === 0" class="welcome">
        <text class="welcome-title">惠福灵犀 · AI 健康助手</text>
        <text class="welcome-sub">我是您的专属健康顾问，可以帮您解答孕期、产后、育儿的各种问题</text>
        <view class="faq-grid">
          <view v-for="faq in faqList" :key="faq" class="faq-chip" @click="tapFAQ(faq)">{{ faq }}</view>
        </view>
      </view>

      <view v-for="(msg, i) in messages" :key="i"
        :class="['msg', msg.role === 'user' ? 'msg-user' : 'msg-ai']">
        <text>{{ msg.content }}</text>
        <text v-if="msg.emergency?.is_emergency" class="emergency-alert">
          紧急关键词已触发，我们的医疗管家将在10分钟内联系你
        </text>
      </view>
      <view v-if="thinking" class="msg msg-ai">
        <text class="thinking">思考中...</text>
      </view>
    </scroll-view>

    <view class="input-bar">
      <input v-model="input" placeholder="输入你的问题..." class="chat-input"
        confirm-type="send" @confirm="send" />
      <button class="send-btn" @click="send" :disabled="!input.trim()">发送</button>
    </view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return {
      messages: [],
      input: '',
      thinking: false,
      sessionId: '',
      scrollTop: 0,
      faqList: ['孕期可以运动吗？', 'NT检查是什么？', '产后多久可以洗澡？', '宝宝发烧怎么办？', '母乳不足怎么办？', '糖耐量要空腹吗？']
    }
  },
  onLoad() {
    this.sessionId = 'sess_' + Date.now()
  },
  methods: {
    tapFAQ(q) {
      this.input = q
      this.send()
    },
    async send() {
      const text = this.input.trim()
      if (!text) return
      this.messages.push({ role: 'user', content: text })
      this.input = ''
      this.thinking = true
      this.scrollToBottom()

      try {
        const res = await api.chat(text, this.sessionId)
        const msg = { role: 'ai', content: res.reply, source: res.source, emergency: res.emergency }
        this.messages.push(msg)
      } catch {
        this.messages.push({ role: 'ai', content: '抱歉，AI 服务暂未配置，请先设置 DEEPSEEK_API_KEY。您可以浏览 FAQ 获取常见问题解答。' })
      } finally {
        this.thinking = false
        this.scrollToBottom()
      }
    },
    scrollToBottom() {
      this.$nextTick(() => { this.scrollTop = 99999 })
    }
  }
}
</script>

<style>
.chat-page { display: flex; flex-direction: column; height: 100vh; background: #F5F7FA; }
.messages { flex: 1; padding: 12px 16px; overflow-y: auto; }
.msg { max-width: 80%; padding: 10px 14px; border-radius: 12px; margin-bottom: 12px;
  font-size: 14px; line-height: 1.5; }
.msg-user { align-self: flex-end; background: #2E75B6; color: #fff; margin-left: auto; }
.msg-ai { align-self: flex-start; background: #fff; color: #333; }
.thinking { color: #999; font-style: italic; }
.emergency-alert { color: #E53E3E; font-weight: bold; display: block; margin-top: 6px; font-size: 12px; }
.input-bar { display: flex; padding: 10px 12px; background: #fff; border-top: 1px solid #eee; }
.chat-input { flex: 1; background: #F5F5F5; border-radius: 20px; padding: 8px 16px;
  font-size: 14px; border: none; }
.send-btn { background: #2E75B6; color: #fff; border: none; border-radius: 20px;
  padding: 8px 20px; margin-left: 8px; font-size: 14px; }
.welcome { padding: 30px 16px; text-align: center; }
.welcome-title { font-size: 20px; font-weight: bold; color: #2E75B6; display: block; margin-bottom: 8px; }
.welcome-sub { font-size: 13px; color: #999; display: block; margin-bottom: 20px; }
.faq-grid { display: flex; flex-wrap: wrap; gap: 10px; justify-content: center; }
.faq-chip { background: #E8F4FD; color: #2E75B6; padding: 8px 16px; border-radius: 16px; font-size: 13px; }
</style>
