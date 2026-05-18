<template>
  <view class="chat-page">
    <scroll-view class="messages" scroll-y :scroll-top="scrollTop">
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
      scrollTop: 0
    }
  },
  onLoad() {
    this.sessionId = 'sess_' + Date.now()
  },
  methods: {
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
        this.messages.push({ role: 'ai', content: '抱歉，服务暂时不可用，请稍后再试。' })
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
</style>
