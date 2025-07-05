import customtkinter as ctk
from auth import AuthWindow
from dashboard import DashboardApp

# KONFIGURASI API
API_SERVER_IP = "192.168.39.89"
API_SERVER_PORT = "8080"
DEVICE_ID = "4"
API_BASE_URL = f"http://{API_SERVER_IP}:{API_SERVER_PORT}/api"
API_ENDPOINT = f"{API_BASE_URL}/sensor-readings/device/{DEVICE_ID}/latest"

# KONFIGURASI
REFRESH_INTERVAL = 1  # 1 detik
REQUEST_TIMEOUT = 2

# Set appearance
ctk.set_appearance_mode("dark")
ctk.set_default_color_theme("green")

class SmartGardenMain(ctk.CTk):
    def __init__(self):
        super().__init__()
        
        self.title("🌿 Smart Garden - Starting...")
        self.geometry("300x200")
        self.resizable(False, False)
        
        # Center window
        self.center_window()
        
        # Setup simple UI
        self.setup_ui()
        
        # Show auth window immediately
        self.after(100, self.show_auth_window)
        
    def center_window(self):
        """Center window on screen"""
        self.update_idletasks()
        width = self.winfo_width()
        height = self.winfo_height()
        x = (self.winfo_screenwidth() // 2) - (width // 2)
        y = (self.winfo_screenheight() // 2) - (height // 2)
        self.geometry(f"{width}x{height}+{x}+{y}")
        
    def setup_ui(self):
        """Setup simple splash screen"""
        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(0, weight=1)
        
        # Main frame
        main_frame = ctk.CTkFrame(self)
        main_frame.grid(row=0, column=0, sticky="nsew", padx=20, pady=20)
        main_frame.grid_columnconfigure(0, weight=1)
        
        # Logo
        logo = ctk.CTkLabel(main_frame, text="🌿", font=ctk.CTkFont(size=48))
        logo.pack(pady=20)
        
        # Title
        title = ctk.CTkLabel(main_frame, text="Smart Garden", 
                            font=ctk.CTkFont(size=20, weight="bold"))
        title.pack()
        
        # Status
        self.status_label = ctk.CTkLabel(main_frame, text="Initializing...", 
                                        font=ctk.CTkFont(size=12))
        self.status_label.pack(pady=10)
        
    def show_auth_window(self):
        """Show authentication window"""
        self.status_label.configure(text="Opening login window...")
        self.withdraw()  # Hide main window
        
        auth_window = AuthWindow(self, self.on_auth_success, API_BASE_URL, REQUEST_TIMEOUT)
        
    def on_auth_success(self, auth_token, user_data):
        """Handle successful authentication"""
        self.status_label.configure(text="Login successful! Opening dashboard...")
        
        # Close main window and open dashboard
        self.destroy()
        
        # Create and show dashboard
        dashboard = DashboardApp(auth_token, user_data, API_ENDPOINT, 
                                REQUEST_TIMEOUT, REFRESH_INTERVAL)
        dashboard.protocol("WM_DELETE_WINDOW", dashboard.on_closing)
        
        # Auto start monitoring
        dashboard.auto_var.set(True)
        dashboard.toggle_auto()
        
        dashboard.mainloop()

def main():
    """Main function"""
    try:
        app = SmartGardenMain()
        app.mainloop()
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()