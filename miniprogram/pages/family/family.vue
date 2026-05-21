<template>
  <view class="family-page">
    <view v-show="!familyId && !isRegister" class="empty-state">
      <text>你还没有家庭账户</text>
      <button class="create-btn" @click="isRegister = true">创建我的家庭</button>
    </view>

    <!-- Register -->
    <view v-show="isRegister" class="form-card">
      <text class="form-title">创建家庭账户</text>
      <input type="text" v-model="form.name" placeholder="家庭名称（可选）" class="input" />
      <input type="text" v-model="form.nickname" placeholder="你的昵称" class="input" />
      <input type="number" v-model="form.phone" placeholder="手机号" class="input" />
      <button class="submit-btn" @click="submitFamily">确认创建</button>
    </view>

    <!-- Family detail -->
    <view v-show="familyId && !isRegister">
      <view class="family-header">
        <text class="family-name">{{ familyName || '我的家庭' }}</text>
        <text class="member-count">{{ members.length }} 位成员</text>
      </view>
      <view v-for="m in members" :key="m.member_uuid" class="member-card">
        <view class="avatar">{{ (m.nickname || '成')[0] }}</view>
        <view class="member-info">
          <text class="member-name">{{ m.nickname || '成员' }}</text>
          <text class="member-role">{{ roleLabel(m.relation) }}</text>
        </view>
      </view>
      <button class="add-btn" @click="showAdd = true">+ 添加家庭成员</button>
    </view>

    <!-- Add member dialog -->
    <view v-show="showAdd" class="form-card">
      <text class="form-title">添加家庭成员</text>
      <input type="text" v-model="addForm.nickname" placeholder="昵称" class="input" />
      <picker :range="relations" @change="onRelationChange">
        <view class="input picker">{{ addForm.relation || '选择关系' }}</view>
      </picker>
      <input type="number" v-model="addForm.phone" placeholder="手机号" class="input" />
      <button class="submit-btn" @click="submitAddMember">确认添加</button>
    </view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return {
      familyId: '', familyName: '', members: [], isRegister: false,
      form: { name: '', nickname: '', phone: '' },
      showAdd: false,
      addForm: { nickname: '', relation: 'child', phone: '' },
      relations: ['配偶', '子女', '父母', '其他'],
      relMap: { '配偶': 'spouse', '子女': 'child', '父母': 'parent', '其他': 'other' }
    }
  },
  onLoad(options) {
    if (options?.action === 'register') this.isRegister = true
    this.loadFamily()
  },
  methods: {
    roleLabel(r) {
      const m = { self: '本人', spouse: '配偶', child: '子女', parent: '父母', other: '其他' }
      return m[r] || r
    },
    loadFamily() {
      const fid = uni.getStorageSync('family_id')
      if (fid) {
        this.familyId = fid
        api.getFamily(fid).then(res => {
          this.familyName = res.name
          this.members = res.members || []
        })
      }
    },
    submitFamily() {
      api.createFamily({
        name: this.form.name || '我的家庭',
        nickname: this.form.nickname,
        phone_hash: this.form.phone
      }).then(res => {
        uni.setStorageSync('family_id', res.family_id)
        uni.setStorageSync('member_id', res.member_uuid)
        this.familyId = res.family_id
        this.isRegister = false
        this.loadFamily()
      }).catch(err => {
        uni.showToast({ title: '创建失败: ' + (err.error || ''), icon: 'none' })
      })
    },
    onRelationChange(e) {
      this.addForm.relation = this.relations[e.detail.value]
    },
    submitAddMember() {
      api.addMember(this.familyId, {
        nickname: this.addForm.nickname,
        relation: this.relMap[this.addForm.relation] || 'other',
        phone_hash: this.addForm.phone
      }).then(() => {
        this.showAdd = false
        this.loadFamily()
      }).catch(err => {
        uni.showToast({ title: '添加失败', icon: 'none' })
      })
    }
  }
}
</script>

<style>
.family-page { padding: 20px 16px; }
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 80px 0; color: #999; }
.create-btn { margin-top: 20px; background: #2E75B6; color: #fff; border: none; border-radius: 8px;
  padding: 10px 30px; font-size: 15px; }
.form-card { background: #fff; padding: 24px 16px; border-radius: 10px; margin-bottom: 16px; }
.form-title { font-size: 18px; font-weight: bold; margin-bottom: 16px; display: block; }
.input { display: block; border: 1px solid #E5E5E5; border-radius: 8px; padding: 10px 12px;
  margin-bottom: 12px; font-size: 14px; width: 100%; box-sizing: border-box;
  background: #fff; color: #333; }
.picker { color: #999; }
.submit-btn { display: block; width: 100%; background: #2E75B6; color: #fff; border: none;
  border-radius: 8px; padding: 14px; font-size: 16px; margin-top: 12px; }
.family-header { background: linear-gradient(135deg, #2E75B6, #1A4F7E); padding: 24px;
  border-radius: 12px; margin-bottom: 16px; }
.family-name { color: #fff; font-size: 20px; font-weight: bold; display: block; }
.member-count { color: rgba(255,255,255,0.7); font-size: 13px; margin-top: 4px; display: block; }
.member-card { display: flex; align-items: center; background: #fff; padding: 14px;
  border-radius: 8px; margin-bottom: 10px; }
.avatar { width: 44px; height: 44px; border-radius: 50%; background: #E8F4FD;
  display: flex; align-items: center; justify-content: center; font-size: 18px; color: #2E75B6; margin-right: 12px; }
.member-info { display: flex; flex-direction: column; }
.member-name { font-size: 15px; color: #333; }
.member-role { font-size: 12px; color: #999; margin-top: 2px; }
.add-btn { background: #fff; color: #2E75B6; border: 1px dashed #2E75B6; border-radius: 8px;
  padding: 12px; font-size: 15px; margin-top: 8px; }
</style>
