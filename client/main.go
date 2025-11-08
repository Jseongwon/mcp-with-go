package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// server/main.go 기준: :8080 에서 /mcp 엔드포인트 사용
	c, err := client.NewStreamableHttpClient("http://localhost:8080/mcp")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("close client: %v", err)
		}
	}()

	// MCP Initialize 호출
	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2025-03-26", // 현재 mcp-go가 사용하는 최신 MCP 스펙 버전
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "go-demo-client",
				Version: "0.1.0",
			},
		},
	}

	if _, err := c.Initialize(ctx, initReq); err != nil {
		log.Fatalf("initialize failed: %v", err)
	}
	fmt.Println("✅ Initialized against MCP server")

	// echo 툴 호출
	callReq := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "echo",
			Arguments: map[string]any{
				"message": "이대로 붙여 넣으면 바로 돌아가는 구조예요.  \\n다음 단계로, Jeong의 NMS/쇼핑몰/알림 서비스를 MCP Tool 세트로 노출하는 버전도 원하시면, 지금 이 구조에 맞춰서 바로 설계/코드까지 만들어줄게요.\\n::contentReference[oaicite:4]{index=4}",
			},
		},
	}

	result, err := c.CallTool(ctx, callReq)
	if err != nil {
		log.Fatalf("call tool failed: %v", err)
	}

	// ✅ 여기서부터가 핵심: result는 *mcp.CallToolResult 이고,
	// Content는 result.Content 로 바로 접근해야 합니다 (result.Result.Content 아님)
	if result.IsError {
		fmt.Println("Tool reported an error result:")
		for i, content := range result.Content {
			fmt.Printf("  [%d] %s\n", i, mcp.GetTextFromContent(content))
		}
		return
	}

	fmt.Println("✅ Tool call succeeded. Response content:")
	for i, content := range result.Content {
		// GetTextFromContent: TextContent든 map이든 알아서 string으로 뽑아주는 helper
		text := mcp.GetTextFromContent(content)
		fmt.Printf("  [%d] %s\n", i, text)
	}
}
