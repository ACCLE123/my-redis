package server

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "strconv"
    "strings"
	"myredis/internal/storage"
)

type RESPHandler struct {
    conn net.Conn
	server *Server
}

func NewRESPHandler(conn net.Conn, server *Server) *RESPHandler {
    return &RESPHandler{
		conn: conn,
		server: server,
	}
}

func (h *RESPHandler) Handle() error {
    reader := bufio.NewReader(h.conn)

    line, err := reader.ReadString('\n')
    if err != nil {
        return err
    }
    line = strings.TrimSpace(line)

    switch line[0] {
    case '*':
        count, err := strconv.Atoi(line[1:])
		fmt.Println(line)
        if err != nil {
            return fmt.Errorf("invalid array length: %v", err)
        }

        var command []string
        for i := 0; i < count; i++ {
            line, err = reader.ReadString('\n')
            if err != nil {
                return err
            }
            line = strings.TrimSpace(line)

            if strings.HasPrefix(line, "$") {
                length, err := strconv.Atoi(line[1:])
                if err != nil {
                    return fmt.Errorf("invalid bulk string length: %v", err)
                }

                bulkString := make([]byte, length)
                _, err = io.ReadFull(reader, bulkString)
                if err != nil {
                    return err
                }
                reader.Discard(2)
                command = append(command, string(bulkString))
            }
        }

        fmt.Printf("Received command: %v\n", command)
        h.handleCommand(command)

    default:
        return fmt.Errorf("unsupported command type")
    }
    return nil
}

func (h *RESPHandler) handleCommand(command []string) {
    if len(command) == 0 {
        h.conn.Write([]byte("-ERR unknown command\r\n"))
        return
    }

    switch strings.ToUpper(command[0]) {
    case "PING":
        h.conn.Write([]byte("+PONG\r\n"))
    case "ECHO":
        if len(command) > 1 {
            h.conn.Write([]byte(fmt.Sprintf("+%s\r\n", command[1])))
        } else {
            h.conn.Write([]byte("-ERR wrong number of arguments for 'echo' command\r\n"))
        }
	case "SET":
        if len(command) != 3 {
            h.conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
            return
        }
        h.server.Set(command[1], storage.NewStringObject(command[2]))
        h.conn.Write([]byte("+OK\r\n"))

    case "GET":
        if len(command) != 2 {
            h.conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
            return
        }
        h.server.mu.Lock()
        value, exists := h.server.store[command[1]]
        h.server.mu.Unlock()
        if exists {
            h.conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", value.Len(), value.String())))
        } else {
            h.conn.Write([]byte("$-1\r\n")) // Null Bulk String if key doesn't exist
        }
    default:
        h.conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command[0])))
    }
}