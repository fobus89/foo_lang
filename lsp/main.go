package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

// stdinStdoutCloser объединяет stdin и stdout для LSP
type stdinStdoutCloser struct {
	io.Reader
	io.Writer
}

func (c *stdinStdoutCloser) Close() error {
	return nil
}

func main() {
	log.SetOutput(os.Stderr) // LSP логи в stderr
	log.Printf("Starting Foo Language LSP Server...")

	// Создаем поток с правильным закрытием
	stdin := os.Stdin
	stdout := os.Stdout
	
	stream := jsonrpc2.NewBufferedStream(&stdinStdoutCloser{Reader: stdin, Writer: stdout}, jsonrpc2.VSCodeObjectCodec{})
	handler := NewFooLanguageServer()
	
	conn := jsonrpc2.NewConn(context.Background(), stream, handler)
	
	// Ждем завершения подключения
	select {
	case <-conn.DisconnectNotify():
		log.Printf("LSP connection closed")
	}
}

// FooLanguageServer реализует LSP сервер
type FooLanguageServer struct {
	workspace *Workspace
}

// NewFooLanguageServer создает новый LSP сервер
func NewFooLanguageServer() *FooLanguageServer {
	return &FooLanguageServer{
		workspace: NewWorkspace(),
	}
}

// Handle обрабатывает LSP запросы
func (s *FooLanguageServer) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error handling %s: %v", req.Method, r)
			// Пропускаем ответ при панике чтобы не сломать поток
		}
	}()
	
	switch req.Method {
	case "initialize":
		s.handleInitialize(ctx, conn, req)
	case "initialized":
		// ничего не делаем
	case "textDocument/didOpen":
		s.handleDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		s.handleDidChange(ctx, conn, req)
	case "textDocument/didSave":
		s.handleDidSave(ctx, conn, req)
	case "textDocument/completion":
		s.handleCompletion(ctx, conn, req)
	case "textDocument/hover":
		s.handleHover(ctx, conn, req)
	case "textDocument/definition":
		s.handleDefinition(ctx, conn, req)
	case "textDocument/diagnostic":
		s.handleDiagnostic(ctx, conn, req)
	case "shutdown":
		log.Printf("Received shutdown request")
		conn.Reply(ctx, req.ID, nil)
	case "exit":
		log.Printf("Received exit request")
		os.Exit(0)
	default:
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: "method not supported"})
	}
}

func (s *FooLanguageServer) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	response := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"textDocumentSync": map[string]interface{}{
				"openClose": true,
				"change":    2, // Incremental
				"save":      map[string]interface{}{"includeText": true},
			},
			"completionProvider": map[string]interface{}{
				"triggerCharacters": []string{".", "@"},
			},
			"hoverProvider":      true,
			"definitionProvider": true,
			"diagnosticProvider": map[string]interface{}{
				"interFileDependencies": false,
				"workspaceDiagnostics":  false,
				"identifier":            "foo-lang-diagnostics",
			},
		},
	}
	conn.Reply(ctx, req.ID, response)
}

func (s *FooLanguageServer) handleDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)
	text := textDoc["text"].(string)

	s.workspace.AddDocument(uri, text)
	s.publishDiagnostics(ctx, conn, uri)
}

func (s *FooLanguageServer) handleDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)

	changes := params["contentChanges"].([]interface{})
	if len(changes) > 0 {
		change := changes[0].(map[string]interface{})
		text := change["text"].(string)
		s.workspace.UpdateDocument(uri, text)
		s.publishDiagnostics(ctx, conn, uri)
	}
}

func (s *FooLanguageServer) handleDidSave(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)

	s.publishDiagnostics(ctx, conn, uri)
}

func (s *FooLanguageServer) handleCompletion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)

	posMap := params["position"].(map[string]interface{})
	pos := Position{
		Line:      int(posMap["line"].(float64)),
		Character: int(posMap["character"].(float64)),
	}

	completions := s.workspace.GetCompletions(uri, pos)
	conn.Reply(ctx, req.ID, map[string]interface{}{"items": completions})
}

func (s *FooLanguageServer) handleHover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)

	posMap := params["position"].(map[string]interface{})
	pos := Position{
		Line:      int(posMap["line"].(float64)),
		Character: int(posMap["character"].(float64)),
	}

	hover := s.workspace.GetHover(uri, pos)
	conn.Reply(ctx, req.ID, hover)
}

func (s *FooLanguageServer) handleDefinition(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)

	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)

	posMap := params["position"].(map[string]interface{})
	pos := Position{
		Line:      int(posMap["line"].(float64)),
		Character: int(posMap["character"].(float64)),
	}

	definitions := s.workspace.GetDefinitions(uri, pos)
	conn.Reply(ctx, req.ID, definitions)
}

func (s *FooLanguageServer) handleDiagnostic(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params map[string]interface{}
	json.Unmarshal(*req.Params, &params)
	
	textDoc := params["textDocument"].(map[string]interface{})
	uri := textDoc["uri"].(string)
	
	diagnostics := s.workspace.GetDiagnostics(uri)
	
	response := map[string]interface{}{
		"kind":  "full",
		"items": diagnostics,
	}
	
	conn.Reply(ctx, req.ID, response)
}

func (s *FooLanguageServer) publishDiagnostics(ctx context.Context, conn *jsonrpc2.Conn, uri string) {
	diagnostics := s.workspace.GetDiagnostics(uri)

	params := map[string]interface{}{
		"uri":         uri,
		"diagnostics": diagnostics,
	}

	conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
