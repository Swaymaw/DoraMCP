# Dora MCP

- MCP Implementation for web exploration from scratch in Golang for LLMs.

## Instructions to Start searxng server:

---

```zsh
cd searxng_setup
mkdir ./searxng
cp setting.yml ./searxng/setting.yml

docker compose up -d
```

You should be able to access searxng at http://127.0.0.1:8080

## Instructions to Start Dora MCP:

```zsh
go run .
```

Can connect to the MCP server at http://127.0.0.1:8080/mcp
