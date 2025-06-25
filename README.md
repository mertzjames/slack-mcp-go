# Slack MCP Server (Go)

This project is a Model Context Protocol (MCP) server for Slack, written in Go. It allows you to interact with the Slack API through a set of predefined tools.

## Prerequisites

- Go 1.18 or later
- A Slack Bot Token

## Installation

You can install the server using `go install`:

```bash
go install github.com/rwatts3/slack-mcp-go@latest
```

## Setup

1.  **Set the Slack Bot Token:**

    Export your Slack bot token as an environment variable:

    ```bash
    export SLACK_BOT_TOKEN="your-slack-bot-token"
    ```

2.  **Configure your MCP client:**

    Add the server to your MCP client's configuration file. Here is an example configuration:

    ```json
    {
      "servers": [
        {
          "name": "slack-mcp-go",
          "command": ["slack-mcp-go"],
          "transport": "stdio"
        }
      ]
    }
    ```

3.  **Run the server:**

    The server will be started automatically by your MCP client when a tool is called.

## Available Tools

The following tools are available through the MCP server:

| Tool                  | Description                       |
| --------------------- | --------------------------------- |
| `add_reaction`        | Add a reaction to a message.      |
| `get_channel_history` | Get channel history.              |
| `get_thread_replies`  | Get thread replies.               |
| `get_user_profile`    | Get user profile.                 |
| `get_users`           | Get all users.                    |
| `list_channels`       | List all channels.                |
| `post_message`        | Post a message to a channel.      |
| `reply_to_thread`     | Reply to a thread.                |

## Logging

The server logs to `/tmp/slack-mcp-go.log`.
