package controllers

type Members struct {
	Base
}

// @Title Login interface
// @Description Login interface
// @Param	token header string	true "token值"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /login [get]
func (this *Members) Login() {

}
