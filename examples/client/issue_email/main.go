package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

func main() {
	fmt.Println("💬 Simple Chat Demo")
	fmt.Println("===================")
	fmt.Println("This demo shows basic chat messaging between peers.")
	fmt.Println()

	// Step 1: Create a Self client
	chatClient := createClient()
	defer chatClient.Close()

	fmt.Printf("🆔 Your DID: %s\n", chatClient.DID())
	fmt.Println()

	// Step 3: Discover and connect to a peer
	peer := discoverPeer(chatClient)

	// Step 4: Demonstrate chat messaging
	sendChatMessages(chatClient, peer)

	// Step 5: Send email credential
	sendEmailCredential(chatClient, peer)

	// Keep running to demonstrate ongoing chat capabilities
	select {}
}

// createClient sets up a Self client for chat messaging
func createClient() *client.Client {
	fmt.Println("🔧 Setting up chat client...")

	// Use the simplified client creation - much easier!
	chatClient, err := client.NewSimplified("./simple_chat_storage")
	if err != nil {
		log.Fatal("Failed to create chat client:", err)
	}

	fmt.Println("✅ Chat client created successfully")
	return chatClient
}

// discoverPeer establishes a connection with another peer via QR code
func discoverPeer(chatClient *client.Client) *client.Peer {
	fmt.Println("🔍 Discovering peer for chat...")
	fmt.Println("🔑 Generating QR code for secure connection...")

	// Generate QR code for peer discovery
	qr, err := chatClient.Discovery().GenerateQR()
	if err != nil {
		log.Fatal("Failed to generate QR code:", err)
	}

	fmt.Println("\n📱 SCAN THIS QR CODE with another Self client:")
	fmt.Println("   • Run another instance of this program")
	fmt.Println("   • Use the Self mobile app")
	fmt.Println("   • Any Self SDK application")

	qrCode, err := qr.Unicode()
	if err != nil {
		log.Fatal("Failed to render QR code:", err)
	}
	fmt.Println(qrCode)

	fmt.Println("⏳ Waiting for peer to scan QR code...")
	fmt.Println("   🔐 Secure encrypted connection will be established")
	fmt.Println("   🛑 Press Ctrl+C to cancel")
	fmt.Println()

	// Wait for peer connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Fatal("❌ No peer connected within timeout. Try running another instance of this program.")
		}
		log.Fatal("❌ Failed to connect to peer:", err)
	}

	fmt.Printf("✅ Peer connected: %s\n", peer.DID())
	fmt.Println("🔐 Secure encrypted channel established")
	fmt.Println()

	return peer
}

// sendChatMessages shows basic chat functionality with the connected peer
func sendChatMessages(chatClient *client.Client, peer *client.Peer) {
	fmt.Println("💬 Demonstrating chat messaging...")

	// Send initial greeting
	greeting := fmt.Sprintf("🎉 Hello! Chat demo started at %s. This message is end-to-end encrypted!",
		time.Now().Format("15:04:05"))

	fmt.Println("📤 Sending email credential...")
	err := chatClient.Chat().Send(peer.DID(), greeting)
	if err != nil {
		log.Printf("Failed to send greeting: %v", err)
		return
	}
	fmt.Printf("✅ Greeting sent: \"%s\"\n", greeting)
}

// sendEmailCredential sends an email credential using the fluent API
func sendEmailCredential(client *client.Client, peer *client.Peer) {
	log.Printf("📤 Creating and sending email credential (fluent API) to %s...", peer.DID())

	// Option 2: Fluent API - Issue and Send in one chain
	issueAndSend, err := client.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "EmailCredential"}).
		Subject(peer.DID()).
		Issuer(client.DID()).
		Claim("emailAddress", "test@example.com").
		Claim("verified", true).
		Claim("verificationDate", time.Now().Format("2006-01-02")).
		ValidFrom(time.Now()).
		SignWith(client.DID(), time.Now()).
		IssueAndSend(client)

	if err != nil {
		log.Printf("Failed to create credential: %v", err)
		return
	}

	err = issueAndSend.Send(peer.DID())
	if err != nil {
		log.Printf("Failed to send credential: %v", err)
		return
	}

	log.Printf("✅ Email credential sent successfully (fluent API)")
}
