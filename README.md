# ZStack Webhook
ZStack Webhook is a lightweight and highly customizable Go application that transforms ZStack's webhook alerts into messages for various communication platforms, such as Slack and Telegram.

The application serves as a simple HTTP server that listens for incoming ZStack webhook payloads, parses them, and forwards them to your configured destinations.

### Features
- Flexible Output: Display alerts as clean, formatted text or raw JSON.

- Multiple Targets: Simultaneously send alerts to different platforms.
  - slack

  - telegram

  - dingtalk

- Easy to Configure: All settings are managed through a single YAML file.

- Lightweight & Fast: Built with Go for high performance and minimal resource usage.

### Getting Started
Prerequisites
- Go (version 1.18 or higher)

- A ZStack installation with a configured webhook endpoint.

### Installation
1. Clone the repository to your local machine:

```Bash

git clone https://github.com/chijiajian/zstack-webhook.git
cd zstack-webhook
```

2. Build the application:

```Bash

go build -o zstack-webhook .
```

### Usage
#### Configuration
Create a config.yaml file in the same directory as your application. Below is an example configuration:


```YAML

server:
  port: 8080
  https: false
  #cert_file: "/path/to/your/cert.pem"
  #key_file: "/path/to/your/key.pem"

webhooks:
  - type: "slack"
    config:
      url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  - type: "telegram"
    config:
      bot_token: "123456789:AABBCCDDEEFF-GGCCHHIIJJKKLLMMNNOOPP"
      chat_id: "1234567890" # For a group, it starts with a '-'.
  - type: "dingtalk"
    config:
      url: "https://oapi.dingtalk.com/robot/send?access_token=your_access_token"
      # Optional: Add the secret if your bot is configured with signature verification.
      # secret: "your_secret_key"
      fields:
        - alertname
        - status
        - current_value
```

Note: You must obtain the Bot Token and Chat ID from Telegram's @BotFather and a Slack Webhook URL from your Slack workspace.

#### Running the Server
Start the webhook server using the serve command.

```Bash

./zstack-webhook serve --config config.yaml
```

The server will listen for incoming webhook payloads at the /webhook endpoint. Make sure your ZStack webhook is configured to send payloads to http://<your-server-ip>:8080/webhook.

#### Command-Line Arguments
The serve command supports the following flags:

- --config or -c: Specifies the path to the configuration file (default: config.yaml).

- --output or -o: Controls the output format of the logs. Use json to see the raw alert payload or text for the formatted message (default: text).

Example:

```Bash

./zstack-webhook serve --config config.yaml --output json
```

#### How It Works
1. The application starts an HTTP server on the specified port.

2. When a payload is received at /webhook, the application reads the JSON body.

3. It then iterates through each alert in the payload and forwards it to all configured webhook targets.

4. For each target, a goroutine is created to send the message concurrently.

5. Based on the --output flag, the log output is either a formatted message or the raw JSON payload.

Feel free to contribute to this project by opening issues or submitting pull requests!