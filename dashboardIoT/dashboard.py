import customtkinter as ctk
import requests
import threading
import time
from datetime import datetime
import math

class DashboardApp(ctk.CTk):
    """
    An ultra-modern dashboard for a Smart Garden device, featuring animated gauges,
    a glassmorphism design, and interactive elements.
    """
    def __init__(self, auth_token, user_data, api_endpoint, request_timeout, refresh_interval, device_id=4):
        super().__init__()
        
        # --- Core Parameters ---
        self.auth_token = auth_token
        self.user_data = user_data
        self.api_base_url = api_endpoint.rsplit('/api/', 1)[0]
        self.sensor_api_endpoint = api_endpoint
        self.request_timeout = request_timeout
        self.refresh_interval = refresh_interval
        self.device_id = device_id
        
        # --- State Variables ---
        self.auto_refresh_enabled = False
        self.device_info = {}
        
        # --- Role-based Access ---
        user_details = self.user_data.get('user', self.user_data)
        self.is_admin = user_details.get('role') == 'admin'
        self.logged_in_user_id = user_details.get('id')
        self.username = user_details.get('username', 'User')

        # --- Super Modern UI Configuration ---
        self.title("Smart Garden Dashboard")
        self.geometry("1200x850")
        
        # --- Colors & Fonts ---
        self.COLOR_PRIMARY = "#20BF55"
        self.COLOR_GRADIENT = "#0B4F6C"
        self.COLOR_BACKGROUND = "#171A21"
        self.COLOR_CARD_BG = "gray14"
        self.COLOR_CARD_BORDER = "gray20"
        self.COLOR_SECONDARY = "gray25" 
        self.COLOR_TEXT = "#FFFFFF"
        self.COLOR_TEXT_SECONDARY = "#A0A0A0"
        self.COLOR_MANUAL_MODE = "#FFC107" # Amber for manual mode

        self.FONT_TITLE = ("Roboto", 28, "bold")
        self.FONT_SUBTITLE = ("Roboto", 12)
        self.FONT_NORMAL = ("Roboto", 11)
        self.FONT_GAUGE = ("Roboto", 32, "bold")
        
        self.configure(fg_color=self.COLOR_BACKGROUND)

        # --- Initial Setup ---
        self.setup_ui()
        self.fetch_device_info()
        self.start_auto_refresh()

    def setup_ui(self):
        self.grid_columnconfigure(1, weight=1)
        self.grid_rowconfigure(0, weight=1)

        self.sidebar_frame = self.create_sidebar()
        self.sidebar_frame.grid(row=0, column=0, sticky="nsw", padx=10, pady=10)

        self.main_content_frame = ctk.CTkFrame(self, fg_color="transparent")
        self.main_content_frame.grid(row=0, column=1, sticky="nsew", padx=10, pady=10)
        self.main_content_frame.grid_columnconfigure(0, weight=1)
        self.main_content_frame.grid_rowconfigure(1, weight=1) # Sensor display
        self.main_content_frame.grid_rowconfigure(2, weight=0) # Control and Log row

        self.create_header()
        self.create_sensor_display()
        
        self.bottom_frame = ctk.CTkFrame(self.main_content_frame, fg_color="transparent")
        self.bottom_frame.grid(row=2, column=0, sticky="nsew", pady=(10,0))
        self.bottom_frame.grid_columnconfigure(0, weight=1) # Manual Control
        self.bottom_frame.grid_columnconfigure(1, weight=1) # Event Log

        self.create_manual_control_panel()
        self.create_log()

    def create_sidebar(self):
        sidebar = ctk.CTkFrame(self, width=220, fg_color=self.COLOR_CARD_BG, corner_radius=10, border_width=1, border_color=self.COLOR_CARD_BORDER)
        sidebar.grid_propagate(False)

        ctk.CTkLabel(sidebar, text="🌿 SmartGarden", font=("Roboto", 22, "bold"), text_color=self.COLOR_PRIMARY).pack(pady=(25, 30))

        ctk.CTkButton(sidebar, text="  Refresh Data", anchor="w", font=self.FONT_NORMAL, command=self.manual_refresh, image=self._get_icon("🔄")).pack(fill="x", padx=15, pady=6)

        ctk.CTkLabel(sidebar, text="SETTINGS", font=("Roboto", 10, "bold"), text_color=self.COLOR_TEXT_SECONDARY).pack(fill="x", padx=20, pady=(20, 5))
        ctk.CTkButton(sidebar, text="  Device Settings", anchor="w", font=self.FONT_NORMAL, command=lambda: self.fetch_device_info(callback=self.show_device_info_and_edit_dialog), image=self._get_icon("⚙️")).pack(fill="x", padx=15, pady=6)
        
        if self.is_admin:
            ctk.CTkButton(sidebar, text="  Manage Users", anchor="w", font=self.FONT_NORMAL, command=self.show_users_window, image=self._get_icon("👥")).pack(fill="x", padx=15, pady=6)
        
        ctk.CTkButton(sidebar, text="  Logout", anchor="w", fg_color="#982D2D", hover_color="#C62828", command=self.logout, image=self._get_icon("🚪")).pack(fill="x", padx=15, pady=20, side="bottom")

        self.auto_var = ctk.BooleanVar(value=True)
        ctk.CTkCheckBox(sidebar, text="Auto Refresh", variable=self.auto_var, command=self.toggle_auto).pack(fill="x", padx=15, pady=10, side="bottom")
        
        return sidebar

    def create_header(self):
        header_frame = ctk.CTkFrame(self.main_content_frame, fg_color="transparent")
        header_frame.grid(row=0, column=0, sticky="ew", pady=(0, 20))
        header_frame.grid_columnconfigure(0, weight=1)

        greeting = self.get_greeting()
        ctk.CTkLabel(header_frame, text=f"{greeting}, {self.username}!", font=self.FONT_TITLE, text_color=self.COLOR_TEXT).grid(row=0, column=0, sticky="w")
        
        title_status_frame = ctk.CTkFrame(header_frame, fg_color="transparent")
        title_status_frame.grid(row=1, column=0, sticky="w")

        device_name = self.device_info.get('device_name', 'Loading...')
        self.title_label = ctk.CTkLabel(title_status_frame, text=device_name, font=self.FONT_SUBTITLE, text_color=self.COLOR_TEXT_SECONDARY)
        self.title_label.pack(side="left", anchor="w")

        self.auto_mode_label = ctk.CTkLabel(title_status_frame, text="⚫ LOADING MODE", font=("Roboto", 10, "bold"), text_color="gray")
        self.auto_mode_label.pack(side="left", anchor="w", padx=10)
        
        date_frame = ctk.CTkFrame(header_frame, fg_color="transparent")
        date_frame.grid(row=0, column=1, rowspan=2, sticky="e")
        ctk.CTkLabel(date_frame, text=self.get_weather_icon(), font=("Arial", 32)).pack(side="left", padx=10)
        ctk.CTkLabel(date_frame, text=datetime.now().strftime("%A\n%d %B %Y"), font=self.FONT_NORMAL, text_color=self.COLOR_TEXT_SECONDARY, justify="right").pack(side="left")

    def create_sensor_display(self):
        self.sensor_frame = ctk.CTkFrame(self.main_content_frame, fg_color="transparent")
        self.sensor_frame.grid(row=1, column=0, sticky="nsew")
        self.sensor_frame.grid_columnconfigure((0, 1, 2), weight=1)
        self.sensor_frame.grid_rowconfigure((0, 1), weight=1)
        
        sensors = [
            {"icon": "🌡️", "name": "Suhu", "keys": ["temperature"], "unit": "°C", "color": ["#29B6F6", "#0288D1"], "gauge": True},
            {"icon": "💧", "name": "Kelembapan", "keys": ["humidity"], "unit": "%", "color": ["#66BB6A", "#388E3C"], "gauge": True},
            {"icon": "🌱", "name": "Kelembapan Tanah", "keys": ["soil_moisture_percent", "soil_moisture"], "unit": "%", "color": ["#8D6E63", "#5D4037"], "gauge": True},
            {"icon": "🚰", "name": "Level Air", "keys": ["water_percentage", "water_level"], "unit": "%", "color": ["#42A5F5", "#1976D2"], "gauge": True},
            {"icon": "⚡", "name": "Status Pompa", "keys": ["pump_status"], "unit": "", "color": ["#FFA726", "#F57C00"]},
            {"icon": "🔧", "name": "Status Sistem", "keys": ["system_status"], "unit": "", "color": ["#AB47BC", "#8E24AA"]},
        ]
        
        self.sensor_cards = {}
        for i, config in enumerate(sensors):
            row, col = divmod(i, 3)
            primary_key = config['keys'][0]
            card = self._create_sensor_card(self.sensor_frame, config)
            card['frame'].grid(row=row, column=col, padx=10, pady=10, sticky="nsew")
            self.sensor_cards[primary_key] = card

    def _create_sensor_card(self, parent, config):
        card_frame = ctk.CTkFrame(parent, fg_color=self.COLOR_CARD_BG, corner_radius=15, border_width=1, border_color=self.COLOR_CARD_BORDER)
        card_frame.grid_rowconfigure(1, weight=1)
        card_frame.grid_columnconfigure(0, weight=1)
        
        ctk.CTkLabel(card_frame, text=config['name'], font=("Roboto", 14), text_color=self.COLOR_TEXT_SECONDARY).grid(row=0, column=0, pady=(15, 5), padx=15, sticky="w")
        
        value_frame = ctk.CTkFrame(card_frame, fg_color="transparent")
        value_frame.grid(row=1, column=0, sticky="nsew", pady=10)

        card_widgets = {"frame": card_frame, "keys": config['keys'], "unit": config['unit'], "value_label": None, "canvas": None, "colors": config['color']}

        if config.get("gauge"):
            canvas = ctk.CTkCanvas(value_frame, width=150, height=150, bg=self.COLOR_CARD_BG, highlightthickness=0)
            canvas.pack()
            value_label = ctk.CTkLabel(value_frame, text="--", font=self.FONT_GAUGE, text_color=self.COLOR_TEXT)
            card_widgets.update({"canvas": canvas, "value_label": value_label})
        else:
            value_label = ctk.CTkLabel(value_frame, text="--", font=self.FONT_GAUGE, text_color=self.COLOR_TEXT)
            value_label.pack(expand=True)
            card_widgets['value_label'] = value_label

        original_border_color = self.COLOR_CARD_BORDER
        def on_enter(e): card_frame.configure(border_color=config['color'][0])
        def on_leave(e): card_frame.configure(border_color=original_border_color)
            
        card_frame.bind("<Enter>", on_enter)
        card_frame.bind("<Leave>", on_leave)

        return card_widgets

    def _update_gauge(self, canvas, value_label, value, unit, colors, secondary_color, text_secondary_color, normal_font):
        canvas.delete("all")
        w, h = 150, 150
        start_angle, full_angle = 140, 260
        
        try:
            percent = min(max(float(value) / 100.0, 0.0), 1.0)
            display_value = f"{float(value):.1f}"
        except (ValueError, TypeError):
            percent = 0.0
            display_value = "--"
            
        value_label.configure(text=display_value)
        value_label.place(relx=0.5, rely=0.5, anchor="center")

        canvas.create_text(w/2, h * 0.68, text=unit, font=normal_font, fill=text_secondary_color)
        canvas.create_arc(10, 10, w-10, h-10, start=start_angle, extent=full_angle, style="arc", outline=secondary_color, width=12, tags="bg")

        if percent > 0:
            canvas.create_arc(10, 10, w-10, h-10, start=start_angle, extent=full_angle * percent, style="arc", outline=colors[0], width=12, tags="fg")

    def create_manual_control_panel(self):
        self.control_frame = ctk.CTkFrame(self.bottom_frame, fg_color=self.COLOR_CARD_BG, corner_radius=10, border_width=1, border_color=self.COLOR_CARD_BORDER)
        self.control_frame.grid(row=0, column=0, sticky="nsew", padx=(0, 10))
        self.control_frame.grid_columnconfigure((0, 1, 2), weight=1)

        ctk.CTkLabel(self.control_frame, text="🕹️ Manual Control", font=self.FONT_SUBTITLE, text_color=self.COLOR_TEXT_SECONDARY).grid(row=0, column=0, columnspan=3, sticky="w", padx=15, pady=(10, 5))
        
        self.pump_on_btn = ctk.CTkButton(self.control_frame, text="Pump ON", command=lambda: self.send_device_command("PUMP_ON"), fg_color="#4CAF50", hover_color="#66BB6A", height=40)
        self.pump_on_btn.grid(row=1, column=0, padx=10, pady=10, sticky="ew")

        self.pump_off_btn = ctk.CTkButton(self.control_frame, text="Pump OFF", command=lambda: self.send_device_command("PUMP_OFF"), fg_color="#F44336", hover_color="#E57373", height=40)
        self.pump_off_btn.grid(row=1, column=1, padx=10, pady=10, sticky="ew")

        self.auto_mode_btn = ctk.CTkButton(self.control_frame, text="Set AUTO", command=lambda: self.send_device_command("AUTO_ON"), fg_color="#2196F3", hover_color="#64B5F6", height=40)
        self.auto_mode_btn.grid(row=1, column=2, padx=10, pady=10, sticky="ew")

    def create_log(self):
        log_container = ctk.CTkFrame(self.bottom_frame, fg_color="transparent")
        log_container.grid(row=0, column=1, sticky="nsew", padx=(10, 0))
        log_container.grid_columnconfigure(0, weight=1)
        
        ctk.CTkLabel(log_container, text="Event Log", font=self.FONT_SUBTITLE, text_color=self.COLOR_TEXT_SECONDARY).pack(anchor="w")
        self.log_text = ctk.CTkTextbox(log_container, height=120, fg_color=self.COLOR_CARD_BG, corner_radius=10, border_width=0, font=("Consolas", 10))
        self.log_text.pack(fill="both", expand=True, pady=5)
        self.log("System", "Dashboard UI Initialized.")

    def update_display(self, data):
        self.log("UI_DEBUG", f"Updating display with data: {str(data)[:200]}")
        for card in self.sensor_cards.values():
            if not card['frame'].winfo_exists(): continue

            value = '--'
            for key in card['keys']:
                if key in data:
                    value = data[key]
                    break
            
            try:
                if card.get("canvas"):
                    self._update_gauge(
                        canvas=card["canvas"], value_label=card["value_label"], value=value, 
                        unit=card["unit"], colors=card["colors"], secondary_color=self.COLOR_SECONDARY,
                        text_secondary_color=self.COLOR_TEXT_SECONDARY, normal_font=self.FONT_NORMAL
                    )
                elif card.get("value_label"):
                    display_text = "--" if value in [None, '--'] else str(value).upper()
                    card['value_label'].configure(text=display_text)
            except Exception as e:
                self.log("UI_ERROR", f"Failed to update card for keys {card['keys']}: {e}")

    def get_greeting(self):
        hour = datetime.now().hour
        if 5 <= hour < 12: return "Selamat Pagi"
        if 12 <= hour < 18: return "Selamat Siang"
        return "Selamat Malam"

    def get_weather_icon(self):
        hour = datetime.now().hour
        return "☀️" if 6 <= hour < 18 else "🌙"

    def _get_icon(self, emoji):
        return None

    def fetch_device_info(self, callback=None):
        def worker():
            try:
                headers = {'Authorization': f'Bearer {self.auth_token}'}
                device_endpoint = f"{self.api_base_url}/api/devices/{self.device_id}"
                response = requests.get(device_endpoint, headers=headers, timeout=self.request_timeout)
                if response.status_code == 200:
                    device_data = response.json().get('device', response.json())
                    if device_data:
                        self.device_info = device_data
                        self.after(0, self.update_device_ui, self.device_info)
                        if callback: self.after(0, callback, self.device_info)
                else:
                    self.after(0, self.log, "API", f"Error fetching info: {response.status_code}")
                    if callback: self.after(0, callback, None)
            except Exception as e:
                self.after(0, self.log, "API", f"Exception: {e}")
                if callback: self.after(0, callback, None)
        threading.Thread(target=worker, daemon=True).start()

    def send_device_command(self, command):
        self.log("COMMAND", f"Sending command: {command}")
        
        # Temporarily disable all buttons to prevent spam clicks
        self.pump_on_btn.configure(state="disabled")
        self.pump_off_btn.configure(state="disabled")
        self.auto_mode_btn.configure(state="disabled")

        def worker():
            try:
                headers = {'Authorization': f'Bearer {self.auth_token}', 'Content-Type': 'application/json'}
                command_endpoint = f"{self.api_base_url}/api/devices/{self.device_id}/command"
                payload = {"command": command}
                response = requests.put(command_endpoint, json=payload, headers=headers, timeout=self.request_timeout)

                if response.status_code == 200:
                    self.after(0, self.log, "COMMAND", f"Successfully sent '{command}' command.")
                    # Refresh device info immediately to get the latest state
                    self.after(100, self.fetch_device_info) 
                else:
                    error_msg = response.json().get('error', 'Unknown error')
                    self.after(0, self.log, "COMMAND_ERROR", f"Failed to send command: {error_msg} ({response.status_code})")
            except requests.exceptions.RequestException as e:
                self.after(0, self.log, "COMMAND_ERROR", f"Connection error: {e}")
            finally:
                # Re-enable all buttons after the operation is complete
                self.after(500, self.update_control_buttons_state)
        
        threading.Thread(target=worker, daemon=True).start()

    def update_device_ui(self, device_info):
        new_name = device_info.get('device_name', 'Unnamed Device')
        if hasattr(self, 'title_label') and self.title_label.winfo_exists():
            self.title_label.configure(text=new_name)
        self.title(f"{new_name} - Dashboard")
        
        is_auto = device_info.get('auto_mode', False)
        self.update_auto_mode_status(is_auto)
        self.update_control_buttons_state()

    def update_auto_mode_status(self, is_auto):
        if is_auto:
            self.auto_mode_label.configure(text="🟢 AUTO MODE", text_color=self.COLOR_PRIMARY)
        else:
            self.auto_mode_label.configure(text="🟠 MANUAL MODE", text_color=self.COLOR_MANUAL_MODE)
    
    # === FUNCTION WITH THE FIX ===
    def update_control_buttons_state(self):
        """
        Enables control buttons. In this corrected version, all buttons are always
        enabled to allow the user to switch modes at any time.
        """
        self.pump_on_btn.configure(state="normal")
        self.pump_off_btn.configure(state="normal")
        self.auto_mode_btn.configure(state="normal")
            
    def show_device_info_and_edit_dialog(self, device_info):
        if not device_info: self.log("UI", "Cannot open settings: data missing."); return
        
        dialog = ctk.CTkToplevel(self, fg_color=self.COLOR_BACKGROUND); 
        dialog.title("Pengaturan Perangkat" if self.is_admin else "Informasi Perangkat")
        dialog.geometry("450x500"); dialog.transient(self); dialog.grab_set()

        ctk.CTkLabel(dialog, text="Device Settings", font=("Roboto", 20, "bold"), text_color=self.COLOR_TEXT).pack(pady=(20,10))
        
        form_frame = ctk.CTkFrame(dialog, fg_color=self.COLOR_CARD_BG, corner_radius=10)
        form_frame.pack(fill="x", expand=True, padx=20, pady=10)

        # Removed 'ip_address' from the fields dictionary
        fields = {"device_name": "Nama Perangkat", "device_type": "Tipe", "location": "Lokasi"}
        entry_widgets = {}
        for i, (key, text) in enumerate(fields.items()):
            ctk.CTkLabel(form_frame, text=text, font=self.FONT_NORMAL, text_color=self.COLOR_TEXT_SECONDARY).pack(anchor="w", padx=20, pady=(15, 2))
            entry = ctk.CTkEntry(form_frame, placeholder_text=f"Masukkan {text}...", height=35, fg_color=self.COLOR_SECONDARY, border_width=0)
            entry.pack(fill="x", padx=20)
            entry.insert(0, str(device_info.get(key, '')))
            entry_widgets[key] = entry
            if not self.is_admin: entry.configure(state="disabled")

        status_label = ctk.CTkLabel(dialog, text=""); status_label.pack(pady=5)
        
        btn_frame = ctk.CTkFrame(dialog, fg_color="transparent"); btn_frame.pack(pady=10, padx=20, fill="x")
        btn_frame.grid_columnconfigure((0, 1), weight=1)
        if self.is_admin:
            ctk.CTkButton(btn_frame, text="Tutup", command=dialog.destroy, fg_color=self.COLOR_SECONDARY, hover_color="gray25").grid(row=0, column=0, padx=(0,5), sticky="ew")
            ctk.CTkButton(btn_frame, text="Simpan", command=lambda: self.save_device_info(entry_widgets, dialog, status_label), fg_color=self.COLOR_PRIMARY).grid(row=0, column=1, padx=(5,0), sticky="ew")
        else:
            ctk.CTkButton(btn_frame, text="Tutup", command=dialog.destroy).grid(row=0, column=0, columnspan=2, sticky="ew")

    def save_device_info(self, entries, dialog, status_label):
        if not self.is_admin: return
        new_data = {k: v.get().strip() for k, v in entries.items()}
        if not new_data.get("device_name"):
            status_label.configure(text="Nama perangkat tidak boleh kosong!", text_color="red"); return
        status_label.configure(text="Menyimpan...", text_color="cyan")
        threading.Thread(target=self._update_device_info_worker, args=(new_data, dialog, status_label), daemon=True).start()

    def _update_device_info_worker(self, update_data, dialog, status_label):
        try:
            headers = {'Authorization': f'Bearer {self.auth_token}', 'Content-Type': 'application/json'}
            response = requests.put(f"{self.api_base_url}/api/devices/{self.device_id}/", json=update_data, headers=headers, timeout=self.request_timeout)
            if response.status_code < 300:
                self.device_info.update(update_data)
                self.after(0, self.update_device_ui, self.device_info)
                self.after(0, status_label.configure, {"text": "Berhasil disimpan!", "text_color": "lightgreen"})
                self.after(1500, dialog.destroy)
            else:
                self.after(0, status_label.configure, {"text": f"Gagal: {response.json().get('error', response.status_code)}", "text_color": "red"})
        except Exception as e:
            self.after(0, status_label.configure, {"text": f"Error: {e}", "text_color": "red"})

    def show_users_window(self):
        if not self.is_admin: return
        win = ctk.CTkToplevel(self, fg_color=self.COLOR_BACKGROUND); win.title("Kelola Pengguna"); win.geometry("600x500"); win.transient(self); win.grab_set()
        
        header = ctk.CTkFrame(win, fg_color="transparent"); header.pack(fill="x", padx=20, pady=(20,10))
        ctk.CTkLabel(header, text="Manajemen Pengguna", font=self.FONT_TITLE).pack(side="left")
        
        frame = ctk.CTkScrollableFrame(win, fg_color=self.COLOR_SECONDARY, corner_radius=10); frame.pack(fill="both", expand=True, padx=20, pady=10)
        
        def fetch():
            for w in frame.winfo_children(): w.destroy()
            ctk.CTkLabel(frame, text="Memuat...").pack(pady=20)
            threading.Thread(target=self._fetch_users_worker, args=(frame,), daemon=True).start()
        
        ctk.CTkButton(header, text="Refresh", command=fetch, width=100, fg_color=self.COLOR_PRIMARY).pack(side="right")
        fetch()

    def _fetch_users_worker(self, frame):
        try:
            headers = {'Authorization': f'Bearer {self.auth_token}'}
            res = requests.get(f"{self.api_base_url}/api/users/", headers=headers, timeout=self.request_timeout)
            if res.status_code == 200:
                data = res.json(); user_list = data.get('users', data.get('data', data if isinstance(data, list) else []))
                self.after(0, self._populate_user_list, frame, user_list)
            else: self.after(0, self._populate_user_list, frame, None, f"Gagal: {res.status_code}")
        except Exception as e: self.after(0, self._populate_user_list, frame, None, f"Error: {e}")

    def _populate_user_list(self, frame, user_list, error_msg=None):
        for w in frame.winfo_children(): w.destroy()
        if error_msg: ctk.CTkLabel(frame, text=error_msg, text_color="red").pack(); return
        if not user_list: ctk.CTkLabel(frame, text="Tidak ada pengguna ditemukan.").pack(); return
        
        for user in user_list:
            user_id = user.get('id')
            row = ctk.CTkFrame(frame, fg_color=self.COLOR_CARD_BG, corner_radius=10, border_width=1, border_color=self.COLOR_CARD_BORDER)
            row.pack(fill="x", pady=5, padx=5); row.grid_columnconfigure(0, weight=1)
            
            info_frame = ctk.CTkFrame(row, fg_color="transparent")
            info_frame.grid(row=0, column=0, sticky="ew", padx=10, pady=10)
            
            ctk.CTkLabel(info_frame, text=user.get('username', 'N/A'), font=("Roboto", 14, "bold")).pack(anchor="w")
            ctk.CTkLabel(info_frame, text=user.get('email', 'N/A'), font=self.FONT_NORMAL, text_color=self.COLOR_TEXT_SECONDARY).pack(anchor="w")
            ctk.CTkLabel(info_frame, text=f"Role: {user.get('role', 'N/A').upper()}", font=("Roboto", 10, "bold"), text_color=self.COLOR_PRIMARY).pack(anchor="w", pady=(5,0))
            
            btn = ctk.CTkButton(row, text="Hapus", fg_color="#982D2D", hover_color="#C62828", width=80, command=lambda uid=user_id: self._confirm_delete_user(uid, frame))
            btn.grid(row=0, column=1, padx=10)
            if user_id == self.logged_in_user_id: btn.configure(state="disabled", text="Anda", fg_color="gray25")

    def _confirm_delete_user(self, user_id, frame):
        dialog = ctk.CTkToplevel(self, fg_color=self.COLOR_CARD_BG); dialog.title("Konfirmasi"); dialog.geometry("300x150"); dialog.transient(self); dialog.grab_set()
        ctk.CTkLabel(dialog, text="Yakin ingin menghapus pengguna ini?", wraplength=280).pack(pady=20, padx=20)
        btn_frame = ctk.CTkFrame(dialog, fg_color="transparent"); btn_frame.pack(pady=10)
        def do_delete():
            dialog.destroy()
            threading.Thread(target=self._delete_user_worker, args=(user_id, frame), daemon=True).start()
        ctk.CTkButton(btn_frame, text="Batal", command=dialog.destroy, fg_color=self.COLOR_SECONDARY).pack(side="left", padx=10)
        ctk.CTkButton(btn_frame, text="Hapus", fg_color="#982D2D", command=do_delete).pack(side="left", padx=10)

    def _delete_user_worker(self, user_id, frame):
        try:
            headers = {'Authorization': f'Bearer {self.auth_token}'}
            res = requests.delete(f"{self.api_base_url}/api/users/{user_id}/", headers=headers, timeout=self.request_timeout)
            if res.status_code < 300:
                self.after(0, self.log, "Admin", f"User {user_id} deleted.")
                self.after(0, lambda: self._fetch_users_worker(frame))
            else: self.after(0, self.log, "Admin", f"Failed to delete {user_id}: {res.status_code}")
        except Exception as e: self.after(0, self.log, "Admin", f"Error deleting {user_id}: {e}")
    
    def log(self, source, message):
        entry = f"[{datetime.now().strftime('%H:%M:%S')}] [{source}] {message}\n"
        def update():
            if hasattr(self, 'log_text') and self.log_text.winfo_exists():
                self.log_text.insert("end", entry); self.log_text.see("end")
        self.after(0, update)

    def toggle_auto(self):
        self.auto_refresh_enabled = self.auto_var.get()
        self.log("Auto", "Auto refresh " + ("dimulai." if self.auto_refresh_enabled else "dihentikan."))
        if self.auto_refresh_enabled: self.start_auto_refresh()

    def start_auto_refresh(self):
        if not self.auto_refresh_enabled: return
        self.fetch_data()
        self.fetch_device_info()
        self.after(self.refresh_interval * 1000, self.start_auto_refresh)
        
    def manual_refresh(self): 
        self.log("Data", "Meminta refresh manual..."); 
        self.fetch_data()
        self.fetch_device_info()
        
    def fetch_data(self):
        try:
            headers = {'Authorization': f'Bearer {self.auth_token}'}
            response = requests.get(self.sensor_api_endpoint, headers=headers, timeout=self.request_timeout)
            if response.status_code == 200:
                data = response.json()
                self.log("API_DEBUG", f"Sensor data received: {str(data)[:200]}")
                self.after(0, self.update_display, data.get('data', data))
            elif response.status_code == 401: self.after(0, self.handle_token_expired)
        except Exception as e: self.log("API", f"Fetch Error: {e}")
            
    def handle_token_expired(self): self.log("Auth", "Token kedaluwarsa."); self.logout()
    
    def logout(self): 
        self.auto_refresh_enabled=False; self.destroy()
        try: from main import main; main()
        except ImportError: print("Could not re-open main login window.")
    
    def on_closing(self): self.auto_refresh_enabled=False; self.destroy()

if __name__ == '__main__':
    current_user_role = 'admin'
    mock_auth_token = "your_test_auth_token"
    mock_user_data = {"user": {"id": 1, "username": "TestAdmin", 'role': 'admin'}} if current_user_role == 'admin' else {"user": {"id": 2, "username": "TestUser", 'role': 'user'}}
    mock_api_endpoint = "http://127.0.0.1:8000/api/sensor-readings/device/4/latest"
    app = DashboardApp(auth_token=mock_auth_token, user_data=mock_user_data, api_endpoint=mock_api_endpoint, request_timeout=5, refresh_interval=5, device_id=4)
    app.protocol("WM_DELETE_WINDOW", app.on_closing)
    app.mainloop()