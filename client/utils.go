package client

import (
	"context"
	"time"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/message"
)

// QuickCredentialExchange performs a simple credential exchange between two clients
func QuickCredentialExchange(requester, responder *Client, credentialType []string, timeout time.Duration) (*CredentialResponse, error) {
	// Set up a simple handler on the responder
	responseSent := make(chan bool, 1)
	responder.Credentials().OnPresentationRequest(func(req *IncomingCredentialRequest) {
		// For demo purposes, just reject
		req.Reject()
		responseSent <- true
	})

	// Send request
	details := []*CredentialDetail{
		{
			CredentialType: credentialType,
			Parameters: []*CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "id",
					Value:    "",
				},
			},
		},
	}

	req, err := requester.Credentials().RequestPresentationWithTimeout(responder.DID(), details, timeout)
	if err != nil {
		return nil, err
	}

	// Wait for response
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return req.WaitForResponse(ctx)
}

// CreateSimpleEmailCredential creates a basic email credential
func CreateSimpleEmailCredential(client *Client, subjectDID, emailAddress string) (*credential.VerifiableCredential, error) {
	return client.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).
		Subject(subjectDID).
		Issuer(client.DID()).
		Claim("emailAddress", emailAddress).
		Claim("verified", true).
		ValidFrom(time.Now()).
		SignWith(client.DID(), time.Now()).
		Issue(client)
}

// CreateSimpleProfileCredential creates a basic profile credential
func CreateSimpleProfileCredential(client *Client, subjectDID, firstName, lastName, country string) (*credential.VerifiableCredential, error) {
	return client.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeProfileName).
		Subject(subjectDID).
		Issuer(client.DID()).
		Claim("firstName", firstName).
		Claim("lastName", lastName).
		Claim("country", country).
		ValidFrom(time.Now()).
		SignWith(client.DID(), time.Now()).
		Issue(client)
}

// CreateSimpleEducationCredential creates a basic education credential
func CreateSimpleEducationCredential(client *Client, subjectDID, degree, institution string, graduationYear int, gpa float64) (*credential.VerifiableCredential, error) {
	return client.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "EducationCredential"}).
		Subject(subjectDID).
		Issuer(client.DID()).
		Claim("degree", degree).
		Claim("institution", institution).
		Claim("graduationYear", graduationYear).
		Claim("gpa", gpa).
		ValidFrom(time.Now()).
		SignWith(client.DID(), time.Now()).
		Issue(client)
}

// IsCredentialOfType checks if a credential is of a specific type
func IsCredentialOfType(cred *credential.VerifiableCredential, credentialType []string) bool {
	credType := cred.CredentialType()
	if len(credType) != len(credentialType) {
		return false
	}
	for i, t := range credType {
		if i >= len(credentialType) || t != credentialType[i] {
			return false
		}
	}
	return true
}

// ExtractEmailFromCredential extracts email address from an email credential
func ExtractEmailFromCredential(cred *credential.VerifiableCredential) (string, bool) {
	if !IsCredentialOfType(cred, credential.CredentialTypeEmail) {
		return "", false
	}

	claims, err := cred.CredentialSubjectClaims()
	if err != nil {
		return "", false
	}

	if email, exists := claims["emailAddress"]; exists {
		if emailStr, ok := email.(string); ok {
			return emailStr, true
		}
	}
	return "", false
}

// ExtractNameFromCredential extracts full name from a profile credential
func ExtractNameFromCredential(cred *credential.VerifiableCredential) (string, string, bool) {
	if !IsCredentialOfType(cred, credential.CredentialTypeProfileName) {
		return "", "", false
	}

	claims, err := cred.CredentialSubjectClaims()
	if err != nil {
		return "", "", false
	}

	firstName, firstOk := claims["firstName"].(string)
	lastName, lastOk := claims["lastName"].(string)

	if firstOk && lastOk {
		return firstName, lastName, true
	}
	return "", "", false
}

// ExtractEducationFromCredential extracts education details from an education credential
func ExtractEducationFromCredential(cred *credential.VerifiableCredential) (string, string, int, float64, bool) {
	if !IsCredentialOfType(cred, []string{"VerifiableCredential", "EducationCredential"}) {
		return "", "", 0, 0.0, false
	}

	claims, err := cred.CredentialSubjectClaims()
	if err != nil {
		return "", "", 0, 0.0, false
	}

	degree, degreeOk := claims["degree"].(string)
	institution, instOk := claims["institution"].(string)
	graduationYear, yearOk := claims["graduationYear"].(int)
	gpa, gpaOk := claims["gpa"].(float64)

	if degreeOk && instOk && yearOk && gpaOk {
		return degree, institution, graduationYear, gpa, true
	}
	return "", "", 0, 0.0, false
}
