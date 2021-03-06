package main

import (
	"bytes"
	"fmt"
	"github.com/dermicha/goutils/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	configUtil "github.com/ubirch/ubirch-niomon-proxy/pkg/confutil"
	"github.com/ubirch/ubirch-niomon-proxy/pkg/model/token"
	"net/http"
	"os"
)

var (
	appName    = "Ubirch Niomon Proxy"
	appVersion = "v0.0.1"
	config     = configUtil.GetConfig()
)

func aboutService(c *fiber.Ctx) error {
	err := c.SendString(fmt.Sprintf("%s %s", appName, appVersion))
	return err
}

func InitTokens(c *fiber.Ctx) error {

	sizes := []int{1, 10, 50, 100}
	for _, s := range sizes {
		filename := fmt.Sprintf("authtokens_%d_%s.csv", s, utils.UUID())
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		for i := int32(0); i < (1024); i++ {
			u, _ := uuid.NewV4()

			atok := token.AnkerToken{}
			atok.Token = u.String()

			atok.UsedCounter = s
			db := database.GetDb()
			res := db.Create(&atok)
			if res.Error != nil {
				return res.Error
			}

			_, err := f.WriteString(fmt.Sprintf("%s\n", atok.Token))
			if err != nil {
				f.Close()
				return err
			}
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	err := c.SendString("OK: tokens created")
	return err
}

func PostUpp(c *fiber.Ctx) error {
	xat := c.Get("x-token")
	zuat := c.Get("x-ubirch-auth-type")
	zuhi := c.Get("x-ubirch-hardware-id")
	zuc := c.Get("x-ubirch-credential")
	ct := c.Get("Content-Type")
	url := config.NiomonUrl

	log.Infof("got request from: %s", zuhi)

	if token.IsValidToken(xat) {
		log.Infof("valid token: %s", xat)
		hc := http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(c.Body()))
		req.Header.Add("x-ubirch-auth-type", zuat)
		req.Header.Add("x-ubirch-hardware-id", zuhi)
		req.Header.Add("x-ubirch-credential", zuc)
		req.Header.Add("Content-Type", ct)
		resp, err := hc.Do(req)

		if err != nil {
			err = c.SendStatus(http.StatusBadRequest)
			return err
		}

		if resp.StatusCode < 300 {
			token.UseToken(xat)
		}

		err = c.SendStatus(resp.StatusCode)
		return err
	} else {
		log.Errorf("invalid token: %s", xat)
		err := c.SendStatus(http.StatusPaymentRequired)
		return err
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)

	database.InitDatabase(config.DbFile, nil)
	database.MigrateDatabase(&token.AnkerToken{})

	log.Info("Welcome!")
	app := fiber.New()

	app.Get("/", aboutService)

	apiV1 := app.Group("/ubproxy/api/v1")

	apiV1.Get("/init/", InitTokens)
	apiV1.Post("/upp/", PostUpp)

	err := app.Listen(":3000")
	if err != nil {
		log.Error("could not start server: ", err.Error())
	}
}
