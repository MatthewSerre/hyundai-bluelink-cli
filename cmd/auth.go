package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/authentication/v1"
)

type Auth struct {
	Username   string
	PIN        string
	JWTToken  string
	JWTExpiry int64
}

var Method string

const Authentication_Address = "localhost:50051"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Hyundai",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.OpenFile("token.txt", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("failed to create or open token.txt with error: %v", err)
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		authConn, err := grpc.Dial(Authentication_Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("failed to connect to the authentication service:", err)
			os.Exit(1)
		}
	
		defer authConn.Close()
	
		c := authv1.NewAuthenticationServiceClient(authConn)

		if lines[3] != "" {
			i, _ := strconv.ParseInt(lines[3], 10, 64)
			tm := time.Unix(int64(i), 0)
		
			now:=time.Now()
			if now.After(tm) || now.Equal(tm) {
				log.Println("previous token expired; authenticating again")
				if err := os.Truncate("token.txt", 0); err != nil {
					log.Printf("failed to truncate token.txt with error: %v", err)
				}
				err := authenticate(c, f)
				if err != nil {
					log.Printf("failed to authenticate with error: %v", err)
				}
			} else {
				log.Println("previous token still valid")
			}
		} else {
			err := authenticate(c, f)
			if err != nil {
				log.Printf("failed to authenticate with error: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&Method, "method", "m", "", "Authentication method")
	authCmd.MarkFlagRequired("method")
}

func authenticate(c authv1.AuthenticationServiceClient, f *os.File) (error) {
	var username, password, pin string

	switch Method {
	case "manual":
		username = getInput("username\n")
		log.Println("password")
		password_input, _ := terminal.ReadPassword(0)
		password = string(password_input)
		log.Println("")
		pin = getInput("PIN\n")
	case "env":
		envFile, err := godotenv.Read(".env")
		if err != nil {
			log.Println("failed to read credentials from environment; create a .env file in the root directory and add USERNAME, PASSWORD, and PIN variables with the correct values")
		}
		if err == nil {
			username = envFile["USERNAME"]
			password = envFile["PASSWORD"]
			pin = envFile["PIN"]
		}
	}
	
	res, err := c.Authenticate(context.Background(), &authv1.AuthenticationRequest{
		Username: username,
		Password: password,
		Pin: pin,
	})

	if err != nil {
		return err
	}

	f.WriteString(res.Username + "\n")
	f.WriteString(res.Pin + "\n")
	f.WriteString(res.JwtToken + "\n")
	f.WriteString(strconv.Itoa(int(res.JwtExpiry)))
	log.Println("username, PIN, token, and expiry written to token.txt")

	return nil
}

func getInput(message string) (input string) {
	log.Println(message)
	input_scanner := bufio.NewScanner(os.Stdin)
	input_scanner.Scan()
	return input_scanner.Text()
}

func ReadLine(f *os.File, lineNum int) (line string, lastLine int, err error) {
    sc := bufio.NewScanner(f)
    for sc.Scan() {
        lastLine++
		log.Printf("%v\n", lastLine)
		log.Printf("%v\n", lineNum)
		log.Printf(sc.Text())
        if lastLine == lineNum {
            return sc.Text(), lastLine, sc.Err()
        }
    }
    return line, lastLine, io.EOF
}