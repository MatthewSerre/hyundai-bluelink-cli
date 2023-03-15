package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/authentication/v1"
	infov1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/information/v1"
	remote_actionv1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/remote_action/v1"
	"github.com/joho/godotenv"
)

type Auth struct {
	Username   string
	PIN        string
	JWT_Token  string
	JWT_Expiry int64
}

type Vehicle struct {
	RegistrationID string
	VIN string
	Generation string
	Mileage string
}

type RemoteActionResponse struct {
	Result string
	FailMsg string
	ResponseString *remote_actionv1.ResponseString
}

const Authentication_Address = "localhost:50051"
const Information_Address = "localhost:50052"
const Remote_Action_Address = "localhost:50053"

func main() {
	log.Println("Welcome to the unofficial Hyundai Bluelink CLI!")

	log.Println("Establishing connection to the authentication service...")

	authConn, err := grpc.Dial(Authentication_Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("failed to connect to the authentication service:", err)
		log.Println("Shutting down...")
		os.Exit(1)
	}

	defer authConn.Close()

	c := authv1.NewAuthenticationServiceClient(authConn)

	var authStub Auth;

	exit := false;
	for !exit {
		auth, err := authenticate(c)

		exit = true

		if err != nil {
			log.Println("authentication failed with error:", err)
			exit = false
		}

		if (Auth{}) == auth {
			log.Println("authentication failed")
			exit = false
		}

		if err == nil && auth != (Auth{}) {
			log.Println("Authentication successful!")
			authStub = auth
		}
	}

	log.Println("Establishing connection to the information service...")

	infoConn, err := grpc.Dial(Information_Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("failed to connect to the information service:", err)
		os.Exit(1)
	}

	defer infoConn.Close()

	d := infov1.NewInformationServiceClient(infoConn)

	log.Println("Obtaining vehicle information...")
		
	info, err := getVehicleInfo(d, authStub)

	if err != nil {
		log.Println("vehicle information request failed with error:", err)
		os.Exit(1)
	}

	remoteActionConn, err := grpc.Dial(Remote_Action_Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	e := remote_actionv1.NewRemoteActionServiceClient(remoteActionConn)

	exitMenu := false;
	for !exitMenu {
		var command string
		log.Println("Select an option from the list below and enter the corresponding number. Otherwise, enter 'exit' to close the application.")
		log.Println("1. Display vehicle registration ID, VIN, and mileage)")
		log.Println("2. Lock vehicle")
		log.Println("3. Unlock vehicle")
		fmt.Scan(&command)
		switch command {
		case "1":
			log.Println("Registration ID:", info.RegistrationID)
			log.Println("VIN:", info.VIN)
			log.Println("Mileage:", info.Mileage)
		case "2":
			_, err := toggleLock(e, authStub, info, remote_actionv1.LockAction_LOCK_ACTION_LOCK)

			if err != nil {
				log.Println("vehicle lock request failed with error:", err)
			}

			if err == nil {
				log.Println("No error returned from the server. Sleeping for three minutes to buffer requests...")
				time.Sleep(180 * time.Second)
				log.Println("Vehicle locked!")
			}
		case "3":
			_, err := toggleLock(e, authStub, info, remote_actionv1.LockAction_LOCK_ACTION_UNLOCK)

			if err != nil {
				log.Println("vehicle unlock request failed with error:", err)
			}

			if err == nil {
				log.Println("No error returned from the server. Sleeping for three minutes to buffer requests...")
				time.Sleep(180 * time.Second)
				log.Println("Vehicle unlocked!")
			}
		case "exit":
			exitMenu = true
		default:
			continue
		}
	}

	fmt.Println("Thank you for using the unofficial Hyundai Bluelink CLI!")
}


func authenticate(c authv1.AuthenticationServiceClient) (Auth, error) {
	var username, password, pin string

	exit := false;
	for !exit {
		var command string
		log.Println("Enter 1 to input your credentials or 2 to have them read from the environment.")
		fmt.Scan(&command)
		switch command {
		case "1":
			username = getInput("Enter your username")
			log.Println("Enter your password")
			password_input, _ := terminal.ReadPassword(0)
			password = string(password_input)
			log.Println("\n")
			pin = getInput("Enter your PIN")
			exit = true
		case "2":
			envFile, err := godotenv.Read(".env")
			if err != nil {
				log.Println("Please create a .env file in the root directory and add USERNAME, PASSWORD, and PIN variables with the correct values.")
			}
			if err == nil {
				username = envFile["USERNAME"]
				password = envFile["PASSWORD"]
				pin = envFile["PIN"]
				exit = true
			}
		default:
			continue
		}
	}

	log.Println("Authenticating!")
	
	res, err := c.Authenticate(context.Background(), &authv1.AuthenticationRequest{
		Username: username,
		Password: password,
		Pin: pin,
	})

	if err != nil {
		return Auth{}, err
	}

	return Auth{ Username: res.Username, PIN: res.Pin, JWT_Token: res.JwtToken, JWT_Expiry: res.JwtExpiry }, nil
}

func getInput(message string) (input string) {
	log.Println(message)
	input_scanner := bufio.NewScanner(os.Stdin)
	input_scanner.Scan()
	return input_scanner.Text()
}

func getVehicleInfo(c infov1.InformationServiceClient, auth Auth) (Vehicle, error) {
	res, err := c.GetVehicleInfo(context.Background(), &infov1.VehicleInfoRequest{
		Username: auth.Username,
		Pin: auth.PIN,
		JwtToken: auth.JWT_Token,
		JwtExpiry: auth.JWT_Expiry,
	})

	if err != nil {
		return Vehicle{}, err
	}

	return Vehicle{ RegistrationID: res.RegistrationId, VIN: res.Vin, Generation: res.Generation, Mileage: res.Mileage }, nil
}

func toggleLock(c remote_actionv1.RemoteActionServiceClient, auth Auth, vehicle Vehicle, lockAction remote_actionv1.LockAction) (RemoteActionResponse, error) {
	res, err := c.ToggleLock(context.Background(), &remote_actionv1.ToggleLockRequest{
		Username: auth.Username,
		Pin: auth.PIN,
		JwtToken: auth.JWT_Token,
		RegistrationId: vehicle.RegistrationID,
		Vin: vehicle.VIN,
		Generation: vehicle.Generation,
		LockAction: lockAction,
	})

	if err != nil {
		return RemoteActionResponse{}, err
	}

	return RemoteActionResponse{Result: res.Result, FailMsg: res.FailMsg, ResponseString: res.ResponseString}, nil
}