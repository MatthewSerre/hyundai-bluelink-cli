/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"time"

	infov1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/information/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Vehicle struct {
	RegistrationID string
	VIN string
	Generation string
	Mileage string
}

const InformationAddress = "localhost:50052"

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Obtain information about a vehicle",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.OpenFile("token.txt", os.O_RDONLY, 0644)
		if err != nil {
			log.Printf("failed to open token.txt with error: %v", err)
		}
		defer f.Close()

		auth, err := GetAuth(f)
		if err != nil {
			log.Fatalf("failed to get auth from token.txt with error; run hb auth --method [env|manual]: %v", err)
		}

		if auth.JWTExpiry != 0 {
			tm := time.Unix(int64(auth.JWTExpiry), 0)
		
			now:=time.Now()
			if now.After(tm) || now.Equal(tm) {
				log.Println("previous token expired; run hb auth --method [env|manual]")
				os.Exit(1)
			}
		}
	
		infoConn, err := grpc.Dial(InformationAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("failed to connect to the information service: ", err)
			os.Exit(1)
		}
	
		defer infoConn.Close()
	
		d := infov1.NewInformationServiceClient(infoConn)

		vehicle, err := getVehicleInfo(d, auth)
		if err != nil {
			log.Printf("failed to get vehicle info with error: %v", err)
			os.Exit(1)
		}

		g, err := os.OpenFile("info.txt", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("failed to create or open token.txt with error: %v", err)
		}
		defer g.Close()

		if err := os.Truncate("info.txt", 0); err != nil {
			log.Printf("failed to truncate info.txt with error: %v", err)
		}

		g.WriteString(vehicle.RegistrationID + "\n")
		g.WriteString(vehicle.VIN+ "\n")
		g.WriteString(vehicle.Generation + "\n")
		g.WriteString(vehicle.Mileage)

		log.Println("Registration ID:", vehicle.RegistrationID)
		log.Println("VIN:", vehicle.VIN)
		log.Println("Mileage:", vehicle.Mileage)
		log.Println("Generation:", vehicle.Generation)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func getVehicleInfo(c infov1.InformationServiceClient, auth Auth) (Vehicle, error) {
	res, err := c.GetVehicleInfo(context.Background(), &infov1.VehicleInfoRequest{
		Username: auth.Username,
		Pin: auth.PIN,
		JwtToken: auth.JWTToken,
		JwtExpiry: auth.JWTExpiry,
	})

	if err != nil {
		return Vehicle{}, err
	}

	return Vehicle{ RegistrationID: res.RegistrationId, VIN: res.Vin, Generation: res.Generation, Mileage: res.Mileage }, nil
}

func GetVehicleInfo(f *os.File) (Vehicle, error) {
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != 4 {
		err := errors.New("emit macho dwarf: elf header corrupted")
		return Vehicle{}, err
	}

	return Vehicle{RegistrationID: lines[0], VIN: lines[1], Generation: lines[2], Mileage: lines[3]}, nil
}