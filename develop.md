# AutoPowerHub Development Specification

Version: v0.1.0

---

# Project Goal

AutoPowerHub 是一套部署於 Raspberry Pi 的個人設備控制系統。

主要用途：

* 遠端控制 Notebook Power Button
* 遠端控制 Server Power Button
* 管理多台 ESP32 控制器
* 提供 Web UI
* 提供安全登入機制
* 保留操作紀錄

---

# System Architecture

```
Browser
    │
HTTP/HTTPS
    │
Vue3 + Element Plus
    │
Gin REST API
    │
JWT Authentication
    │
Service Layer
    │
SQLite
    │
BLE Manager
    │
ESP32-C3
    │
Servo / Relay
```

---

# Technology Stack

## Backend

* Golang 1.24+
* Gin
* GORM
* SQLite
* JWT
* BLE

## Frontend

* Vue3
* TypeScript
* Element Plus
* Axios
* Vue Router

## Firmware

* ESP32-C3 Super Mini
* Arduino Framework
* NimBLE
* ESP32Servo

## Platform

* Raspberry Pi 3B
* Debian / Raspberry Pi OS

---

# Project Structure

```
autopowerhub/

backend/

    cmd/

    api/

    router/

    middleware/

    service/

        auth/

        ble/

    repository/

    database/

    models/

    config/

frontend/

firmware/

docs/

deploy/
```

---

# Milestone 1

## Authentication

* Login
* Logout
* JWT
* Password Hash (bcrypt)

---

## Device

Read Device List

Power Device

Test Device

---

## BLE

Scan Device

Connect

Send Command

Disconnect

---

## Audit Log

Record

* User
* Device
* Command
* Result
* Timestamp

---

# Database

## users

| Field    | Type    |
| -------- | ------- |
| id       | INTEGER |
| username | TEXT    |
| password | TEXT    |

---

## devices

| Field               | Type    |
| ------------------- | ------- |
| id                  | INTEGER |
| name                | TEXT    |
| mac                 | TEXT    |
| service_uuid        | TEXT    |
| characteristic_uuid | TEXT    |
| enabled             | INTEGER |

---

## logs

| Field      | Type     |
| ---------- | -------- |
| id         | INTEGER  |
| username   | TEXT     |
| device     | TEXT     |
| command    | TEXT     |
| result     | TEXT     |
| created_at | DATETIME |

---

# REST API

## Login

POST

```
/api/login
```

Response

```
JWT Token
```

---

## Device List

GET

```
/api/device
```

---

## Power

POST

```
/api/device/{id}/power
```

---

## Test

POST

```
/api/device/{id}/test
```

---

# BLE Protocol

Current Version

```
PRESS
```

Reserved Commands

```
PING

STATUS

TEST

CONFIG

REBOOT
```

Future

```
PRESS:30:300

ANGLE:20

SET:18:30:300
```

---

# ESP32 Behaviour

PRESS

```
Release Angle

↓

Press Angle

↓

Hold

↓

Release Angle
```

---

# Security

JWT Required

All APIs except

```
/api/login
```

must verify JWT.

Password stored using bcrypt.

No plaintext password.

---

# Configuration

config.yaml

```
server:
    port: 80

jwt:
    secret: CHANGE_ME

sqlite:
    path: ./data.db
```

---

# Initial Account

username

```
admin
```

password

```
admin
```

Password must be changed after first login.

---

# Frontend Pages

Login

Dashboard

No additional pages in MVP.

---

# Dashboard

Display

Device Name

Online Status

Buttons

```
Power

Test
```

---

# Development Rules

* Backend follows Layered Architecture.
* API must not access database directly.
* BLE communication must be encapsulated inside Service Layer.
* Repository handles database operations only.
* Configuration must not be hardcoded.
* All APIs return JSON.
* All code must be documented when necessary.

---

# Future Features

* Device CRUD
* Scheduler
* Telegram Notification
* MQTT
* Wake-on-LAN
* WebSocket
* Home Assistant Integration
* OTA Firmware Update
* HTTPS
* Multi-user
* RBAC Permission
* Docker Deployment
* Nginx Reverse Proxy

---

# Current Version Scope

Only implement:

* Login
* JWT
* SQLite
* Device List
* Power
* Test
* BLE Communication
* ESP32 Servo
* Audit Log

Everything else belongs to future milestones.
