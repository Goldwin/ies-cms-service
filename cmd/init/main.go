package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Goldwin/ies-pik-cms/internal/config"
	authData "github.com/Goldwin/ies-pik-cms/internal/data/auth"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

type output struct {
}

func (o output) OnError(err out.AppErrorDetail) {
	log.Fatal(err.Error())
}

func (o output) OnSuccess(result dto.AuthData) {
	log.Default().Printf("Successfully granted %s admin access", result.Email)
}

func main() {
	config := config.LoadConfigEnv()

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	authDataLayer := authData.NewAuthDataLayerComponent(config.DataConfig["AUTH"], infraComponent)
	authComponent := auth.NewAuthComponent(authDataLayer, config.Secret)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username to be granted admin access: ")
	text, _ := reader.ReadString('\n')

	output := output{}
	authComponent.GrantAdminRole(context.Background(), strings.TrimSpace(text), output).Wait()
}
