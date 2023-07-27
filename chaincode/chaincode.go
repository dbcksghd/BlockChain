package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "log"
    "math/rand"
    "strconv"
)

type SmartContract struct {
    contractapi.Contract
}

type User struct {
	LuckyScore int
	Money      int
	Name       string
	Warning    int
	Ban        bool
}

type StrongBox struct {
	Money int
}

type QueryResult struct {
	Key    string `json:"key"`
	Record *User
}

func isTrueUserBan(user *User) bool {
	if user.Ban == true {
			return true
	}
	return false
}

func (s *SmartContract) Register(ctx contractapi.TransactionContextInterface, name string, money string, id string) error {
	moneyAsInt, _ := strconv.Atoi(money)
	user := User{
			LuckyScore: rand.Intn(50),
			Money:      moneyAsInt,
			Name:       name,
			Warning:    2,
			Ban:        false,
	}
	userAsBytes, err := json.Marshal(user)
	if err != nil {
			log.Fatal(err)
	}
	return ctx.GetStub().PutState(id, userAsBytes)
}

func (s *SmartContract) MakeBank(ctx contractapi.TransactionContextInterface, money string) error {
	strongBox := StrongBox{Money : 1000}
	boxAsBytes, err := json.Marshal(strongBox)
	if err != nil {
			log.Fatal(err)
	}
	return ctx.GetStub().PutState("bank", boxAsBytes)
}

func (s *SmartContract) TurnRoulette(ctx contractapi.TransactionContextInterface, money string, id, box string) (int, error) {
	userAsBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
			return -1, err
	}
	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)
	if isTrueUserBan(user) == true {
			return -1, fmt.Errorf("룰렛 금지")
	}
	boxAsBytes, err := ctx.GetStub().GetState(box)
	if err != nil {
			return -1, err
	}
	strongBox := new(StrongBox)
	_ = json.Unmarshal(boxAsBytes, strongBox)
	moneyAsInt, _ := strconv.Atoi(money)
	strongBox.Money += moneyAsInt
	user.Money -= moneyAsInt
	roulette := 100 - user.LuckyScore
	randomValue := rand.Intn(roulette)
	if randomValue == 1 {
			user.LuckyScore = 0
			strongBox.Money -= user.LuckyScore * 2
			user.Money += user.LuckyScore * 2
			userAsBytes, err = json.Marshal(user)
			boxAsBytes, err = json.Marshal(strongBox)
			_ = ctx.GetStub().PutState(box, boxAsBytes)
			_ = ctx.GetStub().PutState(id, userAsBytes)
			return user.LuckyScore * 2, nil
	} else {
			user.LuckyScore += 5
			userAsBytes, err = json.Marshal(user)
			boxAsBytes, err = json.Marshal(strongBox)
			_ = ctx.GetStub().PutState(box, boxAsBytes)
			_ =  ctx.GetStub().PutState(id, userAsBytes)
			return 0, nil
	}
}

func (s *SmartContract) QueryAllUser(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey, endKey := "", ""
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
			return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
			queryRes, err := resultsIterator.Next()

			if err != nil {
					return nil, err
			}

			user := new(User)
			_ = json.Unmarshal(queryRes.Value, user)

			queryResult := QueryResult{Key: queryRes.Key, Record: user}
			if user.Name != ""{
				results = append(results, queryResult)
			}
	}

	return results, nil
}

func (s *SmartContract) borrowMoney(ctx contractapi.TransactionContextInterface, money string, id, box string) error {
	userAsBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
			return fmt.Errorf(err.Error())
	}
	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)
	if isTrueUserBan(user) == true {
			return fmt.Errorf("돈 못빌림")
	}
	boxAsBytes, err := ctx.GetStub().GetState(box)
	if err != nil {
			return fmt.Errorf(err.Error())
	}
	strongBox := new(StrongBox)
	_ = json.Unmarshal(boxAsBytes, strongBox)
	moneyAsInt, _ := strconv.Atoi(money)
	if moneyAsInt > strongBox.Money {
			return fmt.Errorf("요청한 돈이 매우 큼")
	}
	user.Money += moneyAsInt
	user.Warning--
	boxAsBytes, err = json.Marshal(strongBox)
	_ = ctx.GetStub().PutState(box, boxAsBytes)
	userAsBytes, err = json.Marshal(user)
	return ctx.GetStub().PutState(id, userAsBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
			fmt.Println(err.Error())
			return
	}

	if err = chaincode.Start(); err != nil {
			fmt.Println(err.Error())
	}
}