package serv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"signer/service"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Verification struct {
	Username  string `json:"username"`
	Signature string `json:"signature"`
}

type QAs struct {
	QustionsAnswers string `json:"qas"`
}

type handler struct {
	user service.UserService
	test service.TestService
}

func NewHandler(us service.UserService, ts service.TestService) handler {
	return handler{user: us, test: ts}
}

func (h *handler) GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	user := User{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("The user request value %v", user)

	err := h.user.CheckUser(ctx, user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tokenString, err := createToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, tokenString)
}

func (h *handler) Sign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	qas := QAs{}
	json.NewDecoder(r.Body).Decode(&qas)

	user, err := getUsername(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sig, err := h.test.Sign(ctx, user, qas.QustionsAnswers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}
	w.Write([]byte(sig))

}

func (h *handler) CheckSignature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	verification := Verification{}
	json.NewDecoder(r.Body).Decode(&verification)

	timestamp, QAs, err := h.test.CheckSignature(ctx, verification.Username, verification.Signature)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Timestamp: %s \nQAs: \n%s", timestamp, QAs)))
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	user := User{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("The user request value %v", user)

	if err := h.user.RegisterUser(ctx, user.Username, user.Password); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := createToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, tokenString)
}
