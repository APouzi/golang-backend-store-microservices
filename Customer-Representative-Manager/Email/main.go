package email

// import (
// 	"context"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"os"

// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/gmail/v1"
// )

// func createMessage(from, to, subject, body string) *gmail.Message {
// 	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)
// 	raw := base64.URLEncoding.EncodeToString([]byte(msg))
// 	return &gmail.Message{Raw: raw}
// }

// func main() {
// 	// Load your credentials
// 	ctx := context.Background()
// 	b, err := os.ReadFile("credentials.json") // OAuth2 credentials
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// Create a config from the credentials
// 	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}

// 	// Authenticate and get a token
// 	client := GetClient(config) // Implement your token retrieval logic
// 	srv, err := gmail.New(client)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Gmail client: %v", err)
// 	}

// 	// Create and send the email
// 	from := "your-email@gmail.com"
// 	to := "recipient@gmail.com"
// 	subject := "Order Summary"
// 	body := "Thank you for your purchase! Your order is being processed."

// 	message := createMessage(from, to, subject, body)
// 	_, err = srv.Users.Messages.Send("me", message).Do()
// 	if err != nil {
// 		log.Fatalf("Unable to send email: %v", err)
// 	}

// 	fmt.Println("Email sent successfully!")
// }
