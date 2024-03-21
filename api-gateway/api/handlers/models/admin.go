package models

type AdminReq struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Age      int64  `json:"age"`
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AdminResp struct {
	Id           string `json:"id"`
	FullName     string `json:"full_name"`
	Age          int64  `json:"age"`
	Email        string `json:"email"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
}

type AdminLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResp struct {
	AccessToken string `json:"access_token"`
}

type SuperAdminMessage struct {
	Message string `json:"message"`
}

type DeleteAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RbacAllRolesResp struct {
	Roles []string `json:"roles"`
}

type Policy struct {
	Role     string `json:"role"`
	EndPoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type ListRolePolicyResp struct {
	Policies []*Policy
}

type AddPolicyRequest struct {
	Policy Policy
}
