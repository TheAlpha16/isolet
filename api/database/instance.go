package database


// func CanStartInstance(c *fiber.Ctx, userid int, level int) bool {
// 	var runid int
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT runid FROM running WHERE userid = $1 AND level = $2`, userid, level).Scan(&runid); err == nil {
// 		return false
// 	}

// 	if _, err := DB.QueryContext(ctx, `INSERT INTO running (userid, level) VALUES ($1, $2)`, userid, level); err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }

// func DeleteRunning(c *fiber.Ctx, userid int, level int) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func NewFlag(c *fiber.Ctx, userid int, level int, password string, flag string, port int32, hostname string, deadline int64) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `INSERT INTO flags (userid, level, flag, password, port, hostname, deadline) VALUES ($1, $2, $3, $4, $5, $6, $7)`, userid, level, flag, password, port, hostname, deadline); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DeleteFlag(c *fiber.Ctx, userid int, level int) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ValidChallenge(c *fiber.Ctx, level int) bool {
// 	var chall_name string
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT chall_name FROM challenges WHERE level = $1`, level).Scan(&chall_name); err != nil {
// 		return false
// 	}
// 	return true
// }

// func GetInstances(c *fiber.Ctx, userid int) ([]models.Instance, error) {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	instances := make([]models.Instance, 0)
// 	rows, err := DB.QueryContext(ctx, `SELECT userid, level, password, port, verified, hostname, deadline from flags WHERE userid = $1`, userid)
// 	if err != nil {
// 		return instances, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		instance := new(models.Instance)
// 		if err := rows.Scan(&instance.UserID, &instance.Level, &instance.Password, &instance.Port, &instance.Verified, &instance.Hostname, &instance.Deadline); err != nil {
// 			return instances, err
// 		}
// 		instances = append(instances, *instance)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return instances, err
// 	}
// 	return instances, nil
// }


// func AddTime(c *fiber.Ctx, userid int, level int) (bool, string, int64) {
// 	var current int
// 	var deadline int64
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT extended, deadline FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&current, &deadline); err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}

// 	if (current + 1) > (config.MAX_INSTANCE_TIME / config.INSTANCE_TIME) {
// 		return false, "limit reached", 1
// 	}

// 	_, err := DB.QueryContext(ctx, `UPDATE flags SET extended = $1 WHERE userid = $2 AND level = $3`, current+1, userid, level)
// 	if err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}

// 	newdeadline := time.UnixMilli(deadline).Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()

// 	_, err = DB.QueryContext(ctx, `UPDATE flags SET deadline = $1 WHERE userid = $2 AND level = $3`, newdeadline, userid, level)
// 	if err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}
// 	return true, "", newdeadline
// }
