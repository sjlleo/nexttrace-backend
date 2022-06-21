package dbtools

type Users struct {
	Uid int
	Ip  string
}

func SearchUsers(ip string) (int, error) {
	db := GetDB()
	u := Users{}
	res := db.Where("ip = ?", ip).Take(&u)
	return u.Uid, res.Error
}

func AddUsers(ip string) (int, error) {
	db := GetDB()
	u := Users{
		Ip: ip,
	}
	res := db.Create(&u)
	return u.Uid, res.Error
}
