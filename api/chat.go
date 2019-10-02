/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 11:17
 */
package api

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

//api server,Also, you can use gin,echo ... framework wrap
func (c *Chat) Run() {
	//login
	//register
	//CheckAuth
	//Logout

	//Push
	//PushRoom
	//Count
	//GetRoomInfo

}
