package common

import (
	"bytes"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	Contestant Contestant
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	protocol Protocol
	done chan bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		done : make(chan bool, 1),
	}

	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"[CLIENT %v] Could not connect to server. Error: %v",
			c.config.ID,
			err,
		)
	}
	
	c.protocol = Protocol{conn, bytes.Buffer{}}
	return nil
}

func (c *Client) gracefullExit() {
	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	
	select {
		case <- sigs:
			log.Infof("[CLIENT %v] Signal Interruption", c.config.ID)
		}
	c.done <- true
		
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// Establish connection
	c.createClientSocket()
	
	// Wait for interruption
	go c.gracefullExit()
	
	// Process client
	go c.processLottery()

	// Wait for event to happen
	<- c.done

	log.Infof("[CLIENT %v] Closing connection", c.config.ID)
	c.protocol.Close()
}

func (c *Client) processLottery() {
	
	err := c.checkContestant()
	
	if(err != nil){
		log.Errorf(
			"[CLIENT %v] Error writing to socket. %v.",
			c.config.ID,
			err,
		)
		c.done <- false
	}
	
	err = c.checkResponse()
	
	if(err != nil){
		log.Errorf(
			"[CLIENT %v] Error reading from socket. %v.",
			c.config.ID,
			err,
		)
		c.done <- false
	}

	c.done <- true
}

func (c *Client) checkContestant() error {
	con := c.config.Contestant
	
	log.Infof(`[CLIENT %v] Checking Contestant:
	Contestant Id: %v
	Contestant Name: %v
	Contestant Last Name: %v
	Contestant Birth: %v`, 
		c.config.ID,
		con.Id,
		con.Name,
		con.LastName,
		con.Birth,
	)
	
	return c.protocol.CheckContestant(&con)
}

func (c *Client) checkResponse() error {
	con := c.config.Contestant
	res, err := c.protocol.CheckResponse()
	
	if(err != nil){
		return err
	}

	if res {
		log.Infof(`[CLIENT %v] It's a Winner!:
		Contestant Id: %v`, 
		c.config.ID,
		con.Id,
		)
	} else {
		log.Infof(`[CLIENT %v] Best Luck next time!:
		Contestant Id: %v`, 
		c.config.ID,
		con.Id,
		)
	}
	return nil
}

