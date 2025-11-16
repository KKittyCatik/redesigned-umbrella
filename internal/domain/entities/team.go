package entities

type Team struct {
	Name    string
	Members []*User
}

func NewTeam(name string, members []*User) *Team {
	if name == "" {
		return nil
	}
	return &Team{
		Name:    name,
		Members: members,
	}
}

func (t *Team) AddMember(user *User) error {
	if user == nil {
		return ErrUserNotFound
	}
	if t.HasMember(user.ID) {
		return ErrMemberExists
	}
	t.Members = append(t.Members, user)
	return nil
}

func (t *Team) GetActiveMembers() []*User {
	var activeMembers []*User
	for _, member := range t.Members {
		if member.IsActive {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

func (t *Team) HasMember(userID string) bool {
	for _, member := range t.Members {
		if member.ID == userID {
			return true
		}
	}
	return false
}

func (t *Team) RemoveMember(userID string) error {
	for i, member := range t.Members {
		if member.ID == userID {
			t.Members = append(t.Members[:i], t.Members[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}
