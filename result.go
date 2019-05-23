package gomaxcompute

type odpsResult struct {
	affectedRows int64
	insertId     int64
}

func (res *odpsResult) LastInsertId() (int64, error) {
	return res.insertId, nil
}

func (res *odpsResult) RowsAffected() (int64, error) {
	return res.affectedRows, nil
}
