// 权益证明机制最开始是由点点币提出并应用（出块概率=代币数量 * 币龄）
// 简单来说谁的币多，谁就有更大的出块概率。
// 但是深挖下去，这个出块概率谁来计算？
// 碰到无成本利益关系问题怎么办?
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"
)

// 区块结构
type block struct {
	//上一个块的hash
	prehash string
	//本块hash
	hash string
	//时间戳
	timestamp string
	//区块内容
	data string
	//区块高度
	height int
	//挖出本块的地址
	address string
}

//用于存储区块链
var blockchain []block

//代表挖矿节点
type node struct {
	//代币数量
	tokens int
	//质押时间
	days int
	//节点地址
	address string
}

//挖矿节点 用来存放指定的挖矿节点
var mineNodesPool []node

//概率节点池 用于存入挖矿节点的代币数量*币龄获得的概率
var probabilityNodesPool []node

//初始化
func init() {
	//手动添加两个节点
	mineNodesPool = append(mineNodesPool, node{1000, 1, "AAAAAAAAAA"})
	mineNodesPool = append(mineNodesPool, node{100, 3, "BBBBBBBBBB"})
	//初始化随机节点池（挖矿概率与代币数量和币龄有关）
	for _, v := range mineNodesPool {
		for i := 0; i <= v.tokens*v.days; i++ {
			probabilityNodesPool = append(probabilityNodesPool, v)
		}
	}
}

//生成新的区块
func generateNewBlock(oldBlock block, data string, address string) block {
	newBlock := block{}
	newBlock.prehash = oldBlock.hash
	newBlock.data = data
	newBlock.timestamp = time.Now().Format("2006-01-02 15:04:05")
	newBlock.height = oldBlock.height + 1
	newBlock.address = getMineNodeAddress()
	newBlock.getHash()
	return newBlock
}

//对自身进行散列
func (b *block) getHash() {
	sumString := b.prehash + b.timestamp + b.data + b.address + strconv.Itoa(b.height)
	hash := sha256.Sum256([]byte(sumString))
	b.hash = hex.EncodeToString(hash[:])
}

// 每次挖矿都会从概率节点池中随机选出获得出块权的节点地址
// 随机得出挖矿地址（挖矿概率跟代币数量与币龄有关）
func getMineNodeAddress() string {
	bInt := big.NewInt(int64(len(probabilityNodesPool)))
	//得出一个随机数，最大不超过随机节点池的大小
	rInt, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		log.Panic(err)
	}
	return probabilityNodesPool[int(rInt.Int64())].address
}

func main() {
	//创建创世区块
	genesisBlock := block{"0000000000000000000000000000000000000000000000000000000000000000", "", time.Now().Format("2006-01-02 15:04:05"), "我是创世区块", 1, "0000000000"}
	genesisBlock.getHash()
	//把创世区块添加进区块链
	blockchain = append(blockchain, genesisBlock)
	fmt.Println(blockchain[0])
	i := 0
	for {
		time.Sleep(time.Second)
		newBlock := generateNewBlock(blockchain[i], "我是区块内容", "00000")
		blockchain = append(blockchain, newBlock)
		fmt.Println(blockchain[i+1])
		i++
	}
}
