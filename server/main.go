package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 1) MCP 서버 생성
	s := server.NewMCPServer(
		"example-streamable-http",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithLogging(),
	)

	// 2) 간단한 echo 툴 정의
	echoTool := mcp.NewTool(
		"echo",
		mcp.WithDescription("입력받은 text를 그대로 돌려줍니다."),
		mcp.WithString(
			"text",
			mcp.Required(),
			mcp.Description("echo 대상 문자열"),
		),
	)

	// 3) echo 툴 핸들러 등록 (타입 세이프 헬퍼 사용)
	s.AddTool(echoTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		text, err := req.RequireString("text")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("echo: " + text), nil
	})

	// 4) Streamable HTTP 서버 래핑 (SSE 아님)
	//    /mcp 경로를 MCP Streamable HTTP 엔드포인트로 사용
	streamSrv := server.NewStreamableHTTPServer(
		s,
		server.WithEndpointPath("/mcp"),
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp", streamSrv)
	mux.Handle("/mcp/", streamSrv)

	log.Println("MCP Streamable HTTP server listening on :3000/mcp")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}