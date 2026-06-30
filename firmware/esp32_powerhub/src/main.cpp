/**
 * AutoPowerHub – ESP32-C3 Super Mini
 *
 * BLE peripheral that accepts text commands on a writable characteristic
 * and drives a servo motor to press/release a power button.
 *
 * Libraries required (Arduino Library Manager):
 *   - NimBLE-Arduino  (h2zero)
 *   - ESP32Servo      (Kevin Harrington)
 *
 * Wiring:
 *   Servo signal → GPIO 4 (configurable via SERVO_PIN)
 *   Servo power  → 5 V rail, GND to common ground
 */

#include <NimBLEDevice.h>
#include <ESP32Servo.h>
#include "esp_pm.h"

// ── BLE identifiers ──────────────────────────────────────────────────────────
// These UUIDs must match the values stored in the devices table.
#define SERVICE_UUID        "12345678-1234-1234-1234-1234567890ab"
#define CHARACTERISTIC_UUID "abcdefab-1234-5678-1234-abcdefabcdef"

// ── Hardware ─────────────────────────────────────────────────────────────────
#define SERVO_PIN          4
#define SERVO_RELEASE_DEG  0    // resting angle (degrees)
#define SERVO_PRESS_DEG    35   // press angle   (degrees)
#define PRESS_HOLD_MS      500  // how long to hold the pressed position

// ── Globals ───────────────────────────────────────────────────────────────────
Servo servo;
NimBLECharacteristic* pChar = nullptr;

// ── Actions ───────────────────────────────────────────────────────────────────
void doPress() {
    servo.attach(SERVO_PIN, 1000, 2000);
    servo.write(SERVO_PRESS_DEG);
    delay(PRESS_HOLD_MS);
    servo.write(SERVO_RELEASE_DEG);
    delay(300);
    servo.detach();
    Serial.println("[PRESS] done");
}

void doTest() {
    servo.attach(SERVO_PIN, 1000, 2000);
    servo.write(SERVO_PRESS_DEG);
    delay(150);
    servo.write(SERVO_RELEASE_DEG);
    delay(300);
    servo.detach();
    Serial.println("[TEST] done");
}

// ── BLE callbacks ─────────────────────────────────────────────────────────────
class CommandCallbacks : public NimBLECharacteristicCallbacks {
    void onWrite(NimBLECharacteristic* pC) override {
        std::string raw = pC->getValue();
        String cmd(raw.c_str());
        cmd.trim();

        Serial.print("[CMD] ");
        Serial.println(cmd);

        if (cmd == "PRESS") {
            doPress();
        } else if (cmd == "TEST") {
            doTest();
        } else if (cmd == "PING") {
            pC->setValue("PONG");
            pC->notify();
            Serial.println("[PING] pong sent");
        } else if (cmd == "REBOOT") {
            Serial.println("[REBOOT] restarting...");
            delay(200);
            ESP.restart();
        } else {
            Serial.print("[WARN] unknown command: ");
            Serial.println(cmd);
        }
    }
};

// ── Setup ─────────────────────────────────────────────────────────────────────
void setup() {
    #ifdef DEBUG_SERIAL
    Serial.begin(115200);
    #endif
    Serial.println("[BOOT] AutoPowerHub starting...");

    // Servo initialisation – restrict to standard 1000–2000 µs range to prevent
    // startup jitter caused by out-of-range pulses from the wider default window.
    servo.attach(SERVO_PIN, 1000, 2000);
    servo.write(SERVO_RELEASE_DEG);
    delay(500);
    servo.detach(); // cut PWM when idle to avoid light-sleep induced jitter

     // 1. 初始化並設定自動電源管理
    esp_pm_config_esp32c3_t pm_config = {
        .max_freq_mhz = 80,        // 與 platformio.ini 設定一致
        .min_freq_mhz = 10,        // 閒置時允許降到 10 MHz
        .light_sleep_enable = true // 關鍵：開啟自動輕度睡眠
    };
    
    // 2. 啟動電源管理
    esp_err_t err = esp_pm_configure(&pm_config);
    if (err == ESP_OK) {
        Serial.println("電源管理配置成功！晶片將在閒置時自動降溫。");
    } else {
        Serial.println("電源管理配置失敗。");
    }

    // BLE initialisation.
    NimBLEDevice::init("PowerHub");
    NimBLEDevice::setPower(ESP_PWR_LVL_N0);

    NimBLEServer* pServer = NimBLEDevice::createServer();
    NimBLEService* pService = pServer->createService(SERVICE_UUID);

    pChar = pService->createCharacteristic(
        CHARACTERISTIC_UUID,
        NIMBLE_PROPERTY::WRITE |
        NIMBLE_PROPERTY::WRITE_NR |
        NIMBLE_PROPERTY::NOTIFY
    );
    pChar->setCallbacks(new CommandCallbacks());
    pChar->setValue("READY");

    pService->start();

    NimBLEAdvertising* pAdv = NimBLEDevice::getAdvertising();
    pAdv->addServiceUUID(SERVICE_UUID);
    pAdv->setScanResponse(true);
    pAdv->start();

    Serial.println("[BOOT] BLE advertising as 'PowerHub'");
}

// ── Loop ──────────────────────────────────────────────────────────────────────
void loop() {
    // Nothing to poll; all work is driven by BLE callbacks.
    delay(1000);
}
