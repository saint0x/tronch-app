package interactions

import (
  "context"
  "log"
  "math/big"

  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/ethclient"
  "smart_contract/pkg/smart-contract/contract" // Import the compiled contract binding
)

// ContractAddress should be the address where your deployed contract resides
var ContractAddress = common.HexToAddress("YOUR_CONTRACT_ADDRESS_HERE")

// Interactor provides functionalities to trigger interactions with the smart contract
type Interactor struct {
  ethClient *ethclient.Client
}

// NewInteractor creates a new instance of Interactor
func NewInteractor(nodeURL string) (*Interactor, error) {
  client, err := ethclient.Dial(nodeURL)
  if err != nil {
    return nil, err
  }
  return &Interactor{
    ethClient: client,
  }, nil
}

// bindNewTransactor creates a new transactor
func bindNewTransactor() (*bind.TransactOpts, error) {
  privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY_HERE")
  if err != nil {
    return nil, err
  }

  auth := bind.NewKeyedTransactor(privateKey)
  return auth, nil
}



// ExecuteContract triggers the deployment and execution of the generated smart contract
func (i *Interactor) ExecuteContract(ctx context.Context, contractCodeFilePath string) (string, error) {
  // Read the contract code from the provided file path
  contractCode, err := os.ReadFile(contractCodeFilePath)
  if err != nil {
    return "", fmt.Errorf("failed to read contract code: %v", err)
  }

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return "", fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return "", fmt.Errorf("failed to parse private key: %v", err)
  }

  // Deploy the smart contract
  auth := bind.NewKeyedTransactor(adminPrivateKey)
  contractAddress, _, _, err := deployContract(auth, client, contractCode)
  if err != nil {
    return "", fmt.Errorf("failed to deploy contract: %v", err)
  }

  log.Printf("Smart contract deployed at address: %s", contractAddress.Hex())

  // Execute the smart contract
  contract, err := NewEscrowService(contractAddress, client)
  if err != nil {
    return "", fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Call the executeContract function
  tx, err := contract.ExecuteContract(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  })
  if err != nil {
    return "", fmt.Errorf("failed to execute contract: %v", err)
  }

  log.Printf("Executing contract. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return "", fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  return contractAddress.Hex(), nil
}

// deployContract deploys the smart contract to the Polygon chain
func deployContract(auth *bind.TransactOpts, client *ethclient.Client, contractCode []byte) (common.Address, *types.Transaction, *EscrowService, error) {
  // Compile and deploy the contract
  parsedContract, err := solidity.ABI("EscrowService", string(contractCode))
  if err != nil {
    return common.Address{}, nil, nil, fmt.Errorf("failed to parse contract ABI: %v", err)
  }

  address, tx, contract, err := bind.DeployContract(auth, parsedContract, client, parsedContract.Constructor.Inputs[0].Args[0].Value.(*bind.TransactOpts).From)
  if err != nil {
    return common.Address{}, nil, nil, fmt.Errorf("failed to deploy contract: %v", err)
  }

  return address, tx, contract, nil
}




// MarkRequirementsComplete triggers the interaction to mark requirements as complete
func (i *Interactor) MarkRequirementsComplete(ctx context.Context, contractAddress string) error {
  // Implement interaction to mark requirements as complete
  log.Println("Marking requirements as complete...")

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Load the smart contract
  contract, err := NewEscrowService(common.HexToAddress(contractAddress), client)
  if err != nil {
    return fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return fmt.Errorf("failed to parse private key: %v", err)
  }

  // Call the markRequirementsComplete function
  tx, err := contract.MarkRequirementsComplete(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  })
  if err != nil {
    return fmt.Errorf("failed to mark requirements as complete: %v", err)
  }

  log.Printf("Marking requirements as complete. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  // Now, let's confirm the requirements
  return i.ConfirmReqs(ctx, contractAddress)
}




// ConfirmReqs triggers the interaction for the buyer to confirm the requirements
func (i *Interactor) ConfirmReqs(ctx context.Context, contractAddress string) error {
  // TODO: Implement interaction for buyer to confirm requirements
  log.Println("Confirming requirements...")

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Load the smart contract
  contract, err := NewEscrowService(common.HexToAddress(contractAddress), client)
  if err != nil {
    return fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return fmt.Errorf("failed to parse private key: %v", err)
  }

  // Call the confirmReqs function
  tx, err := contract.ConfirmReqs(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  })
  if err != nil {
    return fmt.Errorf("failed to confirm requirements: %v", err)
  }

  log.Printf("Confirming requirements. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  return nil
}



// InitiateDispute triggers the interaction to initiate a dispute
func (i *Interactor) InitiateDispute(ctx context.Context, contractAddress string) error {
  // TODO: Implement interaction to initiate a dispute
  log.Println("Initiating dispute...")

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Load the smart contract
  contract, err := NewEscrowService(common.HexToAddress(contractAddress), client)
  if err != nil {
    return fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return fmt.Errorf("failed to parse private key: %v", err)
  }

  // Call the initiateDispute function
  tx, err := contract.InitiateDispute(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  })
  if err != nil {
    return fmt.Errorf("failed to initiate dispute: %v", err)
  }

  log.Printf("Initiating dispute. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  return nil
}



// ResolveDispute triggers the interaction to resolve a dispute
func (i *Interactor) ResolveDispute(ctx context.Context, contractAddress string, resolution string) error {
  // TODO: Implement interaction to resolve a dispute
  log.Printf("Resolving dispute with resolution: %s...", resolution)

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Load the smart contract
  contract, err := NewEscrowService(common.HexToAddress(contractAddress), client)
  if err != nil {
    return fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return fmt.Errorf("failed to parse private key: %v", err)
  }

  // Call the resolveDispute function
  tx, err := contract.ResolveDispute(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  }, resolution)
  if err != nil {
    return fmt.Errorf("failed to resolve dispute: %v", err)
  }

  log.Printf("Resolving dispute. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  return nil
}



// UpdateContractProgress triggers the interaction to update the contract's progression
func (i *Interactor) UpdateContractProgress(ctx context.Context, contractAddress string, progressStatus uint8) error {
  // TODO: Implement interaction to update contract's progression
  log.Printf("Updating contract progression to status: %d...", progressStatus)

  // Initialize the Ethereum client
  client, err := ethclient.Dial("https://polygon-mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
  if err != nil {
    return fmt.Errorf("failed to connect to Ethereum client: %v", err)
  }

  // Load the smart contract
  contract, err := NewEscrowService(common.HexToAddress(contractAddress), client)
  if err != nil {
    return fmt.Errorf("failed to instantiate smart contract: %v", err)
  }

  // Admin address and private key (this should be stored securely and not hardcoded)
  adminAddress := common.HexToAddress("YOUR_ADMIN_ADDRESS")
  adminPrivateKey, err := crypto.HexToECDSA("YOUR_ADMIN_PRIVATE_KEY")
  if err != nil {
    return fmt.Errorf("failed to parse private key: %v", err)
  }

  

  // Call the updateContractProgress function
  tx, err := contract.UpdateContractProgress(&bind.TransactOpts{
    From: adminAddress,
    Signer: func(signer bind.SignerFn, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
      return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(137)), adminPrivateKey)
    },
  }, progressStatus)
  if err != nil {
    return fmt.Errorf("failed to update contract progression: %v", err)
  }

  log.Printf("Updating contract progression. Transaction hash: %s", tx.Hash().Hex())

  // Wait for the transaction to be mined
  receipt, err := bind.WaitMined(ctx, client, tx)
  if err != nil {
    return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
  }

  log.Printf("Transaction mined. Receipt: %v", receipt)

  return nil
}

