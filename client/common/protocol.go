package common

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
)

type OpCode byte
const BUFFER = 4096
const SEPARATOR = ';'
const CON_SEPARATOR = '|'

/**
Base Message structure

| Opcode | 
| 1 byte | 

**/

const (
	CheckClient OpCode = 1
)

/**
CheckClient Message Encoded

| Opcode | Length in bytes | Client1 info (ID;FirstName;LastName;)	|
| 1 byte | 	4 bytes 	   |		Dynamic						|

**/

type CheckContestantMessage struct{
	Con []Contestant
	Op OpCode
}

type Protocol struct {
	conn   net.Conn
	read_buffer bytes.Buffer
}

type Contestant struct {
	Name string
	LastName string
	Id string
	Birth string
}

// Encode
func (m *CheckContestantMessage) encode() []byte {
	con := m.Con

	elem := bytes.Buffer{}
	for _, record := range con {
       // Save fields
		elem.WriteString(record.Id)
		elem.WriteByte(SEPARATOR)
		elem.WriteString(record.Name)
		elem.WriteByte(SEPARATOR)
		elem.WriteString(record.LastName)
		elem.WriteByte(SEPARATOR)
		elem.WriteString(record.Birth)
		
		if record != con[len(con)-1] {
			elem.WriteByte(CON_SEPARATOR)
		}
    }

	buf := bytes.Buffer{}
	
	// Store Opcode
	buf.WriteByte(byte(m.Op))

	len_buf := make([]byte, 4)
    binary.BigEndian.PutUint32(len_buf, uint32(elem.Len()))
	buf.Write(len_buf)
	buf.Write(elem.Bytes())
	
	return buf.Bytes()
}

func (c *Protocol) CheckContestant(contestant []Contestant) error {
	message := CheckContestantMessage{
		Op: CheckClient,
		Con: contestant,
	}	
	buf := message.encode()
	
	return c.send(buf)
}

/* Check response for check contestant message. */
func (c *Protocol) CheckResponse() ([]string, error) {

	
	buf, err := c.recv(4)
	if err != nil {
		return nil, err	
	}
    
	length := binary.BigEndian.Uint32(buf)
	
	buf, err = c.recv(int(length))
	
	parsed := string(buf)
	ids := strings.Split(parsed, string(CON_SEPARATOR))
	
	return ids, nil
}

/* Close socket associated to protocol. */
func (c *Protocol) Close() error {
	return c.conn.Close()
}

/* Wrapper for socker send. */
func (c *Protocol) send(buf []byte) error {
	sent_bytes := 0
	for sent_bytes < len(buf) {
		sent, err := c.conn.Write(buf[sent_bytes:])
		sent_bytes += sent
		
		if err != nil {
			return err	
		}
	}
	
	return nil
}

/* Wrapper for socker recv. Internal buffer to avoid short read. */
func (c *Protocol) recv(size int) ([]byte, error) {

	for c.read_buffer.Len() < size {
		buf := make([]byte, BUFFER)
		_, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		c.read_buffer.Write(buf)
	}

	return c.read_buffer.Next(size), nil
}