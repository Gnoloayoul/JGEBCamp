package homework

import (
	_ "embed"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

//go:embed testInit.sql
var testInitSQL string

func TestGenSQL(t *testing.T) {
	file, err := os.OpenFile("testData.sql",
		os.O_RDWR|os.O_APPEND|os.O_TRUNC|os.O_CREATE, 0666)
	require.NoError(t, err)
	defer file.Close()

	const rows = 20
	const prefix = "\n\nINSERT INTO testTables (id, A, B, C)\nVALUES "

	_, err = file.WriteString(testInitSQL)
	require.NoError(t, err)
	_, err = file.WriteString(prefix)
	require.NoError(t, err)

	for i := 0; i < rows; i++ {
		if i > 0 {
			file.Write([]byte{',', '\n'})
			file.WriteString("       ")
		}
		file.Write([]byte{'('})
		file.WriteString(strconv.Itoa(i + 1))
		file.Write([]byte{','})
		file.WriteString(strconv.Itoa(i + 2))
		file.Write([]byte{','})
		file.WriteString(strconv.Itoa(int(rand.Int31n(1000))))
		file.Write([]byte{','})
		file.WriteString(strconv.Itoa(int(rand.Int31n(1000))))
		file.Write([]byte{')'})
	}

	file.Write([]byte{';'})
}
