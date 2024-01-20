package packages

import (
	"encoding/json"
	"fmt"
	"os"
)

type Account struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

func Json() {
	fmt.Println("Marshalling:")
	acc := Account{
		Owner:   "Buddhi",
		Balance: 0.10,
	}
	res, err := json.Marshal(acc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Marshalled: from %v (type: %T) to %v (binary form) or %v (binary to string)\n", acc, acc, res, string(res))

	fmt.Println("\nMarshalling directly as string to /sample-dir/new-file.json...")

	file, err := os.Create("./sample-dir/new-file.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(acc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Marshalled to file successfully!")

	fmt.Println("Marshalling to Stdout:")
	json.NewEncoder(os.Stdout).Encode(acc)

	fmt.Println("\nUnmarshalling:")
	jsonString := `{"owner":"Dracula","balance":100000.00}`
	var acc2 Account
	err = json.Unmarshal([]byte(jsonString), &acc2)
	fmt.Printf("Unmarshalled from string: %v\n", jsonString)
	fmt.Printf("To: %v (type: %T)\n", acc2, acc2)
}
