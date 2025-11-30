package role

type Role int

const (
    User     Role = iota
    Moderator
    Guest         
)