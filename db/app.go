package db

import "errors"

//保存用户信息
func SaveUserInfo(openid, nickname string, gender int, city, province, country, avatarurl, unionid string) (uid int, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()
	_, err = tx.Exec(`call save_user(?,?,?,?,?,?,?,?,@a)`, openid, nickname, gender, city, province, country, avatarurl, unionid)
	if err != nil {
		return
	}
	err = tx.QueryRow(`select @a`).Scan(&uid)
	return
}

//创建活动
func CreateActivity(title, descriptiom, act_start_time, act_end_time, join_start_time, join_end_time string, originator int) (err error) {
	_, err = db.Exec(`call createactivity(?,?,?,?,?,?,?)`, title, descriptiom, act_start_time, act_end_time, join_start_time, join_end_time, originator)
	return
}

//发起人取消活动
func CancelActivity(activityId, userId int) (err error) {
	_, err = db.Exec(`call cancelactivity(?,?)`, activityId, userId)
	return
}

//查看自己发起的活动
type Activity struct {
	ID        int
	Title     string
	State     int
	Avatarurl string
	SatrtTime string
}

func GetCreateActivity(userId int) (list []Activity, err error) {
	rows, err := db.Query(`call get_create_activity(?)`, userId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a Activity
		err = rows.Scan(&a.ID, &a.Title, &a.State, &a.Avatarurl, &a.SatrtTime)
		if err != nil {
			return
		}
		list = append(list, a)
	}
	return
}

//报名参加活动
func JoinActivity(activityId, userId int, formid string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()
	_, err = tx.Exec(`call join_activity(?,?,?,@a)`, userId, activityId, formid)
	if err != nil {
		return
	}
	var errmsg string
	err = tx.QueryRow(`select @a`).Scan(&errmsg)
	if err != nil {
		return
	}
	if errmsg != "ok" {
		err = errors.New(errmsg)
	}
	return
}

//取消报名
func CancelJoinActivity(activityId, userId int) (err error) {
	_, err = db.Exec(`delete from TBL_ACTIVITY_USER  where FUSERID=? and FACTIVITYID = ?`, userId, activityId)
	return
}

//查看自己报名的活动
func GetJoinActivity(userId int) (list []Activity, err error) {
	rows, err := db.Query(`call get_join_activity(?)`, userId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a Activity
		err = rows.Scan(&a.ID, &a.Title, &a.State, &a.Avatarurl, &a.SatrtTime)
		if err != nil {
			return
		}
		list = append(list, a)
	}
	return
}

//查看活动具体内容
type User struct {
	Uid       int
	NickName  string
	AvatarUrl string
}
type ActivityDetails struct {
	ID                int
	Title             string
	Description       string
	ActivityStartTime string
	ActivityEndTime   string
	JoinStartTime     string
	JoinEndTime       string
	Originator        string
	HeadUrl           string
	State             int
	PersonnelNum      int
	Personnel         []User
	IsJoin            int
	Time              string
}

func GetActivityByID(id int, uid int) (a ActivityDetails, err error) {
	err = db.QueryRow(`call get_activity_by_id(?)`, id).Scan(
		&a.ID,
		&a.Title,
		&a.Description,
		&a.ActivityStartTime,
		&a.ActivityEndTime,
		&a.JoinStartTime,
		&a.JoinEndTime,
		&a.Originator, &a.State, &a.Time, &a.HeadUrl)
	if err != nil {
		return
	}
	rows, err := db.Query(`call get_joinactivity_member_by_id(?)`, id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Uid, &u.NickName, &u.AvatarUrl)
		if err != nil {
			return
		}
		if u.Uid == uid {
			a.IsJoin = 1
		}
		a.Personnel = append(a.Personnel, u)
	}
	a.PersonnelNum = len(a.Personnel)
	return
}

//获取要发送的模板消息
func GetMessage() (m []string, err error) {
	rows, err := db.Query(`call get_message()`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		err = rows.Scan(&a)
		if err != nil {
			return
		}
		m = append(m, a)
	}
	return
}
