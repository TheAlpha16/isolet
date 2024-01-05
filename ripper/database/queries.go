package database

func DeleteFlag(userid int, level int) error {
	if _, err := DB.Query(`DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func DeleteRunning(userid int, level int) error {
	if _, err := DB.Query(`DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}