# Slack MCP Server (Go)

This project is a Model Context Protocol (MCP) server for Slack, written in Go. It allows you to interact with the Slack API through a set of predefined tools.

## Prerequisites

- Go 1.18 or later
- A Slack Bot Token

## Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/slack-mcp-go.git
    cd slack-mcp-go
    ```

2.  **Set the Slack Bot Token:**

    Export your Slack bot token as an environment variable:

    ```bash
    export SLACK_BOT_TOKEN="your-slack-bot-token"
    ```

3.  **Build the server:**

    ```bash
    go build
    ```

4.  **Run the server:**

    The server communicates over stdio.

    ```bash
    ./slack-mcp-go
    ```

## Available Tools

The following tools are available through the MCP server:

-   `add_reaction`: Add a reaction to a message.
-   `get_channel_history`: Get channel history.
-   `get_thread_replies`: Get thread replies.
-   `get_user_profile`: Get user profile.
-   `get_users`: Get all users.
-   `list_channels`: List all channels.
-   `post_message`: Post a message to a channel.
-   `reply_to_thread`: Reply to a thread.

## Logging

The server logs to `/tmp/slack-mcp-go.log`.
