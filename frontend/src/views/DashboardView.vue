<template>
  <el-container class="layout">
    <el-header class="header">
      <span class="logo">AutoPowerHub</span>
      <div class="user-area">
        <span class="username">{{ authStore.username }}</span>
        <el-button type="danger" plain size="small" @click="handleLogout">Logout</el-button>
      </div>
    </el-header>

    <el-main class="main">
      <div v-if="loading" class="center">
        <el-text>Loading devices…</el-text>
      </div>

      <el-empty v-else-if="devices.length === 0" description="No devices configured" />

      <el-row v-else :gutter="20">
        <el-col
          v-for="device in devices"
          :key="device.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card class="device-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="device-name">{{ device.name }}</span>
                <el-tag :type="device.enabled ? 'success' : 'info'" size="small">
                  {{ device.enabled ? 'Online' : 'Offline' }}
                </el-tag>
              </div>
            </template>

            <p class="device-mac">{{ device.mac }}</p>

            <div class="actions">
              <el-button
                type="primary"
                :loading="busy[device.id] === 'power'"
                :disabled="!device.enabled || !!busy[device.id]"
                @click="handlePower(device.id)"
              >
                Power
              </el-button>
              <el-button
                type="warning"
                :loading="busy[device.id] === 'test'"
                :disabled="!device.enabled || !!busy[device.id]"
                @click="handleTest(device.id)"
              >
                Test
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getDevices, powerDevice, testDevice, type Device } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const devices = ref<Device[]>([])
const loading = ref(false)
const busy = reactive<Record<number, string>>({})

onMounted(fetchDevices)

async function fetchDevices() {
  loading.value = true
  try {
    const res = await getDevices()
    devices.value = res.data.devices ?? []
  } catch {
    ElMessage.error('Failed to load devices')
  } finally {
    loading.value = false
  }
}

async function handlePower(id: number) {
  await runCommand(id, 'power', () => powerDevice(id))
}

async function handleTest(id: number) {
  await runCommand(id, 'test', () => testDevice(id))
}

async function runCommand(id: number, key: string, fn: () => Promise<unknown>) {
  busy[id] = key
  try {
    await fn()
    ElMessage.success('Command sent successfully')
  } catch (err: unknown) {
    const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
    ElMessage.error(msg ?? 'Command failed')
  } finally {
    delete busy[id]
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
  background: #f0f2f5;
}
.header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}
.logo {
  font-size: 18px;
  font-weight: 700;
  color: #409eff;
}
.user-area {
  display: flex;
  align-items: center;
  gap: 12px;
}
.username {
  font-size: 14px;
  color: #606266;
}
.main {
  padding: 24px;
}
.center {
  display: flex;
  justify-content: center;
  padding-top: 80px;
}
.device-card {
  margin-bottom: 20px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.device-name {
  font-weight: 600;
  font-size: 15px;
}
.device-mac {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
  margin: 0 0 16px;
}
.actions {
  display: flex;
  gap: 8px;
}
</style>
