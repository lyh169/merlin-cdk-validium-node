package jsonrpc

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func Test_StrogeCostTime(t *testing.T) {
	for i := 0; i < 5; i++ {
		testOnce(t, i)
	}
}

func testOnce(t *testing.T, num int) {
	//connurl := "postgresql://prover_user:prover_pass@127.0.0.1:5432/prover_db"
	connurl := "postgresql://prover_user:prover_password@54.179.47.175:3711/prover_db"
	conn, err := pgx.Connect(context.Background(), connurl)
	require.NoError(t, err)
	defer conn.Close(context.Background())

	//sql := "INSERT INTO " + "state.nodes" + " ( hash, data ) VALUES ($1, $2)"

	const nodesCount = 25000
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO state.nodes (hash, data) VALUES ")
	for i := 0; i < nodesCount; i++ {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString(fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
	}
	sql := sqlBuilder.String()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	prik := common.Bytes2Hex(crypto.FromECDSA(privateKey)[:5])
	priv := common.Bytes2Hex(crypto.FromECDSA(privateKey)[:5])
	hash := generate(prik, 32)
	length := 150
	data := generate(priv, length)
	var hashes []string
	var datas []string
	for i := 0; i < nodesCount; i++ {
		hashes = append(hashes, fmt.Sprintf("%s%d", hash, i))
		datas = append(datas, fmt.Sprintf("%s%d", data, i))
	}
	var params []interface{}
	for i := 0; i < nodesCount; i++ {
		params = append(params, hashes[i], datas[i])
	}

	start := time.Now()
	tx, err := conn.Begin(context.Background())
	require.NoError(t, err)

	_, err = tx.Exec(context.Background(), sql, params...)
	require.NoError(t, err)

	err = tx.Commit(context.Background())
	require.NoError(t, err)
	fmt.Println("the num of", num, "cost the time", time.Now().Sub(start), "nodesCount", nodesCount, "data length", length)
}

func generate(prefix string, length int) string {
	n := length/len(prefix) - 1
	result := prefix + "."
	for i := 0; i < n; i++ {
		result += prefix
	}
	return result
}

/*
SELECT length(hash) AS hash_length, length(data) AS data_length
FROM your_table_name
ORDER BY hash ASC -- 按照 hash 列的字典顺序排列
LIMIT 1; -- 只返回第一行

SELECT * FROM state.nodes
ORDER BY hash ASC LIMIT 100

SELECT length(hash) AS hash_length
FROM your_table_name
WHERE some_condition;

SELECT length(hash) AS hash_length
FROM state.nodes
ORDER BY hash ASC
LIMIT 1;

SELECT length(data) AS data_length
FROM state.nodes
ORDER BY data ASC
LIMIT 1;


DESC

*/
