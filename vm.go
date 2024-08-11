package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a basic transaction in the blockchain
type Transaction struct {
	ID       string
	Sender   *Account
	Receiver *Account
	Amount   float64
}

// NewTransaction creates a new transaction and generates its ID
func NewTransaction(sender, receiver *Account, amount float64) *Transaction {
	tx := &Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
	tx.ID = tx.hashTransaction()
	return tx
}

// hashTransaction generates a hash ID for the transaction
func (tx *Transaction) hashTransaction() string {
	record := tx.Sender.Username + tx.Receiver.Username + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Block represents a block in the blockchain
type Block struct {
	Timestamp     time.Time
	Transactions  []*Transaction
	PrevBlockHash string
	Hash          string
}

// Blockchain represents the entire chain
type Blockchain struct {
	Blocks []*Block
}

// NewBlock creates a new block containing transactions
func NewBlock(transactions []*Transaction, prevBlockHash string) *Block {
	block := &Block{
		Timestamp:     time.Now(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
	block.Hash = block.hashBlock()
	return block
}

// hashBlock generates a hash for the block
func (b *Block) hashBlock() string {
	record := fmt.Sprintf("%s%s", b.Timestamp.String(), b.PrevBlockHash)
	for _, tx := range b.Transactions {
		record += tx.ID
	}
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock([]*Transaction{}, "")
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(transactions, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

type VirtualMachine struct {
	Blockchain *Blockchain
	Accounts   map[string]*Account
}

// NewVirtualMachine initializes a new VM with an empty blockchain and account map
func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		Blockchain: NewBlockchain(),
		Accounts:   make(map[string]*Account),
	}
}

// CreateAccount creates a new account with the given username
func (vm *VirtualMachine) CreateAccount(username string) *Account {
	if _, exists := vm.Accounts[username]; exists {
		fmt.Printf("Account with username %s already exists.\n", username)
		return nil
	}
	account := NewAccount(username)
	vm.Accounts[username] = account
	fmt.Printf("Account created: %s\n", username)
	return account
}

// ProcessTransaction handles a single transaction
func (vm *VirtualMachine) ProcessTransaction(tx *Transaction) {
	fmt.Printf("Processing Transaction: ID=%s, From=%s, To=%s, Amount=%.2f\n",
		tx.ID, tx.Sender.Username, tx.Receiver.Username, tx.Amount)
	// In a real system, we would update balances, etc.
}

// ExecuteBlock processes all transactions in a block
func (vm *VirtualMachine) ExecuteBlock(block *Block) {
	for _, tx := range block.Transactions {
		vm.ProcessTransaction(tx)
	}
}

// AddBlockToChain adds a block to the blockchain and processes it
func (vm *VirtualMachine) AddBlockToChain(transactions []*Transaction) {
	vm.Blockchain.AddBlock(transactions)
	vm.ExecuteBlock(vm.Blockchain.Blocks[len(vm.Blockchain.Blocks)-1])
}

// GetAccount retrieves an account by username
func (vm *VirtualMachine) GetAccount(username string) *Account {
	return vm.Accounts[username]
}

func main() {
	vm := NewVirtualMachine()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nCommands:")
		fmt.Println("1. create_account [username]")
		fmt.Println("2. send [sender] [receiver] [amount]")
		fmt.Println("3. view_blockchain")
		fmt.Println("4. exit")

		fmt.Print("Enter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		parts := strings.Split(command, " ")

		switch parts[0] {
		case "create_account":
			if len(parts) != 2 {
				fmt.Println("Usage: create_account [username]")
			} else {
				vm.CreateAccount(parts[1])
			}

		case "send":
			if len(parts) != 4 {
				fmt.Println("Usage: send [sender] [receiver] [amount]")
			} else {
				sender := vm.GetAccount(parts[1])
				receiver := vm.GetAccount(parts[2])
				if sender == nil || receiver == nil {
					fmt.Println("Invalid sender or receiver.")
					break
				}
				amount, err := strconv.ParseFloat(parts[3], 64)
				if err != nil {
					fmt.Println("Invalid amount.")
					break
				}
				tx := NewTransaction(sender, receiver, amount)
				vm.AddBlockToChain([]*Transaction{tx})
			}

		case "view_blockchain":
			viewBlockchain(vm)

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

// viewBlockchain prints the entire blockchain
func viewBlockchain(vm *VirtualMachine) {
	for i, block := range vm.Blockchain.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Previous Hash: %s\n", block.PrevBlockHash)
		for _, tx := range block.Transactions {
			fmt.Printf("  TxID: %s | From: %s | To: %s | Amount: %.2f\n", tx.ID, tx.Sender.Username, tx.Receiver.Username, tx.Amount)
		}
	}
}

// Account represents a user account with just a username
type Account struct {
	Username string
}

// NewAccount creates a new account with the given username
func NewAccount(username string) *Account {
	return &Account{
		Username: username,
	}
}
