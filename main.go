package main

import (
	"log"
	"net/http"

	"github.com/farhanswitch/kong-socket/auth"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

type Config struct {
	AuthUrl string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access(kong *pdk.PDK) {
	var staticToken string = "da94d282404dc7870379ce118739b907c3e9663e5f4ed0528ec2908919e82ca4"
	token, err := kong.Request.GetQueryArg("t")

	if err != nil {
		log.Printf("There is no token")
		// Berikan response 401
		kong.Response.Exit(http.StatusUnauthorized, `{"status":"Login first!"}`, map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	if token == "" {
		log.Printf("There is no token")
		// Berikan response 401
		kong.Response.Exit(http.StatusUnauthorized, `{"status":"Login first!"}`, map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	if token == staticToken {
		kong.ServiceRequest.AddHeader("X-Room-Socket", "sendPrice")
	} else {
		// Buat instance dari Auth Provider
		authProvider := auth.AuthProviderFactory(c.AuthUrl)
		// Hit Service Authentikasi untuk verify apakah user dengan token tersebut diperbolehkan access ke resource yang akan dituju
		room, err := authProvider.AuthRequest(token)
		if err != nil {
			log.Printf("Auth Provider Failed\n")
			kong.Response.Exit(http.StatusForbidden, `{"status":"Access Denied!"}`, map[string][]string{
				"Content-Type": {"application/json"},
			})
			return
		}
		kong.ServiceRequest.AddHeader("X-Room-Socket", room)

	}

}
func main() {
	Version := "1.0"
	Priority := 1000
	_ = server.StartServer(New, Version, Priority)
}
