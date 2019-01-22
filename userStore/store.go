package userStore

import (
	"log"
)

type Store struct {
	users map[string]map[string]bool
	voteCount map[string]int
	notifySkip map[string]chan bool 
}

func NewStore() *Store {	
	return &Store { 
		users: make(map[string]map[string]bool), 
		voteCount: make(map[string]int),
		notifySkip: make(map[string]chan bool),
	}
}

func (s *Store) AddChannel(roomCode string) chan bool {
	s.notifySkip[roomCode] = make(chan bool)

	// Return the channel to be used by the song poller
	return s.notifySkip[roomCode]
}

func (s *Store) UserExists(id string, roomCode string) bool {
	log.Println(id)
	if _, ok := s.users[roomCode][id]; ok {
		return true
	}

	return false
}

func (s *Store) AddUser(id string, roomCode string) {
	if _, ok := s.users[roomCode]; !ok {
		// Need to initialize a new map if it has not been created
		s.users[roomCode] = make(map[string]bool)
	}
	
	s.users[roomCode][id] = false
}

func (s *Store) RemoveUser(id string, roomCode string) {
	delete(s.users[roomCode], id)
}

func (s *Store) CastUserVote(id string, roomCode string) {
	prevVal := s.users[roomCode][id]
	s.users[roomCode][id] = true

	// Only update count if they haven't voted yet
	if prevVal == false {
		s.voteCount[roomCode]++
	}

	if(s.GetVoteCount(roomCode) > (s.GetTotalUserCount(roomCode)/2)) {
		s.notifySkip[roomCode] <- true
		s.resetUsersVote(roomCode)
	}
}

func (s *Store) UserHasVoted(id string, roomCode string) bool {
	return s.users[roomCode][id]
}

func (s *Store) resetUsersVote(roomCode string) {
	for id := range s.users[roomCode] {
		s.users[roomCode][id] = false	
	}

	s.voteCount[roomCode] = 0
}

func (s *Store) GetTotalUserCount(roomCode string) int {
	return len(s.users[roomCode])
}

func (s *Store) GetVoteCount(roomCode string) int {
	return s.voteCount[roomCode]
}