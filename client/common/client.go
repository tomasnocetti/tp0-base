package common

import (
	"bytes"
	"encoding/csv"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)


const TIME_LAPSE = 8

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	ContestantsPath string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	protocol Protocol
	reader FileReader
	done chan bool
}

// Contestans File Reader
type FileReader struct {
	file *os.File
	reader *csv.Reader
}

func (c *FileReader) getContestants() ([]Contestant, error) {
	records, err := c.reader.ReadAll()

	if err != nil {
		return nil, err
	}
	
	contestants := make([]Contestant, 0)

	// using for loop
    for _, record := range records {
        contestants = append(contestants, Contestant{
			Name: record[0],
			LastName: record[1],
			Id: record[2],
			Birth: record[3],
		})
    }
	return contestants, nil

}

func (c *FileReader) Close() {
	c.file.Close()
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

func (c *Client) createFileReader() error {
	f, err := os.Open(c.config.ContestantsPath)
    if err != nil {
        log.Fatal(err)
    }
	r := csv.NewReader(f)
	c.reader = FileReader{f, r}

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
	// Init file reader
	c.createFileReader()

	// Establish connection
	c.createClientSocket()
	
	// Wait for interruption
	go c.gracefullExit()
	
	// Process client
	go c.processLottery()

	// Wait for event to happen
	<- c.done

	c.getStats()

	log.Infof("[CLIENT %v] Closing connection", c.config.ID)
	c.protocol.Close()
	c.reader.Close()
}

func (c *Client) processLottery() {
	
	cont, err := c.reader.getContestants()

	return_err := c.checkContestants(cont)

	if(return_err != nil){
		log.Errorf(
			"[CLIENT %v] Error writing to socket. %v.",
			c.config.ID,
			err,
		)
		c.done <- false
	}
	
	
	res, err := c.checkResponse()
	
	if(err != nil){
		log.Errorf(
			"[CLIENT %v] Error reading from socket. %v.",
			c.config.ID,
			err,
		)
		c.done <- false
	}
	log.Debugf(
		"[CLIENT %v] From %v contestants, there were %v winners.",
		c.config.ID,
		len(cont),
		len(res),
	)

	log.Infof(
		"[CLIENT %v] There is a total of %.2f%% winners!",
		c.config.ID,
		float32(len(res)) / float32(len(cont)) * 100,
	)

	c.done <- true
}

func (c *Client) getStats() {
	keep_asking := true
	for keep_asking {
		log.Infof(
			"[CLIENT %v] Getting Stats!",
			c.config.ID,
		)
	
		stats, err := c.protocol.GetStats()
	
		if(err != nil){
			log.Errorf(
				"[CLIENT %v] Error writing to socket. %v.",
				c.config.ID,
				err,
			)
		}
		
		log.Infof(
			"[CLIENT %v] Got Stats. Partial: %v. Total: %v",
			c.config.ID,
			stats.Partial,
			stats.Size,
		)

		if (!stats.Partial) {
			return
		}

		time.Sleep(TIME_LAPSE * time.Second)
	}
	
	
}



func (c *Client) checkContestants(con []Contestant) error {
	log.Infof(
		"[CLIENT %v] Checking information for a total of %v contestants",
		c.config.ID,
		len(con),
	)

	return c.protocol.CheckContestant(con)
}

func (c *Client) checkResponse() ([]string, error) {
	log.Infof(
		"[CLIENT %v] Waiting for central Response.",
		c.config.ID,
	)
	
	res, err := c.protocol.CheckResponse()
	
	if(err != nil){
		return nil, err
	}

	return res, nil
}
