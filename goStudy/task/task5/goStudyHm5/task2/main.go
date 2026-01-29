package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"task2/nftGo" // 注意：这里要替换为你的 NftGo 绑定代码实际包路径

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 配置参数 - 请替换为你自己的信息
const (
	// Sepolia 测试网 RPC 端点（推荐使用 Alchemy/Infura）
	//sepoliaRPC = "https://sepolia.infura.io/v3/8d91817a59984be3b84d4a3cdc950ab3"
	sepoliaRPC = "http://127.0.0.1:8545/"
	// 你的钱包私钥（不含 0x 前缀，仅用于测试网！）
	//privateKeyStr = ""
	privateKeyStr = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// 部署后的合约地址（部署成功后填写，用于直接交互）
	//contractAddress = "0xA62e291C2AAd58EC9fad6f8ED2FF2c4648B20aB6"
	contractAddress = ""
)

func main() {
	// 1. 连接 Sepolia 测试网
	client, err := ethclient.Dial(sepoliaRPC)
	if err != nil {
		//log.Fatalf("❌ 连接 Sepolia 失败: %v", err)
		log.Fatalf("❌ 连接 本地部署hardhat节点 失败: %v", err)
	}
	defer client.Close()
	// fmt.Println("✅ 成功连接 Sepolia 测试网")
	fmt.Println("✅ 成功连接 本地部署hardhat 测试网")

	// 2. 加载私钥并创建交易签名器
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("❌ 解析私钥失败: %v", err)
	}

	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取链 ID 失败: %v", err)
	}

	// 创建交易授权对象
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ 创建交易签名器失败: %v", err)
	}

	// // ========== 修复Gas配置（关键） ==========
	// // 获取最新nonce（避免冲突）
	// nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	// if err != nil {
	// 	log.Fatalf("❌ 获取Nonce失败: %v", err)
	// }
	// log.Fatalf("获取Nonce: %d", nonce)
	// auth.Nonce = big.NewInt(int64(nonce))

	// // EIP-1559 Gas 设置（Sepolia推荐）
	// auth.GasLimit = 6000000 // 进一步提高Gas上限
	// suggestedGasTip, err := client.SuggestGasTipCap(context.Background())
	// if err != nil {
	// 	log.Fatalf("❌ 获取小费失败: %v", err)
	// }
	// // 手动提高小费（确保被矿工优先打包）
	// auth.GasTipCap = new(big.Int).Mul(suggestedGasTip, big.NewInt(3))
	// // 设置基础费上限
	// auth.GasFeeCap = new(big.Int).Mul(auth.GasTipCap, big.NewInt(3))
	// fmt.Printf("🔧 Gas配置 - Tip: %d Gwei, FeeCap: %d Gwei\n",
	// 	new(big.Int).Div(auth.GasTipCap, big.NewInt(1e9)),
	// 	new(big.Int).Div(auth.GasFeeCap, big.NewInt(1e9)))

	// 配置 Gas 参数（根据网络动态调整）
	auth.GasLimit = 5000000 // NFT 操作 Gas 上限（足够覆盖铸造/转账）
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取建议 Gas 价格失败: %v", err)
	}
	auth.GasPrice = gasPrice

	// 3. 部署合约 或 连接已部署合约
	var nftContract *nftGo.NftGo
	var deployedAddr common.Address

	if contractAddress == "" {
		// 3.1 部署新合约
		fmt.Println("\n🚀 开始部署 NftGo 合约...")
		// 部署参数：name, symbol, baseTokenURI
		_, tx, deployContract, err := nftGo.DeployNftGo(
			auth,
			client,
			"MyTestNFT",                         // NFT 名称
			"MTNFT",                             // NFT 符号
			"https://example.com/nft/metadata/", // 基础元数据 URI
		)
		if err != nil {
			log.Fatalf("❌ 部署合约失败: %v", err)
		}

		// 等待部署交易确认
		fmt.Printf("📤 部署交易哈希: %s\n", tx.Hash().Hex())
		fmt.Println("⌛ 等待合约部署确认...")
		deployedAddr, err = bind.WaitDeployed(context.Background(), client, tx)
		if err != nil {
			log.Fatalf("❌ 合约部署确认失败: %v", err)
		}
		nftContract = deployContract
		fmt.Printf("✅ 合约部署成功！地址: %s\n", deployedAddr.Hex())
	} else {
		// 3.2 连接已部署合约
		fmt.Println("\n🔗 连接已部署的合约...")
		deployedAddr = common.HexToAddress(contractAddress)
		nftContract, err = nftGo.NewNftGo(deployedAddr, client)
		if err != nil {
			log.Fatalf("❌ 连接合约失败: %v", err)
		}
		fmt.Printf("✅ 成功连接合约: %s\n", deployedAddr.Hex())
	}

	// 4. 调用合约只读方法（无 Gas 消耗）
	fmt.Println("\n========== 查询合约信息 ==========")
	// 4.1 查询 NFT 名称
	name, err := nftContract.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 查询 NFT 名称失败: %v", err)
	}
	fmt.Printf("NFT 名称: %s\n", name)

	// 4.2 查询 NFT 符号
	symbol, err := nftContract.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 查询 NFT 符号失败: %v", err)
	}
	fmt.Printf("NFT 符号: %s\n", symbol)

	// 4.3 查询合约拥有者
	owner, err := nftContract.Owner(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 查询合约拥有者失败: %v", err)
	}
	fmt.Printf("合约拥有者: %s\n", owner.Hex())

	// 4.4 查询基础 Token URI
	baseURI, err := nftContract.BaseTokenURI(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 查询基础 URI 失败: %v", err)
	}
	fmt.Printf("基础元数据 URI: %s\n", baseURI)

	// 5. 调用写方法 - 铸造 NFT
	fmt.Println("\n========== 铸造 NFT ==========")
	// 铸造接收地址（使用私钥对应的地址）
	mintTo := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("准备铸造 NFT 到地址: %s\n", mintTo.Hex())

	// 发送铸造交易
	mintTx, err := nftContract.Mint(auth, mintTo)
	if err != nil {
		log.Fatalf("❌ 铸造 NFT 失败: %v", err)
	}
	fmt.Printf("📤 铸造交易哈希: %s\n", mintTx.Hash().Hex())
	fmt.Println("⌛ 等待铸造交易确认...")

	// 等待交易确认
	_, err = bind.WaitMined(context.Background(), client, mintTx)
	if err != nil {
		log.Fatalf("❌ 铸造交易确认失败: %v", err)
	}
	fmt.Println("✅ NFT 铸造成功！")

	// 6. 查询铸造后的 NFT 信息
	fmt.Println("\n========== 查询铸造结果 ==========")
	// 6.1 查询地址的 NFT 持有数量
	balance, err := nftContract.BalanceOf(&bind.CallOpts{}, mintTo)
	if err != nil {
		log.Fatalf("❌ 查询 NFT 余额失败: %v", err)
	}
	fmt.Printf("地址 %s 持有 NFT 数量: %d\n", mintTo.Hex(), balance)

	// 6.2 查询 tokenId=1 的拥有者（首次铸造通常是 1）
	tokenId := big.NewInt(0)
	ownerOf, err := nftContract.OwnerOf(&bind.CallOpts{}, tokenId)
	if err != nil {
		log.Fatalf("❌ 查询 NFT #%d 拥有者失败: %v", tokenId, err)
	}
	fmt.Printf("NFT #%d 拥有者: %s\n", tokenId, ownerOf.Hex())

	// 6.3 查询 NFT #1 的元数据 URI
	tokenURI, err := nftContract.TokenURI(&bind.CallOpts{}, tokenId)
	if err != nil {
		log.Fatalf("❌ 查询 NFT #%d URI 失败: %v", tokenId, err)
	}
	fmt.Printf("NFT #%d 元数据 URI: %s\n", tokenId, tokenURI)

	// 7. 调用写方法 - 转账 NFT（可选）
	fmt.Println("\n========== 转账 NFT ==========")
	// 目标接收地址（替换为你要转账的地址）
	recipient := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	fmt.Printf("准备将 NFT #%d 转账到: %s\n", tokenId, recipient.Hex())

	// 发送转账交易（safeTransferFrom）
	transferTx, err := nftContract.SafeTransferFrom(
		auth,
		mintTo,    // 转出地址
		recipient, // 接收地址
		tokenId,   // NFT ID
	)
	if err != nil {
		log.Fatalf("❌ 转账 NFT 失败: %v", err)
	}
	fmt.Printf("📤 转账交易哈希: %s\n", transferTx.Hash().Hex())
	fmt.Println("⌛ 等待转账交易确认...")

	// 等待转账确认
	_, err = bind.WaitMined(context.Background(), client, transferTx)
	if err != nil {
		log.Fatalf("❌ 转账交易确认失败: %v", err)
	}
	fmt.Println("✅ NFT 转账成功！")

	// 8. 验证转账结果
	newOwner, err := nftContract.OwnerOf(&bind.CallOpts{}, tokenId)
	if err != nil {
		log.Fatalf("❌ 验证转账结果失败: %v", err)
	}
	fmt.Printf("✅ NFT #%d 新拥有者: %s\n", tokenId, newOwner.Hex())

	// 9. 可选：修改基础元数据 URI（仅合约拥有者可操作）
	fmt.Println("\n========== 修改基础 URI ==========")
	newBaseURI := "https://new-example.com/nft/metadata/"
	fmt.Printf("准备修改基础 URI 为: %s\n", newBaseURI)

	setURITx, err := nftContract.SetBaseTokenURI(auth, newBaseURI)
	if err != nil {
		log.Fatalf("❌ 修改基础 URI 失败: %v", err)
	}
	fmt.Printf("📤 修改 URI 交易哈希: %s\n", setURITx.Hash().Hex())
	_, err = bind.WaitMined(context.Background(), client, setURITx)
	if err != nil {
		log.Fatalf("❌ 修改 URI 交易确认失败: %v", err)
	}

	// 验证修改结果
	updatedBaseURI, err := nftContract.BaseTokenURI(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 查询修改后的 URI 失败: %v", err)
	}
	fmt.Printf("✅ 基础 URI 已更新为: %s\n", updatedBaseURI)
}
