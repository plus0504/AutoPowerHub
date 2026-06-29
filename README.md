# AutoPowerHub

透過 Raspberry Pi + ESP32-C3 遠端控制電腦電源按鈕的個人設備管理系統。

---

## 系統架構

```
Browser
   │
   │  HTTP
   ▼
Vue3 + Element Plus (SPA)
   │
   │  REST API / JWT
   ▼
Go + Gin (Raspberry Pi 3B)
   │
   │  BLE (tinygo-bluetooth)
   ▼
ESP32-C3 Super Mini
   │
   │  PWM
   ▼
Servo Motor → 電源按鈕
```

---

## 功能特色

- 遠端觸發筆電 / 伺服器電源按鈕
- 管理多台 ESP32 裝置
- JWT 驗證，所有 API 均受保護
- 操作紀錄 Audit Log（使用者、裝置、指令、結果、時間）
- Vue3 Dashboard，一鍵 Power / Test
- ESP32 自動輕度睡眠省電（`esp_pm`）

---

## 技術棧

| 層級      | 技術                                              |
| --------- | ------------------------------------------------- |
| Frontend  | Vue 3, TypeScript, Element Plus, Vite, Pinia      |
| Backend   | Go 1.24, Gin, GORM, SQLite, JWT, tinygo-bluetooth |
| Firmware  | Arduino, NimBLE-Arduino, ESP32Servo               |
| Hardware  | Raspberry Pi 3B, ESP32-C3 Super Mini, SG90 Servo  |

---

## 專案結構

```
AutoPowerHub/
├── backend/
│   ├── cmd/main.go          # 程式進入點
│   ├── api/handler/         # HTTP handlers
│   ├── router/              # 路由設定
│   ├── middleware/          # JWT middleware
│   ├── service/
│   │   ├── auth/            # 認證服務
│   │   ├── ble/             # BLE 管理員
│   │   └── device/          # 裝置控制服務
│   ├── repository/          # 資料庫操作
│   ├── models/              # GORM 資料模型
│   ├── config/              # 設定載入
│   ├── database/            # SQLite 初始化
│   ├── web/                 # 前端建置產物 (frontend build output)
│   └── config.yaml          # 設定檔
├── frontend/
│   └── src/
│       ├── views/           # LoginView, DashboardView
│       ├── api/             # Axios API client
│       ├── stores/          # Pinia stores
│       └── router/          # Vue Router
└── firmware/
    └── esp32_powerhub/
        ├── src/main.cpp     # ESP32 韌體主程式
        └── platformio.ini   # PlatformIO 設定
```

---

## 快速開始

### 前置需求

- Raspberry Pi 3B（或相容 ARM64 Linux 裝置）
- Go 1.24+
- Node.js 18+
- PlatformIO CLI（燒錄 ESP32 韌體）
- ESP32-C3 Super Mini + SG90 Servo

---

### 1. 韌體燒錄（ESP32）

接線：Servo 訊號線 → GPIO 4

```bash
cd firmware/esp32_powerhub
pio run --target upload
pio device monitor   # 確認 BLE 廣播正常
```

---

### 2. 後端

```bash
cd backend

# 編輯設定檔
cp config.yaml config.yaml.local
vim config.yaml

# 取得 ESP32 的 MAC 位址後，填入 devices 區段
# 修改 jwt.secret

# 本機執行（需要 Go 環境）
go run ./cmd/main.go

# 或直接執行預編譯的 ARM64 執行檔（部署於 Raspberry Pi）
./autopowerhub-arm64
```

---

### 3. 前端

```bash
cd frontend
npm install
npm run build          # 產生 dist/
cp -r dist/* ../backend/web/
```

開發模式（含 Hot Reload）：

```bash
npm run dev            # Vite dev server，預設 http://localhost:5173
```

---

### 4. 存取 Dashboard

開啟瀏覽器：`http://<Raspberry-Pi-IP>:<port>`

預設帳號：`admin` / 密碼：`admin`（**首次登入後請立即修改**）

---

## 設定檔

`backend/config.yaml`

```yaml
server:
  port: 8081

jwt:
  secret: CHANGE_ME_IN_PRODUCTION

sqlite:
  path: ./data.db

admin:
  username: admin
  password: admin

devices:
  - name: MyServer
    mac: "AA:BB:CC:DD:EE:FF"
    service_uuid: "12345678-1234-1234-1234-1234567890ab"
    characteristic_uuid: "abcdefab-1234-5678-1234-abcdefabcdef"
    enabled: true
```

---

## REST API

所有 `/api/*` 端點（除 `/api/login`）須在 Header 帶入 JWT：

```
Authorization: Bearer <token>
```

| Method | 路徑                    | 說明             |
| ------ | ----------------------- | ---------------- |
| POST   | `/api/login`            | 登入，回傳 JWT   |
| GET    | `/api/device`           | 取得裝置清單     |
| POST   | `/api/device/:id/power` | 觸發電源按鈕     |
| POST   | `/api/device/:id/test`  | 短按測試（不開機）|
| GET    | `/api/debug/device/:id/scan` | BLE 掃描偵錯 |

---

## BLE 通訊協定

ESP32 透過 NimBLE GATT 服務接收文字指令：

| 指令      | 行為                                     |
| --------- | ---------------------------------------- |
| `PRESS`   | Servo 按壓電源按鈕（500 ms）後回到原位  |
| `TEST`    | Servo 短按測試（150 ms），不觸發開機    |
| `PING`    | 回覆 `PONG`（連線健康確認）              |
| `REBOOT`  | ESP32 重新啟動                           |

**GATT 識別碼（預設）**

```
Service UUID:        12345678-1234-1234-1234-1234567890ab
Characteristic UUID: abcdefab-1234-5678-1234-abcdefabcdef
```

---

## 資料庫 Schema

**users**

| 欄位       | 型別    |
| ---------- | ------- |
| id         | INTEGER |
| username   | TEXT    |
| password   | TEXT    |

**devices**

| 欄位                | 型別    |
| ------------------- | ------- |
| id                  | INTEGER |
| name                | TEXT    |
| mac                 | TEXT    |
| service_uuid        | TEXT    |
| characteristic_uuid | TEXT    |
| enabled             | INTEGER |

**logs**

| 欄位       | 型別     |
| ---------- | -------- |
| id         | INTEGER  |
| username   | TEXT     |
| device     | TEXT     |
| command    | TEXT     |
| result     | TEXT     |
| created_at | DATETIME |

---

## 安全注意事項

- 部署前務必更換 `config.yaml` 中的 `jwt.secret`
- 預設帳號 `admin/admin` 首次登入後請立即修改密碼
- 生產環境建議搭配 Nginx 反向代理並啟用 HTTPS

---

## 未來規劃

- [ ] Device CRUD（Web UI 新增 / 刪除裝置）
- [ ] Scheduler（排程自動開關機）
- [ ] WebSocket 即時狀態推送
- [ ] Telegram / MQTT 通知
- [ ] Wake-on-LAN 整合
- [ ] OTA 韌體更新
- [ ] HTTPS / Nginx 反向代理
- [ ] Docker 一鍵部署
- [ ] 多使用者 / RBAC 權限
- [ ] Home Assistant 整合

---

## License

MIT
