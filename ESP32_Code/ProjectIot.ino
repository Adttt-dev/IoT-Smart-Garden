#include <DHT.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>

// ===== WIFI & API CONFIG =====
const char* ssid = "haha";        
const char* password = "87654321"; 
String api_base_url = "http://192.168.39.89:8080/api";
String auth_token = "";
String device_id = "4";
bool api_connected = false;

// ===== PIN DEFINITIONS =====
#define DHTPIN 2       
#define DHTTYPE DHT11
#define SOIL_PIN 34    
#define POMPA_PIN 15   
#define TRIG_PIN 5     
#define ECHO_PIN 18    

// ===== TIMING CONSTANTS =====
const unsigned long SENSOR_READ_INTERVAL = 1000;   
const unsigned long DHT_READ_INTERVAL = 5000;      
const unsigned long API_SEND_INTERVAL = 30000;     
const unsigned long SIMPLE_DISPLAY_INTERVAL = 1000; 

// ===== HARDCODED SETTINGS (UPDATED FOR 15CM JAR) =====
struct Settings {
  // Calibration - UPDATED FOR 15CM JAR
  int soil_dry = 4095, soil_wet = 0;
  float tank_height = 15.0;  // CHANGED: 15cm jar height
  
  // Fuzzy parameters - 4 SOIL CATEGORIES
  float soil_dry_max = 40.0;
  float soil_medium_min = 35.0, soil_medium_max = 65.0;
  float soil_moist_min = 60.0, soil_moist_max = 85.0;  // BASAH: 60-85%
  float soil_very_wet_min = 80.0;  // SANGAT BASAH: >80%
  
  // Temperature parameters (unchanged)
  float temp_cold_max = 30.0;
  float temp_normal_min = 28.0, temp_normal_max = 35.0;
  float temp_hot_min = 33.0;
  
  // PWM values
  int pump_pwm_med = 127, pump_pwm_high = 204, pump_pwm_max = 255;
  
  // Safety & timing - UPDATED FOR SMALLER CONTAINER
  float min_water_level = 8.0;  // CHANGED: Minimum 8% (about 1.2cm) for 15cm jar
  int max_duration_min = 10, cooldown_min = 3;  // CHANGED: Shorter duration/cooldown for small jar
  
  // Control flags
  bool auto_mode = true, debug_mode = false;
  bool simple_display = true;
  bool enable_fuzzy_logic = true;
} settings;


struct PumpState {
  int pwm = 0, percent = 0;
  String status = "OFF";
  String internal_status = "OFF"; // Status lengkap untuk debugging
  String fuzzy_explanation = "";
  bool active = false;
  unsigned long start_time = 0, cooldown_start = 0;
  bool in_cooldown = false;
  float fuzzy_confidence = 0.0;
} pump;

// ===== GLOBAL VARIABLES =====
DHT dht(DHTPIN, DHTTYPE);
float temp = 25.0, humidity = 60.0, soil_percent = 0.0;
float water_level = 0.0, water_percent = 0.0;
int soil_raw = 0;
bool water_ok = true;
String system_status = "INIT";

// ===== TIMERS =====
unsigned long last_sensor = 0, last_dht = 0, last_api = 0;
unsigned long last_simple_display = 0;

// ===== API CREDENTIALS =====
const String api_username = "admin";
const String api_password = "adminiot123";  
const String api_email = "admin@gmail.com";

// ===== FUNGSI KONVERSI STATUS UNTUK DATABASE =====
String convertStatusForDatabase(String internal_status) {
  // Konversi status internal ke format yang diterima database
  // Database hanya menerima: 'OFF','MED','HIGH','MAX','NO_WATER'
  
  if (internal_status == "OFF" || 
      internal_status == "OFF_FUZZY" || 
      internal_status == "OFF_CLASSIC" ||
      internal_status == "NO_RULE" ||
      internal_status == "COOLDOWN" ||
      internal_status == "TIMEOUT") {
    return "OFF";
  }
  else if (internal_status == "MED_FUZZY" || 
           internal_status == "MED_CLASSIC") {
    return "MED";
  }
  else if (internal_status == "HIGH_FUZZY" || 
           internal_status == "HIGH_CLASSIC") {
    return "HIGH";
  }
  else if (internal_status == "MAX_FUZZY" || 
           internal_status == "MAX_CLASSIC") {
    return "MAX";
  }
  else if (internal_status == "NO_WATER") {
    return "NO_WATER";
  }
  else {
    // Default fallback
    return "OFF";
  }
}

void setup() {
  Serial.begin(115200);
  Serial.println("🌱 Smart Irrigation with Fuzzy Logic Starting...");
  Serial.println("📏 Container: 15cm Jar Configuration");
  
  // WiFi setup
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\n✅ WiFi Connected! IP: " + WiFi.localIP().toString());
  
  // Hardware setup
  dht.begin();
  ledcAttach(POMPA_PIN, 1000, 8);
  ledcWrite(POMPA_PIN, 0);
  
  pinMode(SOIL_PIN, INPUT);
  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);
  
  // API connection
  Serial.println("🔐 Connecting to API...");
  authenticateAPI();
  
  system_status = "READY";
  Serial.println("🚀 System Ready with Fuzzy Logic!");
  Serial.printf("📏 Tank Height: %.1fcm, Min Water Level: %.1f%% (%.1fcm)\n", 
                settings.tank_height, settings.min_water_level, 
                (settings.min_water_level / 100.0) * settings.tank_height);
  printSimpleHeader();
}

void loop() {
  unsigned long now = millis();
  
  // Periodic tasks
  if (now - last_sensor >= SENSOR_READ_INTERVAL) {
    readSensors();
    last_sensor = now;
  }
  
  if (now - last_api >= API_SEND_INTERVAL) {
    if (api_connected) {
      sendData();
    } else {
      authenticateAPI();
    }
    last_api = now;
  }
  
  if (settings.simple_display && now - last_simple_display >= SIMPLE_DISPLAY_INTERVAL) {
    printSimpleData();
    last_simple_display = now;
  }
  
  checkPumpCooldown();
  delay(100);
}

// ===== FUZZY LOGIC FUNCTIONS =====
// Trapezoidal membership function
float trapezoidalMembership(float x, float a, float b, float c, float d) {
  if (x <= a || x >= d) return 0.0;
  if (x >= b && x <= c) return 1.0;
  if (x > a && x < b) return (x - a) / (b - a);
  return (d - x) / (d - a);
}

// Triangular membership function
float triangularMembership(float x, float a, float b, float c) {
  if (x <= a || x >= c) return 0.0;
  if (x == b) return 1.0;
  if (x > a && x < b) return (x - a) / (b - a);
  return (c - x) / (c - b);
}

// Apply fuzzy rules
// ===== UPDATED FUZZY LOGIC FUNCTION =====
// ===== UPDATED FUZZY LOGIC FUNCTION WITH 4 SOIL CATEGORIES =====
void applyFuzzyRules(float soil, float temp) {
  if (!settings.enable_fuzzy_logic) {
    calculatePumpClassic();
    return;
  }
  
  // Calculate soil membership - 4 CATEGORIES
  float soil_kering = trapezoidalMembership(soil, 0, 0, 35, settings.soil_dry_max);
  float soil_sedang = triangularMembership(soil, settings.soil_medium_min, 50, settings.soil_medium_max);
  float soil_basah = trapezoidalMembership(soil, settings.soil_moist_min, 65, 80, settings.soil_moist_max);  // 60-85%
  float soil_sangat_basah = trapezoidalMembership(soil, settings.soil_very_wet_min, 85, 100, 100);  // >80%
  
  // Calculate temperature membership (unchanged)
  float temp_dingin = trapezoidalMembership(temp, 0, 0, 25, settings.temp_cold_max);
  float temp_normal = triangularMembership(temp, settings.temp_normal_min, 31.5, settings.temp_normal_max);
  float temp_panas = trapezoidalMembership(temp, settings.temp_hot_min, 38, 50, 50);
  
  // Initialize rule outputs
  float off_strength = 0.0, med_strength = 0.0, high_strength = 0.0, max_strength = 0.0;
  String explanation = "Rules: ";
  
  // Apply 12 fuzzy rules (3 temp × 4 soil = 12 rules)
  float rules[12] = {
    // KERING + Suhu
    min(soil_kering, temp_dingin),     // R1: Kering + Dingin -> HIGH
    min(soil_kering, temp_normal),     // R2: Kering + Normal -> MAX
    min(soil_kering, temp_panas),      // R3: Kering + Panas -> MAX
    
    // SEDANG + Suhu
    min(soil_sedang, temp_dingin),     // R4: Sedang + Dingin -> MED
    min(soil_sedang, temp_normal),     // R5: Sedang + Normal -> HIGH
    min(soil_sedang, temp_panas),      // R6: Sedang + Panas -> MAX
    
    // BASAH + Suhu
    min(soil_basah, temp_dingin),      // R7: Basah + Dingin -> OFF
    min(soil_basah, temp_normal),      // R8: Basah + Normal -> MED
    min(soil_basah, temp_panas),       // R9: Basah + Panas -> HIGH
    
    // SANGAT BASAH + Suhu (SELALU OFF)
    min(soil_sangat_basah, temp_dingin), // R10: Sangat Basah + Dingin -> OFF
    min(soil_sangat_basah, temp_normal),  // R11: Sangat Basah + Normal -> OFF
    min(soil_sangat_basah, temp_panas)    // R12: Sangat Basah + Panas -> OFF
  };
  
  // Apply rule outputs
  high_strength = max(high_strength, rules[0]);  // R1 -> HIGH
  max_strength = max(max_strength, max(rules[1], rules[2])); // R2,R3 -> MAX
  med_strength = max(med_strength, max(rules[3], rules[7])); // R4,R8 -> MED
  high_strength = max(high_strength, max(rules[4], rules[8])); // R5,R9 -> HIGH
  max_strength = max(max_strength, rules[5]); // R6 -> MAX
  off_strength = max(off_strength, rules[6]); // R7 -> OFF
  
  // SANGAT BASAH SELALU OFF (R10, R11, R12)
  off_strength = max(off_strength, max(rules[9], max(rules[10], rules[11])));
  
  // Build explanation
  for (int i = 0; i < 12; i++) {
    if (rules[i] > 0) {
      explanation += "R" + String(i+1) + "(" + String(rules[i], 2) + ") ";
    }
  }
  
  // Defuzzification using centroid method
  float total_strength = off_strength + med_strength + high_strength + max_strength;
  
  if (total_strength == 0) {
    pump.pwm = 0;
    pump.percent = 0;
    pump.internal_status = "NO_RULE";
    pump.status = convertStatusForDatabase("NO_RULE");
    pump.active = false;
    pump.fuzzy_confidence = 0.0;
  } else {
    // Weighted average (centroid)
    float weighted_pwm = (off_strength * 0 + 
                         med_strength * settings.pump_pwm_med + 
                         high_strength * settings.pump_pwm_high + 
                         max_strength * settings.pump_pwm_max) / total_strength;
    
    pump.pwm = (int)weighted_pwm;
    pump.percent = (pump.pwm * 100) / 255;
    pump.fuzzy_confidence = total_strength;
    
    // Determine dominant rule for status
    if (max_strength >= high_strength && max_strength >= med_strength && max_strength >= off_strength) {
      pump.internal_status = "MAX_FUZZY";
      pump.status = convertStatusForDatabase("MAX_FUZZY");
      pump.active = (max_strength > 0.1);
    } else if (high_strength >= med_strength && high_strength >= off_strength) {
      pump.internal_status = "HIGH_FUZZY";
      pump.status = convertStatusForDatabase("HIGH_FUZZY");
      pump.active = (high_strength > 0.1);
    } else if (med_strength >= off_strength) {
      pump.internal_status = "MED_FUZZY";
      pump.status = convertStatusForDatabase("MED_FUZZY");
      pump.active = (med_strength > 0.1);
    } else {
      pump.internal_status = "OFF_FUZZY";
      pump.status = convertStatusForDatabase("OFF_FUZZY");
      pump.active = false;
      pump.pwm = 0;
      pump.percent = 0;
    }
  }
  
  pump.fuzzy_explanation = explanation;
  
  if (settings.debug_mode) {
    Serial.println("🧠 FUZZY LOGIC (4 SOIL CATEGORIES):");
    Serial.printf("   Soil: Kering=%.2f, Sedang=%.2f, Basah=%.2f, SangatBasah=%.2f\n", 
                  soil_kering, soil_sedang, soil_basah, soil_sangat_basah);
    Serial.printf("   Temp: Dingin=%.2f, Normal=%.2f, Panas=%.2f\n", 
                  temp_dingin, temp_normal, temp_panas);
    Serial.printf("   Current: Soil=%.1f%%, Temp=%.1f°C\n", soil, temp);
    Serial.printf("   Output: OFF=%.2f, MED=%.2f, HIGH=%.2f, MAX=%.2f\n", 
                  off_strength, med_strength, high_strength, max_strength);
    Serial.printf("   Result: PWM=%d, Status=%s, Confidence=%.2f\n", 
                  pump.pwm, pump.internal_status.c_str(), pump.fuzzy_confidence);
  }
}

// Classic pump calculation (fallback)
void calculatePumpClassic() {
  if (soil_percent <= 70.0 && soil_percent >= 60.0 && temp >= 30.0) {
    pump.pwm = settings.pump_pwm_med;
    pump.percent = 50;
    pump.internal_status = "MED_CLASSIC";
    pump.status = convertStatusForDatabase("MED_CLASSIC"); // Convert to "MED"
    pump.active = true;
  } else if (soil_percent <= 61.0 && soil_percent >= 50.0 && temp >= 28.0) {
    pump.pwm = settings.pump_pwm_high;
    pump.percent = 80;
    pump.internal_status = "HIGH_CLASSIC";
    pump.status = convertStatusForDatabase("HIGH_CLASSIC"); // Convert to "HIGH"
    pump.active = true;
  } else if (soil_percent <= 51.0 && temp >= 25.0) {
    pump.pwm = settings.pump_pwm_max;
    pump.percent = 100;
    pump.internal_status = "MAX_CLASSIC";
    pump.status = convertStatusForDatabase("MAX_CLASSIC"); // Convert to "MAX"
    pump.active = true;
  } else {
    pump.pwm = 0; 
    pump.percent = 0; 
    pump.internal_status = "OFF_CLASSIC";
    pump.status = convertStatusForDatabase("OFF_CLASSIC"); // Convert to "OFF"
    pump.active = false;
  }
}

// ===== DISPLAY FUNCTIONS =====
void printSimpleHeader() {
  Serial.println("\n" + String('=').substring(0, 80));
  Serial.println("         SMART IRRIGATION - 15CM JAR CONFIGURATION");
  Serial.println(String('=').substring(0, 80));
  Serial.println("Time\t\tTemp\tHumid\tSoil\tWater\tPump\t\tStatus");
  Serial.println(String('-').substring(0, 80));
}

void printSimpleData() {
  unsigned long runtime_sec = millis() / 1000;
  int minutes = (runtime_sec / 60) % 60;
  int seconds = runtime_sec % 60;
  
  // Tampilkan status internal (lebih detail) di Serial Monitor
  Serial.printf("%02d:%02d\t\t%.1f°C\t%.1f%%\t%.1f%%\t%.1f%%\t%s\t\t%s\n",
                minutes, seconds, temp, humidity, soil_percent, water_percent,
                pump.internal_status.c_str(), system_status.c_str());
  
  if (!water_ok) Serial.println("⚠️  WARNING: Water level too low! (< 1.2cm)");
  if (pump.in_cooldown) {
    unsigned long remaining = (settings.cooldown_min * 60) - ((millis() - pump.cooldown_start) / 1000);
    Serial.printf("❄️  Cooldown: %lu seconds remaining\n", remaining);
  }
}

// ===== SENSOR FUNCTIONS =====
void readSensors() {
  // Soil moisture
  soil_raw = analogRead(SOIL_PIN);
  soil_percent = map(soil_raw, settings.soil_dry, settings.soil_wet, 0, 100);
  soil_percent = constrain(soil_percent, 0, 100);
  
  // Water level - OPTIMIZED FOR 15CM JAR
  float distance = measureDistance();
  if (distance > 0 && distance <= 50) { // Changed max range to 50cm for jar measurement
    water_level = settings.tank_height - distance;
    if (water_level < 0) water_level = 0;
    if (water_level > settings.tank_height) water_level = settings.tank_height; // Cap at jar height
    water_percent = (water_level / settings.tank_height) * 100.0;
    water_ok = water_percent >= settings.min_water_level;
    
    if (settings.debug_mode) {
      Serial.printf("🔍 Distance: %.1fcm, Water Level: %.1fcm, Water %%: %.1f%%\n", 
                    distance, water_level, water_percent);
    }
  } else if (settings.debug_mode) {
    Serial.printf("⚠️ Invalid distance reading: %.1fcm\n", distance);
  }
  
  // DHT with caching
  unsigned long now = millis();
  if (now - last_dht >= DHT_READ_INTERVAL) {
    float t = dht.readTemperature();
    float h = dht.readHumidity();
    if (!isnan(t) && !isnan(h) && t >= -40 && t <= 80 && h >= 0 && h <= 100) {
      temp = t;
      humidity = h;
    }
    last_dht = now;
  }
  
  // Pump control
  if (settings.auto_mode && water_ok && !pump.in_cooldown) {
    applyFuzzyRules(soil_percent, temp);
    updatePump();
  } else if (!water_ok) {
    pump.internal_status = "NO_WATER";
    pump.status = convertStatusForDatabase("NO_WATER"); // "NO_WATER"
    stopPump();
  } else if (pump.in_cooldown) {
    pump.internal_status = "COOLDOWN";
    pump.status = convertStatusForDatabase("COOLDOWN"); // Convert to "OFF"
    stopPump();
  }
  
  updateSystemStatus();
  handleSerialCommand();
}

float measureDistance() {
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);
  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);
  
  // Reduced timeout for shorter distances (jar measurement)
  unsigned long duration = pulseIn(ECHO_PIN, HIGH, 15000); // 15ms timeout for ~50cm max
  if (duration == 0) return -1; // Timeout occurred
  
  return (duration * 0.0343) / 2.0;
}

void updatePump() {
  // Check duration limit - SHORTER FOR SMALL JAR
  if (pump.active && pump.start_time > 0) {
    unsigned long duration = millis() - pump.start_time;
    if (duration > (settings.max_duration_min * 60000UL)) {
      pump.internal_status = "TIMEOUT";
      pump.status = convertStatusForDatabase("TIMEOUT"); // Convert to "OFF"
      stopPump();
      pump.in_cooldown = true;
      pump.cooldown_start = millis();
      return;
    }
  }
  
  // Start timing
  if (pump.active && pump.start_time == 0) {
    pump.start_time = millis();
  } else if (!pump.active) {
    pump.start_time = 0;
  }
  
  ledcWrite(POMPA_PIN, pump.pwm);
}

void stopPump() {
  pump.pwm = 0;
  pump.percent = 0;
  pump.active = false;
  pump.start_time = 0;
  ledcWrite(POMPA_PIN, 0);
}

void checkPumpCooldown() {
  if (pump.in_cooldown) {
    unsigned long duration = millis() - pump.cooldown_start;
    if (duration > (settings.cooldown_min * 60000UL)) {
      pump.in_cooldown = false;
      pump.cooldown_start = 0;
    }
  }
}

void updateSystemStatus() {
  if (!water_ok) system_status = "WATER_LOW";
  else if (pump.in_cooldown) system_status = "COOLDOWN";
  else if (pump.active) system_status = "PUMPING";
  else if (!settings.auto_mode) system_status = "MANUAL";
  else system_status = "MONITORING";
}

// ===== SERIAL COMMANDS =====
void handleSerialCommand() {
  if (Serial.available()) {
    String command = Serial.readString();
    command.trim();
    command.toUpperCase();
    
    if (command == "TOGGLE") {
      settings.simple_display = !settings.simple_display;
      Serial.println(settings.simple_display ? "✅ Simple display ON" : "❌ Simple display OFF");
      if (settings.simple_display) printSimpleHeader();
    }
    else if (command == "STATUS") {
      printDetailedStatus();
    }
    else if (command == "FUZZY") {
      settings.enable_fuzzy_logic = !settings.enable_fuzzy_logic;
      Serial.println(settings.enable_fuzzy_logic ? "🧠 Fuzzy Logic ON" : "📊 Classic Logic ON");
    }
    else if (command == "DEBUG") {
      settings.debug_mode = !settings.debug_mode;
      Serial.println(settings.debug_mode ? "🔍 Debug Mode ON" : "🔍 Debug Mode OFF");
    }
    else if (command == "HELP") {
      printHelp();
    }
    else if (command == "CALIBRATE") {
      calibrateUltrasonic();
    }
  }
}

// ===== NEW CALIBRATION FUNCTION FOR JAR =====
void calibrateUltrasonic() {
  Serial.println("\n🔧 === ULTRASONIC SENSOR CALIBRATION (15CM JAR) ===");
  Serial.println("📏 Make sure the jar is EMPTY, then press any key...");
  
  // Wait for user input
  while (!Serial.available()) {
    delay(100);
  }
  Serial.readString(); // Clear input
  
  // Measure empty jar
  float empty_distance = 0;
  for (int i = 0; i < 10; i++) {
    float dist = measureDistance();
    if (dist > 0) empty_distance += dist;
    delay(200);
  }
  empty_distance /= 10;
  
  Serial.printf("📊 Empty jar distance: %.2fcm\n", empty_distance);
  Serial.printf("📊 Calculated jar height: %.2fcm\n", empty_distance);
  
  Serial.println("💧 Now FILL the jar completely, then press any key...");
  while (!Serial.available()) {
    delay(100);
  }
  Serial.readString(); // Clear input
  
  // Measure full jar
  float full_distance = 0;
  for (int i = 0; i < 10; i++) {
    float dist = measureDistance();
    if (dist > 0) full_distance += dist;
    delay(200);
  }
  full_distance /= 10;
  
  float calculated_height = empty_distance - full_distance;
  
  Serial.printf("📊 Full jar distance: %.2fcm\n", full_distance);
  Serial.printf("📊 Calculated water height: %.2fcm\n", calculated_height);
  
  if (calculated_height > 10 && calculated_height < 20) {
    Serial.printf("✅ Calibration successful! Jar height: %.2fcm\n", calculated_height);
    Serial.println("💡 You can update tank_height in your code to: " + String(calculated_height, 1) + "cm");
  } else {
    Serial.println("⚠️  Calibration seems incorrect. Please check sensor placement.");
  }
  
  Serial.println("🔧 === CALIBRATION COMPLETE ===\n");
}

// ===== UPDATED DETAILED STATUS FUNCTION =====
void printDetailedStatus() {
  Serial.println("\n========== DETAILED STATUS (15CM JAR) ==========");
  Serial.printf("🌡️ Temperature: %.2f°C\n", temp);
  Serial.printf("💧 Humidity: %.2f%%\n", humidity);
  Serial.printf("🌱 Soil Moisture: %.2f%% (Raw: %d)\n", soil_percent, soil_raw);
  
  // Show soil categories
  if (soil_percent <= 40) {
    Serial.printf("   📊 Soil Category: KERING (0-40%%)\n");
  } else if (soil_percent <= 65) {
    Serial.printf("   📊 Soil Category: SEDANG (35-65%%)\n");
  } else if (soil_percent <= 85) {
    Serial.printf("   📊 Soil Category: BASAH (60-85%%)\n");
  } else {
    Serial.printf("   📊 Soil Category: SANGAT BASAH (>80%%) - POMPA MATI!\n");
  }
  
  Serial.printf("🪣 Water Level: %.2fcm (%.2f%%) - JAR HEIGHT: %.1fcm\n", 
                water_level, water_percent, settings.tank_height);
  Serial.printf("   💧 Min Level: %.1f%% (%.1fcm)\n", 
                settings.min_water_level, (settings.min_water_level/100.0)*settings.tank_height);
  Serial.printf("⚙️ Pump Internal: %s\n", pump.internal_status.c_str());
  Serial.printf("📤 Pump DB Status: %s - %d%% (PWM: %d)\n", pump.status.c_str(), pump.percent, pump.pwm);
  Serial.printf("🧠 Fuzzy Logic: %s\n", settings.enable_fuzzy_logic ? "ENABLED (4 Categories)" : "DISABLED");
  Serial.printf("🎯 Confidence: %.2f\n", pump.fuzzy_confidence);
  Serial.printf("⏰ Timing: Max Duration=%dmin, Cooldown=%dmin\n", 
                settings.max_duration_min, settings.cooldown_min);
  Serial.printf("🔄 System Status: %s\n", system_status.c_str());
  Serial.printf("📡 API Connected: %s\n", api_connected ? "YES" : "NO");
  Serial.printf("⏰ Uptime: %lu seconds\n", millis() / 1000);
  Serial.println("====================================================");
}

void printHelp() {
  Serial.println("\n========== SERIAL COMMANDS ==========");
  Serial.println("TOGGLE     - Toggle simple display on/off");
  Serial.println("STATUS     - Show detailed status");
  Serial.println("FUZZY      - Toggle fuzzy logic on/off");
  Serial.println("DEBUG      - Toggle debug mode on/off");
  Serial.println("CALIBRATE  - Run ultrasonic calibration for jar");
  Serial.println("HELP       - Show this help menu");
  Serial.println("====================================");
}

// ===== API FUNCTIONS =====
void authenticateAPI() {
  HTTPClient http;
  http.begin(api_base_url + "/auth/login");
  http.addHeader("Content-Type", "application/json");
  
  StaticJsonDocument<200> doc;
  doc["username"] = api_username;
  doc["password"] = api_password;
  doc["email"] = api_email;
  
  String payload;
  serializeJson(doc, payload);
  
  int code = http.POST(payload);
  if (code == 200) {
    String response = http.getString();
    StaticJsonDocument<300> resp;
    if (!deserializeJson(resp, response) && resp.containsKey("token")) {
      auth_token = resp["token"].as<String>();
      api_connected = true;
      Serial.println("✅ API Connected!");
    }
  } else {
    api_connected = false;
  }
  http.end();
}

void sendData() {
  HTTPClient http;
  http.begin(api_base_url + "/sensor-readings");
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Authorization", "Bearer " + auth_token);
  
  StaticJsonDocument<600> doc;
  doc["device_id"] = device_id.toInt();
  doc["temperature"] = round(temp * 100) / 100.0;
  doc["humidity"] = round(humidity * 100) / 100.0;
  doc["temperature_source"] = "sensor";
  doc["humidity_source"] = "sensor";
  
  doc["soil_moisture_raw"] = soil_raw;
  doc["soil_moisture_percent"] = round(soil_percent * 100) / 100.0;
  doc["water_level_cm"] = round(water_level * 100) / 100.0;
  doc["water_percentage"] = round(water_percent * 100) / 100.0;
  doc["tank_height_cm"] = settings.tank_height;
  
  // Gunakan status yang sudah dikonversi untuk database
  doc["pump_status"] = pump.status; // Sudah dalam format DB-compatible
  doc["pump_pwm_value"] = pump.pwm;
  doc["pump_percentage"] = pump.percent;
  doc["system_status"] = system_status;
  
  String explanation = "Mode=" + String(settings.enable_fuzzy_logic ? "FUZZY" : "CLASSIC") + 
                      ", Internal=" + pump.internal_status + 
                      ", Soil=" + String(soil_percent, 1) + "%, Temp=" + String(temp, 1) + "°C";
  
  if (settings.enable_fuzzy_logic && pump.fuzzy_confidence > 0) {
    explanation += ", Confidence=" + String(pump.fuzzy_confidence, 2);
    doc["fuzzy_confidence"] = pump.fuzzy_confidence;
    doc["fuzzy_explanation"] = pump.fuzzy_explanation;
  }
  
  doc["logic_explanation"] = explanation;
  doc["wifi_rssi"] = WiFi.RSSI();
  doc["uptime_ms"] = millis();
  
  String payload;
  serializeJson(doc, payload);
  
  if (settings.debug_mode) {
    Serial.println("📤 Sending to API:");
    Serial.println("   DB Status: " + pump.status);
    Serial.println("   Internal Status: " + pump.internal_status);
  }
  
  int code = http.POST(payload);
  if (code == 401) {
    api_connected = false;
    auth_token = "";
  } else if (code == 200) {
    if (settings.debug_mode) {
      Serial.println("✅ Data sent successfully!");
    }
  } else {
    if (settings.debug_mode) {
      Serial.println("❌ API Error: " + String(code));
    }
  }
  
  http.end();
}