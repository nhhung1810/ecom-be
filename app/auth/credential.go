package auth

// For parsing the json request
type Credential struct {
	//TODO: Check for UNICODE compatible
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
