# Go MCP Example – Streamable HTTP (mark3labs/mcp-go)

`mark3labs/mcp-go`를 사용해 **Streamable HTTP** 프로토콜 기반 MCP 서버와
테스트 클라이언트를 구현한 최소 예제입니다.

- 서버: `server/main.go`
- 클라이언트: `client/main.go`
- 프로토콜: [MCP Streamable HTTP](https://modelcontextprotocol.io/) 기반
- SSE 사용 안 함 (단, 라이브러리 구현상 Streamable HTTP 내에서 SSE 모드도 지원할 수 있음)

## 1. 요구사항

- Go 1.21 이상
- Git
- 인터넷 연결 (모듈 다운로드용)

## 2. 의존성

이 예제는 다음 모듈을 사용합니다.

```bash
github.com/mark3labs/mcp-go v0.43.0