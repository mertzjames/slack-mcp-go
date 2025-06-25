package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/slack-go/slack"
)

// ToolHandler is an interface for MCP tool handlers.
type ToolHandler interface {
	Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// SlackAPI is an interface to abstract the slack client for testing.
type SlackAPI interface {
	AddReaction(name string, item slack.ItemRef) error
	GetConversationHistory(params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error)
	GetConversationReplies(params *slack.GetConversationRepliesParameters) (msgs []slack.Message, hasMore bool, nextCursor string, err error)
	GetUserProfile(params *slack.GetUserProfileParameters) (*slack.UserProfile, error)
	GetUsers(ctx context.Context) ([]slack.User, error)
	GetConversations(params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error)
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

// SlackClient is a wrapper around the slack.Client that implements the SlackAPI interface.
type SlackClient struct {
	*slack.Client
}

func (c *SlackClient) GetUsers(ctx context.Context) ([]slack.User, error) {
	return c.Client.GetUsersContext(ctx)
}

// BaseHandler is a base struct for Slack tool handlers.
type BaseHandler struct {
	api SlackAPI
}

// AddReactionHandler handles adding a reaction to a message.
type AddReactionHandler struct {
	*BaseHandler
}

func (h *AddReactionHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channel, _ := args["channel"].(string)
	timestamp, _ := args["timestamp"].(string)
	reaction, _ := args["reaction"].(string)
	err := h.api.AddReaction(reaction, slack.NewRefToMessage(channel, timestamp))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText("success"), nil
}

// GetChannelHistoryHandler handles getting channel history.
type GetChannelHistoryHandler struct {
	*BaseHandler
}

func (h *GetChannelHistoryHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channel, _ := args["channel"].(string)
	history, err := h.api.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: channel})
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(history)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// GetThreadRepliesHandler handles getting thread replies.
type GetThreadRepliesHandler struct {
	*BaseHandler
}

func (h *GetThreadRepliesHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channel, _ := args["channel"].(string)
	timestamp, _ := args["timestamp"].(string)
	messages, _, _, err := h.api.GetConversationReplies(&slack.GetConversationRepliesParameters{ChannelID: channel, Timestamp: timestamp})
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// GetUserProfileHandler handles getting a user's profile.
type GetUserProfileHandler struct {
	*BaseHandler
}

func (h *GetUserProfileHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	userID, _ := args["user_id"].(string)
	user, err := h.api.GetUserProfile(&slack.GetUserProfileParameters{UserID: userID})
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// GetUsersHandler handles getting all users.
type GetUsersHandler struct {
	*BaseHandler
}

func (h *GetUsersHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	users, err := h.api.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// ListChannelsHandler handles listing all channels.
type ListChannelsHandler struct {
	*BaseHandler
}

func (h *ListChannelsHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	channels, _, err := h.api.GetConversations(&slack.GetConversationsParameters{})
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(channels)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// PostMessageHandler handles posting a message to a channel.
type PostMessageHandler struct {
	*BaseHandler
}

func (h *PostMessageHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channel, _ := args["channel"].(string)
	text, _ := args["text"].(string)
	respMsg, timestamp, err := h.api.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(map[string]string{"channel": respMsg, "timestamp": timestamp})
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

// ReplyToThreadHandler handles replying to a thread.
type ReplyToThreadHandler struct {
	*BaseHandler
}

func (h *ReplyToThreadHandler) Call(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channel, _ := args["channel"].(string)
	timestamp, _ := args["timestamp"].(string)
	text, _ := args["text"].(string)
	respMsg, respTimestamp, err := h.api.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionTS(timestamp))
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(map[string]string{"channel": respMsg, "timestamp": respTimestamp})
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(res)),
		},
	}, nil
}

func main() {
	// Setup logging
	logFile, err := os.OpenFile("/tmp/slack-mcp-go.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Check if SLACK_BOT_TOKEN is set
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackToken == "" {
		fmt.Fprintln(os.Stderr, "Error: SLACK_BOT_TOKEN environment variable not set")
		log.Println("Error: SLACK_BOT_TOKEN environment variable not set")
		os.Exit(1)
	}

	api := &SlackClient{slack.New(slackToken)}
	baseHandler := &BaseHandler{api: api}

	mcpServer := server.NewMCPServer(
		"slack-mcp-go",
		"1.0.0",
	)

	addTool(mcpServer, "add_reaction", &AddReactionHandler{baseHandler},
		mcp.WithDescription("Add a reaction to a message."),
		mcp.WithString("channel", mcp.Description("Channel ID"), mcp.Required()),
		mcp.WithString("timestamp", mcp.Description("Message timestamp"), mcp.Required()),
		mcp.WithString("reaction", mcp.Description("Reaction name"), mcp.Required()),
	)

	addTool(mcpServer, "get_channel_history", &GetChannelHistoryHandler{baseHandler},
		mcp.WithDescription("Get channel history."),
		mcp.WithString("channel", mcp.Description("Channel ID"), mcp.Required()),
	)

	addTool(mcpServer, "get_thread_replies", &GetThreadRepliesHandler{baseHandler},
		mcp.WithDescription("Get thread replies."),
		mcp.WithString("channel", mcp.Description("Channel ID"), mcp.Required()),
		mcp.WithString("timestamp", mcp.Description("Message timestamp"), mcp.Required()),
	)

	addTool(mcpServer, "get_user_profile", &GetUserProfileHandler{baseHandler},
		mcp.WithDescription("Get user profile."),
		mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
	)

	addTool(mcpServer, "get_users", &GetUsersHandler{baseHandler}, mcp.WithDescription("Get all users."))

	addTool(mcpServer, "list_channels", &ListChannelsHandler{baseHandler}, mcp.WithDescription("List all channels."))

	addTool(mcpServer, "post_message", &PostMessageHandler{baseHandler},
		mcp.WithDescription("Post a message to a channel."),
		mcp.WithString("channel", mcp.Description("Channel ID"), mcp.Required()),
		mcp.WithString("text", mcp.Description("Message text"), mcp.Required()),
	)

	addTool(mcpServer, "reply_to_thread", &ReplyToThreadHandler{baseHandler},
		mcp.WithDescription("Reply to a thread."),
		mcp.WithString("channel", mcp.Description("Channel ID"), mcp.Required()),
		mcp.WithString("timestamp", mcp.Description("Message timestamp"), mcp.Required()),
		mcp.WithString("text", mcp.Description("Message text"), mcp.Required()),
	)

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func addTool(s *server.MCPServer, name string, handler ToolHandler, opts ...mcp.ToolOption) {
	tool := mcp.NewTool(name, opts...)
	s.AddTool(tool, handler.Call)
}
