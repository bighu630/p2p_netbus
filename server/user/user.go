// user是一个用户的全部客户端的ip, Server负责同步所有client的net穿透后的ip
// 客户端需要有一个心跳包，用来同步user下的machines
// 重点式客户端决定什么东西要被分享，什么东西不用
// server类似于 p2p下载中的tracher服务器
package user

type User struct {
	userID   string          // 用户ID,根据用户id匹配用户设备列表
	pwd      string          // 用户密码，加个密码可以保证userID不会重复
	machines map[string]bool // ip列表
}

func (u *User) GetUserID() string {
	return u.userID
}

func (u *User) GetUserMachines() []string {
	var addresses []string
	for k := range u.machines {
		addresses = append(addresses, k)
	}
	return addresses
}

func (u *User) AddMachine(machine string) {
	u.machines[machine] = true
}
