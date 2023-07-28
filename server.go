package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type User struct {
	LuckyScore int    `json:"luckyScore"`
	Money      int    `json:"money"`
	Name       string `json:"name"`
	Warning    int    `json:"warning"`
	Ban        bool   `json:"ban"`
}

type QueryResult struct {
	Key    string `json:"key"`
	Record *User  `json:"record"`
}

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	e := echo.New()

	e.GET("/user", func(c echo.Context) error {
		contract := a()
		result, err := contract.EvaluateTransaction("QueryAllUser")
		if err != nil {
			return c.JSON(500, err)
		}
		users := []QueryResult{}
		err = json.Unmarshal(result, &users)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		fmt.Println(users)
		return c.JSON(200, users)
	})

	e.POST("/user", func(c echo.Context) error {
		type UserDTO struct {
			Money int    `json:"money"`
			Name  string `json:"name"`
			Id    string `json:"id"`
		}
		newUser := UserDTO{}
		_ = c.Bind(&newUser)
		moneyAsString := strconv.Itoa(newUser.Money)
		contract := a()
		_, err := contract.SubmitTransaction("Register", newUser.Name, moneyAsString, newUser.Id)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(201, map[string]string{"massage": "성공적으로 등록되었습니다"})
	})

	e.GET("/bank", func(c echo.Context) error {
		money := c.QueryParam("money")
		contract := a()
		_, err := contract.SubmitTransaction("MakeBank", money)
		if err != nil {
			return c.JSON(500, err)
		}
		return c.JSON(201, map[string]string{"massage": "성공적으로 등록되었습니다"})
	})

	e.GET("/roulette", func(c echo.Context) error {
		type RouletteDTO struct {
			Money string `json:"money"`
			Box   string `json:"box"`
			Id    string `json:"id"`
		}

		type RouletteRes struct {
			Price int   `json:"price"`
			Error error `json:"error"`
		}
		contract := a()
		newRoulette := RouletteDTO{}
		_ = c.Bind(&newRoulette)
		result, err := contract.SubmitTransaction("TurnRoulette", newRoulette.Money, newRoulette.Id, newRoulette.Box)
		if err != nil {
			return c.JSON(500, err)
		}
		rouletteRes := RouletteRes{}
		err = json.Unmarshal(result, &rouletteRes)
		if err != nil {
			return c.JSON(500, err)
		}
		return c.JSON(200, map[string]string{"price": strconv.Itoa(rouletteRes.Price)})
	})

	e.GET("/borrow", func(c echo.Context) error {
		type BorrowDTO struct {
			Money string `json:"money"`
			Id    string `json:"id"`
		}
		borrowDto := BorrowDTO{}
		_ = c.Bind(&borrowDto)
		contract := a()
		_, err := contract.SubmitTransaction("BorrowMoney", borrowDto.Money, borrowDto.Id)
		if err != nil {
			return c.JSON(500, err)
		}
		result, _ := contract.EvaluateTransaction("QueryUser", borrowDto.Id)
		user := User{}
		_ = json.Unmarshal(result, &user)
		return c.JSON(200, map[string]string{"message": "요청이 성공적으로 완료되었습니다", "잔액": strconv.Itoa(user.Money), "카운트": strconv.Itoa(user.Warning)})
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

func a() *gateway.Contract {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")
	return contract
}
