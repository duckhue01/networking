package main

import (
	"fmt"
	"log"
	"net/smtp"
)

// The goal of this programming assignment is to create a simple mail client that sends e-mail to any
// recipient. Your client will need to establish a TCP connection with a mail server (e.g., a Google mail
// server), dialogue with the mail server using the SMTP protocol, send an e-mail message to a recipient(e.g., your friend)
// via the mail server, and finally close the TCP connection with the mail server.
// For this assignment, the Companion Website provides the skeleton code for your client. Your job is to
// complete the code and test your client by sending e-mail to different user accounts. You may also try
// sending through different servers (for example, through a Google mail server and through your
// university mail server).

func main() {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("smtp.gmail.com:25")
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err := c.Mail("khue2001hd@gmail.com"); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt("duckhuejs@gmail.com"); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
