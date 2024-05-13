package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"follower.xws.com/model"
	"follower.xws.com/repository"
	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type UserHandler struct {
	logger *log.Logger
	repo   *repository.UserRepository
}

func NewFollowersHandler(l *log.Logger, r *repository.UserRepository) *UserHandler {
	return &UserHandler{l, r}
}

func (f *UserHandler) CreateUser(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyProduct{}).(*model.User)
	userSaved, err := f.repo.SaveUser(user)
	if err != nil {
		f.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userSaved {
		f.logger.Print("New user saved to database")
		rw.WriteHeader(http.StatusCreated)
	} else {
		rw.WriteHeader(http.StatusConflict)
	}
}

func (f *UserHandler) CreateFollowing(rw http.ResponseWriter, h *http.Request) {

	decoder := json.NewDecoder(h.Body)
	defer h.Body.Close()

	var users []model.User
	if err := decoder.Decode(&users); err != nil {

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error decoding JSON data: " + err.Error()))
		return
	}

	if len(users) != 2 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Expecting exactly two users"))
		return
	}

	User1 := model.User{Id: users[0].Id, Username: users[0].Username, ProfileImage: users[0].ProfileImage}
	User2 := model.User{Id: users[1].Id, Username: users[1].Username, ProfileImage: users[1].ProfileImage}
	err := f.repo.SaveFollowing(&User1, &User2)
	if err != nil {
		f.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	User1 = model.User{}
	jsonData, _ := json.Marshal(User1)
	rw.Write(jsonData)

}
func (f *UserHandler) Unfollow(rw http.ResponseWriter, h *http.Request) {
	userId1 := h.URL.Query().Get("followerId")
	userId2 := h.URL.Query().Get("followedId")
	err := f.repo.DeleteFollowing(userId1, userId2)
	if err != nil {
		f.logger.Print("Database exception: ", err)
		return
	}
	user := model.User{}
	jsonData, _ := json.Marshal(user)
	rw.Write(jsonData)
}

func (f *UserHandler) GetFollowingsForUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["userId"]
	users, err := f.repo.GetFollowingsForUser(id)
	if err != nil {
		f.logger.Print("Database exception: ", err)
	}
	if users == nil {
		users = model.Users{}
		jsonData, _ := json.Marshal(users)
		rw.Write(jsonData)
		return
	}
	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		f.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (f *UserHandler) GetFollowersForUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["userId"]
	users, err := f.repo.GetFollowersForUser(id)
	if err != nil {
		f.logger.Print("Database exception: ", err)
	}
	if users == nil {
		users = model.Users{}
		jsonData, _ := json.Marshal(users)
		rw.Write(jsonData)
		return
	}
	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		f.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (f *UserHandler) Recommendations(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["userId"]
	users, err := f.repo.Recommendations(id)
	if err != nil {
		f.logger.Print("Database exception: ", err)
	}
	if users == nil {
		users = model.Users{}
		jsonData, _ := json.Marshal(users)
		rw.Write(jsonData)
		return
	}
	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		f.logger.Fatal("Unable to convert to json :", err)
		return
	}
}
