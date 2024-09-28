package services

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file}, nil
}

func (l *Logger) LogCommand(clientID string, cmd string) {
	log.SetOutput(l.file)
	log.Printf("[%s] | Client ID :  %s executed command: %s\n", time.Now().Format(time.RFC822), clientID, cmd)

}

func (l *Logger) LogConnection(clientID string) {
	log.SetOutput(l.file)
	log.Printf("[%s] | Client with ID : %s  got connected with the server\n", time.Now().Format(time.RFC822), clientID)

}

func (l *Logger) LogDissconnect(clientID string) {
	log.SetOutput(l.file)
	log.Printf("[%s] | Client ID :%s Dissconnected \n", time.Now().Format(time.RFC822), clientID)

}

func (l *Logger) LogError(clientID string, err error) {
	log.SetOutput(l.file)
	log.Printf("[%s] Client ID : %s got error : %v\n", time.Now().Format(time.RFC822), clientID, err)

}

func (l *Logger) Close() {
	l.file.Close()

}
