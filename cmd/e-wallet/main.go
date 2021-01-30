package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"linkaja.com/e-wallet/config/postgresql"
	accountManagerHTTPHandler "linkaja.com/e-wallet/pkg/account-manager/handler"
	accountManagerRepo "linkaja.com/e-wallet/pkg/account-manager/repository"
	accountManagerUC "linkaja.com/e-wallet/pkg/account-manager/usecase"
)

func main() {
	// if err := config.Load(".env"); err != nil {
	// 	fmt.Println(".env is not loaded properly")
	// 	fmt.Println(err)
	// 	os.Exit(2)
	// }
	r := echo.New()
	r.Debug = true
	r.Use(middleware.Recover())
	r.Use(middleware.Logger())
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	r.GET("/", func(context echo.Context) error {
		return context.HTML(http.StatusOK, "<strong>E-WALLET</strong>")
	})

	dbUser := postgresql.CreateDBConnection("postgres", os.Getenv("DB_USER_URL"), os.Getenv("MAX_CONNECTION_POOL"))
	accountManagerRepo := accountManagerRepo.NewAccountManagerRepo(dbUser)
	accountManagerUC := accountManagerUC.NewAccountManagerUsecase(accountManagerRepo)
	accountManagerHTTPHandler := accountManagerHTTPHandler.NewAccountManagerHandler(accountManagerUC)
	accountManagerGroup := r.Group("/v1/e-wallet")
	accountManagerHTTPHandler.Mount(accountManagerGroup)

	err := r.Start(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err.Error())
	}
}
