package gomaxcompute

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	for i := 0; i < 100; i++ {
		const stmt = `SELECT * FROM iris_test LIMIT 1;`
		_, err = db.Query(stmt)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(-1)
		}
	}
}

//func TestQuery2(t *testing.T) {
//	a := assert.New(t)
//	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
//	a.NoError(err)
//
//	for i := 0; i < 100; i++ {
//		//const stmt = `select * from ant_p13n_dev.softmax_estimator_train_data limit 10;`
//		const stmt = `SELECT * FROM yiyang_test_table1;`
//		_, err := db.Query(stmt)
//		if err != nil {
//			fmt.Printf("%+v\n", err)
//			os.Exit(-1)
//		}
//	}
//}
