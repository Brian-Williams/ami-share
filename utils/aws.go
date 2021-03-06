package utils

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

type SessionKey struct {
	AccountID  string
	AssumeRole string
	Region     string
}

// Returns an AWS session by profile name and region
// Profile name must be present in the credentials file
// Master session must be initialized to assume roles in other accounts
type AWSSessionFactory struct {
	logger        *log.Entry
	MasterSession *session.Session
	SessionCache  map[SessionKey]*session.Session
}

func NewAWSSessionFactory() *AWSSessionFactory {
	return &AWSSessionFactory{
		SessionCache: make(map[SessionKey]*session.Session),
		logger: log.WithFields(log.Fields{
			"context":   "aws-session-factory",
			"operation": "session",
		}),
	}
}

func (sessionFactory *AWSSessionFactory) GenerateMasterSession(sessionKey SessionKey) (*session.Session, error) {
	sessionFactory.logger.Debugf("Creating master session with: %v", sessionKey)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sessionKey.Region),
	})
	sessionFactory.MasterSession = sess
	return sess, err
}

func (sessionFactory *AWSSessionFactory) GetSession(sessionKey SessionKey) (*session.Session, error) {
	var err error
	sess, ok := sessionFactory.SessionCache[sessionKey]
	if !ok {
		if sessionFactory.MasterSession == nil {
			return nil, errors.New("master session not initialized - required to assume role")
		}

		sessionFactory.logger.Debugf("Generating session: %v", sessionKey)
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(sessionKey.Region),
			Credentials: stscreds.NewCredentials(sessionFactory.MasterSession, sessionKey.AssumeRole),
		})
		if err != nil {
			return sess, err
		}
		sessionFactory.SessionCache[sessionKey] = sess
	}

	return sess, err
}

func (sessionKey SessionKey) String() string {
	return fmt.Sprintf("%s in %s as [%s]", sessionKey.AccountID, sessionKey.Region, sessionKey.AssumeRole)
}
