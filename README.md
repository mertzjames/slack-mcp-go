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

    Add the server to your MCP client's configuration file.

    #### Generic MCP Client

    Here is an example configuration for a generic MCP client:

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

    #### Claude

    For Claude, you can add the server to your `~/.anthropic/mcp_servers.json` file. Make sure to use the full path to the binary if it's not in your `PATH`.

    ```json
    {
      "servers": [
        {
          "name": "slack-mcp-go",
          "command": ["/path/to/your/go/bin/slack-mcp-go"],
          "transport": "stdio"
        }
      ]
    }
    ```

    #### VSCode Copilot

    For VSCode Copilot, you can add the server to your `settings.json` file. Make sure to use the full path to the binary if it's not in your `PATH`.

    ```json
    {
      "github.copilot.mcp.servers": {
        "slack-mcp-go": {
          "command": ["/path/to/your/go/bin/slack-mcp-go"],
          "transport": "stdio"
        }
      }
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
