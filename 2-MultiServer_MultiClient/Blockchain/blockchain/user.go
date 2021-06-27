package blockchain

import "github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"

//
// MARK: User struct
//
type User struct {
	Uuid     string  `json:"uuid"`
	Ballance float64 `json:"ballance"`
	Node     string  `json:"node"`
}

//
// note: Methods
//

func (user *User) UpdateUserBalance(amount float64) {
	user.Ballance += amount
}

//
// note: Formatting From and To userpb
//

// Returns a new blockchainpb_user from a user object
func (user User) ToUserpbfmt() *blockchainpb.User {
	return &blockchainpb.User{
		Uuid: user.Uuid,
		Node: user.Node,
		// Ballance: user.Ballance,
	}
}

// Returns a new user object from a blockcainpb_user
func FromUserpbfmt(user_pb *blockchainpb.User) *User {
	return &User{
		Uuid: user_pb.GetUuid(),
		Node: user_pb.GetNode(),
		// Ballance: user_pb.GetBallance(),
	}
}
