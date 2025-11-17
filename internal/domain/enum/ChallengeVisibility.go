package enum

type ChallengeVisibility uint

const (
	VisibilityPublic ChallengeVisibility = iota + 1
	VisibilityPrivate
	VisibilityInvite
)

func (v ChallengeVisibility) String() string {
	switch v {
	case VisibilityPublic:
		return "public"
	case VisibilityPrivate:
		return "private"
	case VisibilityInvite:
		return "invite"
	}
	return "unknown"
}

func GetAllChallengeVisibilities() []ChallengeVisibility {
	return []ChallengeVisibility{
		VisibilityPublic,
		VisibilityPrivate,
		VisibilityInvite,
	}
}
