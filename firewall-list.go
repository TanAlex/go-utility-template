package main

import (
	"context"
	"os"
	"strings"

	"encoding/csv"
	"flag"
	"fmt"
	"log"

	// "os"
	"golang.org/x/oauth2"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func printUsage() {
	fmt.Println("Usage: firewall-list -projectID <projectID> -output <outputFile>")
}

func main() {
	// projectID := "<your-project-id>"
	// Define command-line flags.
	projectP := flag.String("projectID", "", "GCP project ID")
	// vpcName := flag.String("vpcName", "", "Name of VPC network to get firewalls for")
	// credsPath := flag.String("credsPath", "", "Path to GCP credentials JSON file")
	outputFileP := flag.String("output", "output.csv", "Output file name")

	// Parse command-line flags.
	flag.Parse()
	projectID := *projectP
	// fmt.Printf("projectID: %s", projectID)
	// If project ID is not provided, print usage and exit.
	if projectID == "" {
		printUsage()
		os.Exit(1)
	}

	// If credentials file path is provided, use it to set the credentials.
	// Otherwise, use the default credentials.
	var outputFile string
	if *outputFileP != "" {
		outputFile = *outputFileP
	} else {
		outputFile = "output.csv"
	}

	var client *compute.Service
	var err error
	var ctx context.Context

	if accessToken := os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"); accessToken != "" {
		ctx = context.Background()
		// Create an OAuth2 token source using the access token.
		token := &oauth2.Token{
			AccessToken: os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"),
		}
		tokenSource := oauth2.StaticTokenSource(token)
		// Use the token source to create an HTTP client.
		oauthClient := oauth2.NewClient(ctx, tokenSource)
		client, err = compute.NewService(ctx, option.WithHTTPClient(oauthClient))
		if err != nil {
			log.Fatalf("Failed to create GCP client: %v", err)
		}
	} else {
		ctx = context.Background()

		// Use the default credentials
		client, err = compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
		if err != nil {
			log.Fatal(err)
		}
		// defer client.Close()
	}

	// List all the firewalls in the specified project and zone
	firewalls, err := client.Firewalls.List(projectID).Do()
	if err != nil {
		log.Fatal(err)
	}

	// // Print the details of each firewall
	// for _, firewall := range firewalls.Items {
	// 	fmt.Printf("Firewall name: %s\n", firewall.Name)
	// 	fmt.Printf("Firewall description: %s\n", firewall.Description)
	// 	fmt.Printf("Firewall network: %s\n", firewall.Network)

	// 	fmt.Println("Allowed rules:")
	// 	for _, allowed := range firewall.Allowed {
	// 		fmt.Printf("  Protocol: %s\n", allowed.IPProtocol)
	// 		fmt.Printf("  Ports: %s\n", allowed.Ports)
	// 	}

	// 	fmt.Println("Denied rules:")
	// 	for _, denied := range firewall.Denied {
	// 		fmt.Printf("  Protocol: %s\n", denied.IPProtocol)
	// 		fmt.Printf("  Ports: %s\n", denied.Ports)
	// 	}

	// 	fmt.Println()
	// }

	// Write firewalls to CSV file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	headers := []string{
		"Name", "Description", "Network", "Disabled",
		"SourceRanges", "SourceServiceAccounts", "SourceTags",
		"TargetTags", "TargetServiceAccounts", "Allowed", "Denied",
	}
	writer.Write(headers)
	for _, firewall := range firewalls.Items {
		ruleDetails := []string{firewall.Name, firewall.Description, firewall.Network}
		// Show all the attributes of the firewall
		// log.Printf("%+v", firewall)
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.Disabled))
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.SourceRanges))
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.SourceServiceAccounts))
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.SourceTags))
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.TargetTags))
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", firewall.TargetServiceAccounts))

		// Add allowed rules to the details slice
		allowedRules := []string{}
		for _, allowed := range firewall.Allowed {
			protocol := allowed.IPProtocol
			ports := ""
			if allowed.Ports != nil {
				// join the allowed.Ports into a string
				// show how to join a slice of strings into a single string
				// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
				ports = strings.Join(allowed.Ports, ";")
				// ports = fmt.Sprintf("%v", allowed.Ports)
			}
			allowedRules = append(allowedRules, fmt.Sprintf("%s:%s", protocol, ports))
		}
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", allowedRules))

		// Add denied rules to the details slice
		deniedRules := []string{}
		for _, denied := range firewall.Denied {
			protocol := denied.IPProtocol
			ports := ""
			if denied.Ports != nil {
				ports = strings.Join(denied.Ports, ";")
			}
			deniedRules = append(deniedRules, fmt.Sprintf("%s:%s", protocol, ports))
		}
		ruleDetails = append(ruleDetails, fmt.Sprintf("%v", deniedRules))
		fmt.Println(strings.Join(ruleDetails, ","))
		writer.Write(ruleDetails)
	}

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Firewall rules written to %s", outputFile)
}
