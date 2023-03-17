/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	remoteActionV1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/remote_action/v1"
	remote_action_v1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/remote_action/v1"
)

type RemoteActionResponse struct {
	Result string
	FailMsg string
	ResponseString *remoteActionV1.ResponseString
}

const RemoteActionAddress = "localhost:50053"

var Toggle string

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.OpenFile("token.txt", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open token.txt with error: %v", err)
		}
		defer f.Close()

		auth, err := GetAuth(f)
		if err != nil {
			log.Fatalf("failed to get auth from token.txt with error; run hb auth --method [env|manual]: %v", err)
		}

		g, err := os.OpenFile("info.txt", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open info.txt with error: %v", err)
		}
		defer g.Close()

		info, err := GetVehicleInfo(g)
		if err != nil {
			log.Fatalf("failed to get vehicle info from info.txt with error; run hb auth --method [env|manual]: %v", err)
		}

		remoteActionConn, err := grpc.Dial(RemoteActionAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Fatalf("failed to connect to the information service: %v", err)
		}
	
		defer remoteActionConn.Close()

		e := remoteActionV1.NewRemoteActionServiceClient(remoteActionConn)

		var remoteAction remote_action_v1.LockAction

		switch Toggle {
		case "lock":
			remoteAction = remoteActionV1.LockAction_LOCK_ACTION_LOCK
		case "unlock":
			remoteAction = remoteActionV1.LockAction_LOCK_ACTION_UNLOCK
		}

		toggleLock(e, auth, info, remoteAction)

		log.Printf("vehicle %ved", Toggle)
		log.Printf("sleeping for 60 seconds to buffer requests")
		time.Sleep(60 * time.Second)
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)

	lockCmd.Flags().StringVarP(&Toggle, "toggle", "t", "", "Lock or unlock")
	lockCmd.MarkFlagRequired("toggle")
}

func toggleLock(c remoteActionV1.RemoteActionServiceClient, auth Auth, vehicle Vehicle, lockAction remoteActionV1.LockAction) (RemoteActionResponse, error) {
	res, err := c.ToggleLock(context.Background(), &remoteActionV1.ToggleLockRequest{
		Username: auth.Username,
		Pin: auth.PIN,
		JwtToken: auth.JWTToken,
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