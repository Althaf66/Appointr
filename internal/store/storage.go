package store

import (
	"context"
	"database/sql"
	"time"
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		// GauthCreate(context.Context, string, string) (*User, error)
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
		GetByEmail(context.Context, string) (*User, error)
	}
	Expertise interface {
		Create(context.Context, *Expertise) error
		Get(context.Context) ([]*Expertise, error)
		Update(context.Context, *Expertise) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*Expertise, error)
	}
	Discipline interface {
		Create(context.Context, *Discipline) error
		Get(context.Context) ([]*Discipline, error)
		Update(context.Context, *Discipline) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*Discipline, error)
		GetByField(context.Context, string) ([]*Discipline, error)
	}
	Messages interface {
		CreateConversation(ctx context.Context, userID1, userID2 int64) (*Conversation, error)
		GetConversation(ctx context.Context, conversationID int64) (*Conversation, error)
		GetOrCreateConversationByUsers(ctx context.Context, userID1, userID2 int64) (*Conversation, error)
		GetUserConversations(ctx context.Context, userID int64) ([]*Conversation, error)
		CreateMessage(ctx context.Context, message *Message) error
		GetConversationMessages(ctx context.Context, conversationID int64, limit, offset int) ([]*Message, error)
		MarkConversationAsRead(ctx context.Context, conversationID, userID int64) error
		GetUnreadCount(ctx context.Context, userID int64) (int, error)
	}
	Mentor interface {
		CreateMentor(ctx context.Context, mentor *Mentor) error
		GetAllMentors(ctx context.Context, limit, offset int) ([]*Mentor, error)
		GetMentorByName(ctx context.Context, name string) ([]*Mentor, error)
		GetMentorByID(ctx context.Context, id int64) (*Mentor, error)
		UpdateMentor(ctx context.Context, mentor *Mentor) error
		DeleteMentor(ctx context.Context, mentorID int64) error
		GetMentorByUserID(ctx context.Context, userid int64) (*Mentor, error)
		GetMentorsByExpertise(ctx context.Context, expertise string) ([]*Mentor, error)
		GetMentorsByDiscipline(ctx context.Context, discipline string) ([]*Mentor, error)
	}
	Gig interface {
		CreateGig(ctx context.Context, gig *Gig) error
		GetAllGigs(ctx context.Context, limit, offset int) ([]*Gig, error)
		GetGigsByExpertise(ctx context.Context, expertise string) ([]*Gig, error)
		UpdateGig(ctx context.Context, gig *Gig) error
		DeleteGig(ctx context.Context, gigID int64) error
		GetGigByID(ctx context.Context, id int64) (*Gig, error)
		// GetGigByMentorID(ctx context.Context, mentorID int64) ([]*Gig, error)
	}
	Education interface {
		CreateEducation(ctx context.Context, education *Education) error
		GetEducationById(ctx context.Context, id int64) (*Education, error)
		GetEducationByUserId(ctx context.Context, userid int64) ([]*Education, error)
		UpdateEducation(ctx context.Context, education *Education) error
		DeleteEducation(ctx context.Context, educationID int64) error
	}
	Experience interface {
		CreateExperience(ctx context.Context, experience *Experience) error
		GetExperienceById(ctx context.Context, id int64) (*Experience, error)
		GetExperienceByUserId(ctx context.Context, userid int64) ([]*Experience, error)
		UpdateExperience(ctx context.Context, experience *Experience) error
		DeleteExperience(ctx context.Context, experienceID int64) error
	}
	SocialMedia interface {
		CreateSocialMedia(ctx context.Context, socialmedia *SocialMedia) error
		GetSocialMediaById(ctx context.Context, id int64) (*SocialMedia, error)
		GetSocialMediaByUserId(ctx context.Context, userid int64) ([]*SocialMedia, error)
		UpdateSocialMedia(ctx context.Context, socialmedia *SocialMedia) error
		DeleteSocialMedia(ctx context.Context, socialmediaID int64) error
	}
	WorkingAt interface {
		CreateWorkingAt(ctx context.Context, workingat *WorkingAt) error
		GetWorkingAtById(ctx context.Context, id int64) (*WorkingAt, error)
		GetWorkingAtByUserId(ctx context.Context, userid int64) ([]*WorkingAt, error)
		UpdateWorkingAt(ctx context.Context, workingat *WorkingAt) error
		DeleteWorkingAt(ctx context.Context, workingatID int64) error
	}
	Meetings interface {
		CreateMeeting(ctx context.Context, meeting *Meetings) error
		GetAllMeetings(ctx context.Context, limit, offset int) ([]*Meetings, error)
		GetMeetingByUserID(ctx context.Context, userID int64) ([]*Meetings, error)
		GetMeetingMentorNotConfirm(ctx context.Context, mentorID int64) ([]*Meetings, error)
		GetMeetingUserNotPaid(ctx context.Context, userID int64) ([]*Meetings, error)
		GetMeetingUserNotCompleted(ctx context.Context, userID int64) ([]*Meetings, error)
		GetMeetingMentorNotCompleted(ctx context.Context, mentorID int64) ([]*Meetings, error)
		UpdateMeetingConfirm(ctx context.Context, meetingID int64) error
		UpdateMeetingPaid(ctx context.Context, meetingID int64) error
		UpdateMeetingCompleted(ctx context.Context, meetingID int64) error
		UpdateLink(ctx context.Context, meeting *Meetings) error
		DeleteMeeting(ctx context.Context, meetingID int64) error
		GetMeetingByID(ctx context.Context, id int64) (*Meetings, error)
	}
	BookingSlot interface {
		CreateBookingSlot(ctx context.Context, slot *BookingSlot) error
	}
	Country interface {
		GetCountry(ctx context.Context) ([]*Country, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users:       &UserStore{db},
		Expertise:   &ExpertiseStore{db},
		Discipline:  &DisciplineStore{db},
		Messages:    &MessageStore{db},
		Mentor:      &MentorStore{db},
		Gig:         &GigStore{db},
		Education:   &EducationStore{db},
		Experience:  &ExperienceStore{db},
		SocialMedia: &SocialMediaStore{db},
		WorkingAt:   &WorkingAtStore{db},
		BookingSlot: &BookingStore{db},
		Meetings:    &MeetingsStore{db},
		Country:	 &CountryStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
