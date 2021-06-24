package blockchain

type User struct {
	Ballance float64
}

func (user *User) GetBallance() float64 {
	return user.Ballance
}
