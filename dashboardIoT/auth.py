import customtkinter as ctk
import requests
import threading
from PIL import Image, ImageTk

class AuthWindow(ctk.CTkToplevel):
    def __init__(self, parent, callback, api_base_url, request_timeout):
        super().__init__(parent)
        
        self.callback = callback
        self.api_base_url = api_base_url
        self.request_timeout = request_timeout
        self.auth_token = ""
        self.user_data = {}

        # --- Modern UI Configuration ---
        self.title("Smart Garden Access")
        # --- FIX: Increased window height further to ensure visibility ---
        self.geometry("400x650")
        self.resizable(False, False)
        
        self.transient(parent)
        self.grab_set()

        # --- Colors and Fonts ---
        self.COLOR_PRIMARY = "#2FA235"
        self.COLOR_SECONDARY = "#222B36"
        self.COLOR_BACKGROUND = "#1C242D"
        self.COLOR_TEXT = "#E0E0E0"
        self.COLOR_ERROR = "#D32F2F"
        self.COLOR_SUCCESS = "#388E3C"

        self.FONT_BOLD = ("Arial", 16, "bold")
        self.FONT_NORMAL = ("Arial", 12)
        
        self.configure(fg_color=self.COLOR_BACKGROUND)
        self.setup_ui()
        
    def setup_ui(self):
        """Setup the modern user interface."""
        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(2, weight=1)

        # --- Header ---
        header_frame = ctk.CTkFrame(self, fg_color="transparent")
        header_frame.grid(row=0, column=0, pady=(30, 15))
        
        ctk.CTkLabel(
            header_frame, 
            text="🌿 SmartGarden", 
            font=("Arial", 28, "bold"),
            text_color=self.COLOR_PRIMARY
        ).pack()
        ctk.CTkLabel(
            header_frame, 
            text="Welcome Back", 
            font=self.FONT_NORMAL,
            text_color=self.COLOR_TEXT
        ).pack()

        # --- Segmented Button for Login/Register Toggle ---
        toggle_frame = ctk.CTkFrame(self, fg_color=self.COLOR_SECONDARY, corner_radius=15)
        toggle_frame.grid(row=1, column=0, pady=10)

        self.login_toggle_btn = ctk.CTkButton(
            toggle_frame, text="Login", font=self.FONT_BOLD,
            command=lambda: self.show_frame("login"),
            fg_color=self.COLOR_PRIMARY, hover_color=self.COLOR_PRIMARY,
            corner_radius=15
        )
        self.login_toggle_btn.pack(side="left", padx=5, pady=5)

        self.register_toggle_btn = ctk.CTkButton(
            toggle_frame, text="Register", font=self.FONT_BOLD,
            command=lambda: self.show_frame("register"),
            fg_color="transparent", hover_color="#333"
        )
        self.register_toggle_btn.pack(side="left", padx=5, pady=5)

        # --- Main content area for forms ---
        content_frame = ctk.CTkFrame(self, fg_color="transparent")
        content_frame.grid(row=2, column=0, sticky="nsew", padx=20)
        content_frame.grid_rowconfigure(0, weight=1)
        content_frame.grid_columnconfigure(0, weight=1)
        
        self.login_frame = self._create_login_frame(content_frame)
        self.register_frame = self._create_register_frame(content_frame)
        self.login_frame.grid(row=0, column=0, sticky="nsew")
        self.register_frame.grid(row=0, column=0, sticky="nsew")

        # --- Status Label ---
        self.status_label = ctk.CTkLabel(self, text="", font=self.FONT_NORMAL, height=30)
        self.status_label.grid(row=3, column=0, sticky="ew", padx=20, pady=(10, 20))

        self.show_frame("login")

    def _create_login_frame(self, parent):
        frame = ctk.CTkFrame(parent, fg_color="transparent")
        frame.grid_columnconfigure(0, weight=1)

        self.login_email = self._create_entry_with_label(frame, "Email", 0)
        self.login_password = self._create_entry_with_label(frame, "Password", 2, show="*")
        
        # --- FIX: Bind Enter key to focus next widget ---
        self.login_email.bind("<Return>", lambda e: self.login_password.focus())
        self.login_password.bind("<Return>", lambda e: self.handle_login())
        
        self.login_btn = ctk.CTkButton(
            frame, text="Login Securely", font=self.FONT_BOLD, 
            command=self.handle_login, fg_color=self.COLOR_PRIMARY,
            hover_color="#277C2B", height=40, corner_radius=10
        )
        self.login_btn.grid(row=4, column=0, pady=(25, 10), sticky="ew")
        
        return frame

    def _create_register_frame(self, parent):
        frame = ctk.CTkFrame(parent, fg_color="transparent")
        frame.grid_columnconfigure(0, weight=1)

        self.reg_username = self._create_entry_with_label(frame, "Username", 0)
        self.reg_email = self._create_entry_with_label(frame, "Email", 2)
        self.reg_password = self._create_entry_with_label(frame, "Create Password", 4, show="*")
        self.reg_confirm = self._create_entry_with_label(frame, "Confirm Password", 6, show="*")

        # --- FIX: Bind Enter key to focus next widget or submit ---
        self.reg_username.bind("<Return>", lambda e: self.reg_email.focus())
        self.reg_email.bind("<Return>", lambda e: self.reg_password.focus())
        self.reg_password.bind("<Return>", lambda e: self.reg_confirm.focus())
        self.reg_confirm.bind("<Return>", lambda e: self.handle_register())
        
        self.register_btn = ctk.CTkButton(
            frame, text="Create Account", font=self.FONT_BOLD, 
            command=self.handle_register, fg_color=self.COLOR_PRIMARY,
            hover_color="#277C2B", height=40, corner_radius=10
        )
        # --- FIX: Adjusted padding to give more space ---
        self.register_btn.grid(row=8, column=0, pady=(15, 10), sticky="ew")
        return frame
        
    def _create_entry_with_label(self, parent, label_text, row, show=None):
        """Helper to create a label and an entry field."""
        ctk.CTkLabel(parent, text=label_text, font=self.FONT_NORMAL, text_color=self.COLOR_TEXT).grid(
            row=row, column=0, sticky="w", padx=5
        )
        entry = ctk.CTkEntry(
            parent, font=self.FONT_NORMAL, show=show, height=35,
            fg_color=self.COLOR_SECONDARY, border_color="#333",
            corner_radius=8
        )
        entry.grid(row=row + 1, column=0, sticky="ew", pady=(5, 15))
        return entry

    def show_frame(self, frame_name):
        """Shows the selected frame (login/register) and updates toggle button styles."""
        if frame_name == "login":
            self.login_frame.tkraise()
            self.login_toggle_btn.configure(fg_color=self.COLOR_PRIMARY, hover_color=self.COLOR_PRIMARY)
            self.register_toggle_btn.configure(fg_color="transparent", hover_color="#333")
        else:
            self.register_frame.tkraise()
            self.register_toggle_btn.configure(fg_color=self.COLOR_PRIMARY, hover_color=self.COLOR_PRIMARY)
            self.login_toggle_btn.configure(fg_color="transparent", hover_color="#333")

    def handle_login(self):
        email = self.login_email.get().strip()
        password = self.login_password.get().strip()
        if not email or not password:
            self.show_status("Email dan password harus diisi", self.COLOR_ERROR)
            return
        
        self.login_btn.configure(text="Authenticating...", state="disabled")
        self.show_status("Connecting...", self.COLOR_TEXT, clear=False)
        threading.Thread(target=self._login_worker, args=(email, password), daemon=True).start()
        
    def _login_worker(self, email, password):
        """Worker thread for login."""
        try:
            login_data = {
                "email": email,
                "password": password
            }
            
            print(f"DEBUG: Mengirim data login: {login_data}")
            
            response = requests.post(f"{self.api_base_url}/auth/login", json=login_data, timeout=self.request_timeout)
            
            if response.status_code == 200:
                data = response.json()
                if 'token' in data and 'user' in data:
                    self.auth_token = data['token']
                    self.user_data = data 
                    self.after(0, lambda: self.show_status("Login Berhasil!", self.COLOR_SUCCESS))
                    self.after(1000, self.success_login)
                else:
                    self.after(0, lambda: self.show_status("Respon tidak valid dari server", self.COLOR_ERROR))
                    self.after(0, self.reset_login_button)
            else:
                error_msg = f"Login Gagal: {response.json().get('message', 'Cek kembali kredensial Anda')}"
                self.after(0, lambda: self.show_status(error_msg, self.COLOR_ERROR))
                self.after(0, self.reset_login_button)
                
        except requests.exceptions.RequestException:
            self.after(0, lambda: self.show_status("Kesalahan Koneksi. Periksa server.", self.COLOR_ERROR))
            self.after(0, self.reset_login_button)
        except Exception as e:
            self.after(0, lambda: self.show_status(f"Error: {e}", self.COLOR_ERROR))
            self.after(0, self.reset_login_button)
            
    def handle_register(self):
        username = self.reg_username.get().strip()
        email = self.reg_email.get().strip()
        password = self.reg_password.get().strip()
        confirm = self.reg_confirm.get().strip()
        
        if not all([username, email, password]):
            self.show_status("Semua kolom harus diisi!", self.COLOR_ERROR)
            return
        if password != confirm:
            self.show_status("Password tidak cocok!", self.COLOR_ERROR)
            return
        
        self.register_btn.configure(text="Creating Account...", state="disabled")
        self.show_status("Please wait...", self.COLOR_TEXT, clear=False)
        threading.Thread(target=self._register_worker, args=(username, email, password), daemon=True).start()
        
    def _register_worker(self, username, email, password):
        try:
            register_data = {"username": username, "email": email, "password": password}
            response = requests.post(f"{self.api_base_url}/auth/register", json=register_data, timeout=self.request_timeout)
            
            if response.status_code in [200, 201]:
                self.after(0, lambda: self.show_status("Akun berhasil dibuat! Silakan login.", self.COLOR_SUCCESS))
                self.after(1500, lambda: self.show_frame("login"))
                self.after(0, self.reset_register_button)
            else:
                error_msg = f"Registrasi Gagal: {response.json().get('message', 'Kesalahan server')}"
                self.after(0, lambda: self.show_status(error_msg, self.COLOR_ERROR))
                self.after(0, self.reset_register_button)
                
        except requests.exceptions.RequestException:
            self.after(0, lambda: self.show_status("Kesalahan Koneksi. Periksa server.", self.COLOR_ERROR))
            self.after(0, self.reset_register_button)
        except Exception as e:
            self.after(0, lambda: self.show_status(f"Error: {e}", self.COLOR_ERROR))
            self.after(0, self.reset_register_button)
            
    def reset_login_button(self):
        self.login_btn.configure(text="Login Securely", state="normal")
        
    def reset_register_button(self):
        self.register_btn.configure(text="Create Account", state="normal")
        
    def show_status(self, message, color, clear=True):
        self.status_label.configure(text=message, text_color="white", fg_color=color)
        if clear:
            self.status_label.after(3000, lambda: self.status_label.configure(text="", fg_color="transparent"))
        
    def success_login(self):
        self.callback(self.auth_token, self.user_data)
        self.destroy()

