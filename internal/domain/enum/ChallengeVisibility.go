package enum

import (
	"encoding/json"
	"fmt"
)

type ChallengeVisibility string

const (
	VisibilityPublic ChallengeVisibility = "public"
	VisibilityPrivate ChallengeVisibility = "private"
	VisibilityInvite ChallengeVisibility = "invite"
)

func (v ChallengeVisibility) IsValid() bool {
	switch v {
	case VisibilityPublic, VisibilityPrivate, VisibilityInvite:
		return true
	default:
		return false
	}
}

func GetAllChallengeVisibilities() []ChallengeVisibility {
	return []ChallengeVisibility{
		VisibilityPublic,
		VisibilityPrivate,
		VisibilityInvite,
	}
}

func (v *ChallengeVisibility) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cv := ChallengeVisibility(s)
	if !cv.IsValid() {
		return fmt.Errorf("invalid challenge visibility: %s", s)
	}
	*v = cv
	return nil
}

func (v ChallengeVisibility) String() string {
	return string(v)
}

func (v ChallengeVisibility) MarshalJSON() ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("invalid challenge visibility: %s", v)
	}
	return json.Marshal(v.String())
}