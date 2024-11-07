package server

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "strconv"
    "strings"
)

type RESPHandler struct {
    conn net.Conn
}

func NewRESPHandler(conn net.Conn) *RESPHandler {
    return &RESPHandler{conn: conn}
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
    default:
        h.conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command[0])))
    }
}

type RESP interface {
	Decode([]byte, net.Conn)
	Encode()
}