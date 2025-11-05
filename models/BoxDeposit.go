package models

type BoxDeposit struct {
	FirstName      string `bson:"firstName" json:"firstName"`
	LastName       string `bson:"lastName" json:"lastName"`
	Email          string `bson:"email" json:"email"`
	Phone          string `bson:"phone" json:"phone"`
	ContactMethod  string `bson:"contactMethod" json:"contactMethod"`
	BoxSize        string `bson:"boxSize" json:"boxSize"`
	Duration       string `bson:"duration" json:"duration"`
	Referral       string `bson:"referral" json:"referral"`
	AdditionalInfo string `bson:"additionalInfo" json:"additionalInfo"`
}
