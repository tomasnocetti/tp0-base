package common

import (
	"bytes"
	"encoding/binary"
	"net"
)

type OpCode byte
const BUFFER = 4096

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
	Con Contestant
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
	buf := bytes.Buffer{}
	
	// Store Opcode
	buf.WriteByte(byte(m.Op))
	
	con := m.Con
	
	// Calculate length and store it
	length := len(con.Name) + len(con.LastName) + len(con.Birth) + len(con.Id) + 3
	len_buf := make([]byte, 4)
    binary.BigEndian.PutUint32(len_buf, uint32(length))
	buf.Write(len_buf)

	// Save fields
	buf.WriteString(con.Id)
	buf.WriteByte(';')
	buf.WriteString(con.Name)
	buf.WriteByte(';')
	buf.WriteString(con.LastName)
	buf.WriteByte(';')
	buf.WriteString(con.Birth)

	return buf.Bytes()
}

func (c *Protocol) CheckContestant(contestant *Contestant) error {
	message := CheckContestantMessage{
		Op: CheckClient,
		Con: *contestant,
	}	
	buf := message.encode()
	
	return c.send(buf)
}

/* Check response for check contestant message. */
func (c *Protocol) CheckResponse() (bool, error) {

	
	buf, err := c.recv(2)
	if err != nil {
		return false, err	
	}
    
	data := binary.BigEndian.Uint16(buf)
	if data == 1 {
		return true, nil
	}
	
	return false, nil
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
		_, err := c.read_buffer.ReadFrom(c.conn)
		
		if err != nil {
			return nil, err
		}
	}

	return c.read_buffer.Next(size), nil
}