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
      <!-- Power Hub section -->
      <p class="section-title">Power Hub</p>

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
                Press Power Button
              </el-button>              
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- Smart Switches section -->
      <el-divider />
      <p class="section-title">Smart Switches</p>

      <div v-if="switchLoading" class="center">
        <el-text>Loading switches…</el-text>
      </div>

      <el-empty v-else-if="switches.length === 0" description="No switches found in cozylife_devices.json" />

      <el-row v-else :gutter="20">
        <el-col
          v-for="sw in switches"
          :key="sw.ip"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card class="device-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="device-name">{{ sw.dmn }}</span>
                <el-tag type="success" size="small">Switch</el-tag>
              </div>
            </template>

            <p class="device-mac">{{ sw.ip }}</p>

            <div class="actions">
              <el-button
                type="success"
                :loading="switchBusy[sw.ip] === 'on'"
                :disabled="!!switchBusy[sw.ip]"
                @click="handleSwitchOn(sw.ip)"
              >
                On
              </el-button>
              <el-button
                type="danger"
                plain
                :loading="switchBusy[sw.ip] === 'off'"
                :disabled="!!switchBusy[sw.ip]"
                @click="handleSwitchOff(sw.ip)"
              >
                Off
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
import { getDevices, powerDevice, getSwitches, switchOn, switchOff, type Device, type CozySwitch } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const devices = ref<Device[]>([])
const loading = ref(false)
const busy = reactive<Record<number, string>>({})

const switches = ref<CozySwitch[]>([])
const switchLoading = ref(false)
const switchBusy = reactive<Record<string, string>>({})

onMounted(() => { fetchDevices(); fetchSwitches() })

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

async function fetchSwitches() {
  switchLoading.value = true
  try {
    const res = await getSwitches()
    switches.value = res.data.switches ?? []
  } catch {
    ElMessage.error('Failed to load switches')
  } finally {
    switchLoading.value = false
  }
}

async function handlePower(id: number) {
  await runCommand(id, 'power', () => powerDevice(id))
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

async function handleSwitchOn(ip: string) {
  await runSwitchCommand(ip, 'on', () => switchOn(ip))
}

async function handleSwitchOff(ip: string) {
  await runSwitchCommand(ip, 'off', () => switchOff(ip))
}

async function runSwitchCommand(ip: string, key: string, fn: () => Promise<unknown>) {
  switchBusy[ip] = key
  try {
    await fn()
    ElMessage.success('Command sent')
  } catch (err: unknown) {
    const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
    ElMessage.error(msg ?? 'Command failed')
  } finally {
    delete switchBusy[ip]
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
.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 16px;
}
</style>
