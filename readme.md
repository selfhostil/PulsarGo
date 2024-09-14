Pulsarego



Pulsarego is a lightweight server monitoring tool written in Go. It allows you to monitor the status of servers or web services, alerting you via Telegram when a service goes down. The tool provides a simple web interface for managing servers, configuring check intervals, and customizing API settings.
Features

    📡 Real-Time Monitoring: Continuously monitor your servers and web services.
    🟢 Customizable Check Interval: Set monitoring intervals between 1 and 60 minutes.
    ⚠️ Alert System: Sends Telegram alerts if a server becomes unavailable.
    🌐 Web Interface: Manage domains, configure intervals, and view server statuses via a web dashboard.
    🧩 Domain Management: Add and remove servers, track status, and view service names.
    🚀 Cross-Platform: Run as a Docker container or as a standalone binary on Windows, macOS, and Linux.
    🎨 Customizable Theme: Choose between themes to style the web interface.
    🔧 Settings Panel: Edit Telegram API tokens and chat IDs directly in the UI.

Screenshots

Include screenshots of the web interface and dashboard here.
Installation
Option 1: Docker

    Clone the repository:

    bash

git clone https://github.com/selfhostil/pulsarego.git
cd pulsarego

Build and run the Docker container:

bash

docker build -t pulsarego .
docker run -d -p 5001:5001 --name pulsarego pulsarego

Access the web interface:

arduino

    http://localhost:5001

Option 2: Standalone Binary

    Clone the repository:

    bash

git clone https://github.com/selfhostil/pulsarego.git
cd pulsarego

Build the binary:

bash

go build -o pulsarego

Run the application:

bash

./pulsarego

Access the web interface:

arduino

    http://localhost:5001

Option 3: Precompiled Binary

Download the precompiled binary from the Releases page, and run it directly on your machine.
Configuration
domain.conf

Add server URLs and their respective service names in the domain.conf file in the format:

bash

http://example.com # Example Service

http://test.com # Test Service

api.conf

Set the Telegram bot token and chat ID in the api.conf file:

makefile

api_token=YOUR_TELEGRAM_BOT_TOKEN

chat_id=YOUR_CHAT_ID

Usage

    Access the web interface at http://localhost:5001.
    Add servers by clicking "Add New Server" and entering the domain and service name.
    Set the monitoring interval using the slider.
    Monitor server statuses in the dashboard.
    Enable or disable monitoring as needed.
    Configure the Telegram API token and chat ID in the Settings panel.

Contributing

We welcome contributions! To contribute, please follow these steps:

    Fork the project.
    Create a new branch (git checkout -b feature/my-feature).
    Commit your changes (git commit -m 'Add new feature').
    Push to the branch (git push origin feature/my-feature).
    Open a pull request.

License

This project is licensed under the MIT License. See the LICENSE file for more information.
Support

For help and support, feel free to open an issue on the GitHub Issues page.

Feel free to add any specific information such as how to run tests or additional commands if necessary. Let me know if you need more modifications!
