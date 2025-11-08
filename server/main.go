package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// MCP 서버 생성
	s := server.NewMCPServer(
		"go-streamable-demo",
		"0.1.0",
		server.WithToolCapabilities(false), // 기본 MCP tool capabilities 광고
	)

	// echo 툴 정의
	echoTool := mcp.NewTool(
		"echo",
		mcp.WithDescription("Echo back the provided message"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("Message to echo"),
		),
	)

	// echo 툴 핸들러 등록
	s.AddTool(echoTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		msg, err := req.RequireString("message")
		if err != nil {
			// 툴 내부 에러는 MCP 에러로 감싸서 반환
			return mcp.NewToolResultErrorFromErr("missing required argument 'message'", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("echo: %s", msg)), nil
	})

	// Streamable HTTP 서버 생성
	streamServer := server.NewStreamableHTTPServer(
		s,
		server.WithEndpointPath("/mcp"), // 기본도 /mcp 이지만 명시적으로 지정
	)

	log.Println("MCP Streamable HTTP server listening on :8080/mcp")
	if err := streamServer.Start(":8080"); err != nil {
		log.Fatalf("streamable HTTP server error: %v", err)
	}
}
