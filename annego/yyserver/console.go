package yyserver

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"goBase/annego/logger"
)

// ConsoleHandle 接收的参数列表，已空格分割。返回命令执行结果
type ConsoleHandle func([]string) string

type consoleCommand struct {
	command string
	help    string
	handle  ConsoleHandle
}

type Console struct {
	commands map[string]consoleCommand
	cmdList  *list.List
	listener net.Listener
}

func NewConsole() *Console {
	return &Console{make(map[string]consoleCommand), list.New(), nil}
}

func (self *Console) AddCommand(command string, help string, handle ConsoleHandle) {
	if self.listener != nil {
		panic("Console is runing")
	}
	if _, ok := self.commands[command]; !ok {
		self.commands[command] = consoleCommand{command, help, handle}
		self.cmdList.PushBack(command)
	} else {
		panic(fmt.Sprintf("Console add command %s exist", command))
	}
}

func (self *Console) AddDefaultCommand() {
	self.AddCommand("setLogLevel", "set logger level, usage: setLogLevel [0 - 7]", cmdSetLogLevel)
	self.AddCommand("getLogLevel", "get logger level, usage: getLogLevel", cmdGetLogLevel)
}

// Start 之前应该完成AddCommand
func (self *Console) Start(addr string) error {
	if self.listener != nil {
		panic("Console is runing")
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	self.AddCommand("help", "print all command", self.help)

	self.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go self.handleConnect(conn)
			} else {
				logger.Warning("accept %v error %v", addr, err)
			}
		}
	}()
	return nil
}

// StartRange 以此探测从addr开始的，trytime个端口
func (self *Console) StartRange(addr string, trytime int) (err error) {
	if self.listener != nil {
		panic("Console is runing")
	}

	addv := strings.Split(addr, ":")
	if len(addv) < 2 {
		return fmt.Errorf("addr format error %s", addr)
	}
	ip := addv[0]
	port, err := strconv.Atoi(addv[1])
	if err != nil {
		return fmt.Errorf("addr format error %s", addr)
	}
	for i := 0; i < trytime; i++ {
		useaddr := fmt.Sprintf("%s:%d", ip, port+i)
		if self.Start(useaddr) == nil {
			return nil
		}
	}
	return fmt.Errorf("try listen %s range %d fail", addr, trytime)
}

// GetListenAddr 获取监听的地址
func (self *Console) GetListenAddr() net.Addr {
	if self.listener == nil {
		return nil
	}
	return self.listener.Addr()
}

func (self *Console) GetListenPort() int {
	if self.listener == nil {
		return 0
	}
	addr := self.listener.Addr()
	tcpaddr := addr.(*net.TCPAddr)
	return tcpaddr.Port
}

func (self *Console) handleConnect(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		byteline, _, err := reader.ReadLine()
		if err != nil {
			logger.Info("console %v read error %v", conn.RemoteAddr(), err)
			break
		}
		line := string(byteline)
		params := strings.Split(line, " ")
		if len(params) > 0 {
			logger.Info("addr %v command %s", conn.RemoteAddr(), params[0])
			var res string
			if cmd, ok := self.commands[params[0]]; ok {
				res = cmd.handle(params)
				res += "\n"
			} else {
				res = "invalid command\n"
			}
			conn.Write([]byte(res))
		}
	}
}

func (self *Console) help(params []string) string {
	result := strings.Builder{}
	result.WriteString("print all command:\n")
	for elem := self.cmdList.Front(); elem != nil; elem = elem.Next() {
		key := elem.Value.(string)
		val, _ := self.commands[key]
		result.WriteString(fmt.Sprintf("%s: %s\n", key, val.help))
	}
	return result.String()
}

func cmdSetLogLevel(params []string) string {
	if len(params) != 2 {
		return "usage: setLogLevel [0 - 7]\n"
	}
	level, err := strconv.Atoi(params[1])
	if err != nil || level < 0 || level > 7 {
		return "usage: setLogLevel [0 - 7]\n"
	}
	logger.SetLogLevel(level)
	return fmt.Sprintf("setLogLevel to %d\n", level)
}

func cmdGetLogLevel(params []string) string {
	level := logger.GetLogLevel()
	return fmt.Sprintf("getLogLevel: %d\n", level)
}

var DefaultConsole *Console

func init() {
	DefaultConsole = NewConsole()
	DefaultConsole.AddDefaultCommand()
}
