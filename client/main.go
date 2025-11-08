package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	ctx := context.Background()

	// 1) Streamable HTTP 트랜스포트 생성
	t, err := transport.NewStreamableHTTP("http://localhost:3000/mcp")
	if err != nil {
		log.Fatalf("failed to create transport: %v", err)
	}
	defer t.Close()

	// 2) MCP 클라이언트 생성
	c := client.NewClient(t)

	// 3) 시작
	if err := c.Start(ctx); err != nil {
		log.Fatalf("failed to start client: %v", err)
	}
	defer c.Close()

	// 4) initialize 호출 (Streamable HTTP 프로토콜 핸드셰이크)
	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "example-client",
				Version: "0.0.1",
			},
		},
	})
	if err != nil {
		log.Fatalf("initialize failed: %v", err)
	}

	// 5) echo 툴 호출
	resp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "echo",
			Arguments: map[string]any{
				"text": "Hello from Streamable HTTP!",
			},
		},
	})
	if err != nil {
		log.Fatalf("call tool failed: %v", err)
	}

	// 6) 결과 출력 (텍스트 컨텐츠만 단순 파싱)
	for _, content := range resp.Result.Content {
		if txt, ok := content.(mcp.TextContent); ok {
			fmt.Println("Tool result:", txt.Text)
		}
	}
}
