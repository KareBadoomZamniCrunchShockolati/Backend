package enum

type UserFeeling string

const (
	FeelingGood UserFeeling = "good"
	FeelingBad  UserFeeling = "bad"
)

func GetAllUserFeelings() []UserFeeling {
	return []UserFeeling{FeelingGood, FeelingBad}
}