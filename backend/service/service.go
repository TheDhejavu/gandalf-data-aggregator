package service

import (
	"context"
	"fmt"
	"gandalf-data-aggregator/config"
	"gandalf-data-aggregator/models"
	token "gandalf-data-aggregator/pkg/jwt"
	"gandalf-data-aggregator/repository"
	"gandalf-data-aggregator/store"
	"gandalf-data-aggregator/webapi"
	workertask "gandalf-data-aggregator/worker/tasks"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/markbates/goth/providers/twitter"
)

type Service struct {
	twitterProvider *twitter.Provider
	repo            *repository.Postgres
	gandalfClient   *webapi.GandalfClient
	sessionStore    *store.SessionStore
	jwtMaker        token.Maker
	wt              workertask.WorkerTask
	cfg             config.Config
}

func NewService(cfg config.Config, repo *repository.Postgres, gandalClient *webapi.GandalfClient, sessionStore *store.SessionStore, workertask workertask.WorkerTask, jwtMaker token.Maker) *Service {
	return &Service{
		twitterProvider: twitter.New(cfg.Twitter.Key, cfg.Twitter.Secret, cfg.Twitter.Callback),
		repo:            repo,
		gandalfClient:   gandalClient,
		cfg:             cfg,
		sessionStore:    sessionStore,
		jwtMaker:        jwtMaker,
		wt:              workertask,
	}
}

func (s Service) BeginAuth(ctx context.Context) (string, error) {
	sess, err := s.twitterProvider.BeginAuth("")
	if err != nil {
		return "", fmt.Errorf("unable to begin auth for %w", err)
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", fmt.Errorf("unable to get auth URL for %w", err)
	}

	parsedURL, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}
	oauthToken := parsedURL.Query().Get("oauth_token")

	if err = s.sessionStore.SaveSession(oauthToken, sess.Marshal()); err != nil {
		return "", err
	}
	return authURL, nil
}

func (s Service) CompleteAuth(ctx context.Context, rawAuthURL string) (*models.User, string, error) {
	parsedURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing the url: %w", err)
	}

	session, err := s.sessionStore.GetSession(parsedURL.Query().Get("oauth_token"))
	if err != nil {
		return nil, "", err
	}

	sess, err := s.twitterProvider.UnmarshalSession(session)
	if err != nil {
		return nil, "", err
	}

	_, err = sess.Authorize(s.twitterProvider, parsedURL.Query())
	if err != nil {
		return nil, "", err
	}

	authUser, err := s.twitterProvider.FetchUser(sess)
	if err != nil {
		return nil, "", err
	}

	user, err := s.repo.FindOrCreateUser(ctx, &models.User{
		Username:   authUser.NickName,
		Email:      authUser.Email,
		AvatarURL:  authUser.AvatarURL,
		FirstName:  authUser.FirstName,
		LastName:   authUser.LastName,
		ExternalID: authUser.UserID,
	})
	if err != nil {
		return nil, "", err
	}

	token, err := s.jwtMaker.CreateToken(user.Username, user.ID, 720*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s Service) GenerateGandalfCallback(ctx context.Context, userID uuid.UUID) (string, error) {
	state := uuid.NewString()
	state, err := s.sessionStore.FindKeyOrSaveSession(state, userID.String())
	if err != nil {
		return "", err
	}
	callbackURL := fmt.Sprintf("%s/gandalf/callback/netflix/%s", s.cfg.ServerURL, state)
	return callbackURL, nil
}

func (s Service) RegisterUserDataKey(ctx context.Context, state string, key string, source string) error {
	sessionUserID, err := s.sessionStore.GetSession(state)
	if err != nil {
		return fmt.Errorf("no active session for state: %s", state)
	}

	userID, err := uuid.Parse(sessionUserID)
	if err != nil {
		return err
	}

	dataKey := &models.DataKey{
		UserID:   userID,
		DataType: models.DataType(source),
		Key:      key,
	}

	_, err = s.repo.FindOrCreateDataKey(ctx, dataKey)
	if err != nil {
		return err
	}

	err = s.wt.EnqueueActivityDataResolver(workertask.QueuePayload{
		UserID:  userID,
		DataKey: key,
	})

	if err != nil {
		log.Error().Err(err).Msg("EnqueueActivityDataResolver failed")
	}

	return nil
}

func (s *Service) GetActivitySetByUser(ctx context.Context, userID uuid.UUID, limit int, page int) (*models.ActivityDataSet, error) {
	return s.repo.GetActivitySetByUser(ctx, userID, limit, page)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) FetchAndDumpUserActivities(ctx context.Context, userID uuid.UUID, dataKey string) error {
	limit := 300
	page := 1
	for {
		activityResponse, err := s.gandalfClient.QueryActivities(context.Background(), dataKey, limit, page)
		if err != nil {
			log.Error().Err(err).Msg("QueryActivities on gandalf failed.")
			return err
		}

		if len(activityResponse.Data) == 0 {
			break
		}

		var activities []*models.Activity
		for _, activity := range activityResponse.Data {
			var identifiers []models.Identifier

			for _, identifier := range activity.Metadata.Subject {
				identifiers = append(identifiers, models.Identifier{
					Value:          identifier.Value,
					IdentifierType: identifier.IdentifierType,
				})
			}

			date, err := time.Parse("02/01/2006", activity.Metadata.Date)
			if err != nil {
				log.Error().Err(err).Msg("unable to parse date")
			}
			activities = append(activities, &models.Activity{
				UserID:             userID,
				ProviderActivityID: activity.ID,
				Title:              activity.Metadata.Title,
				Date:               date,
				Subject:            identifiers,
			})
		}

		if _, err := s.repo.CreateActivities(ctx, activities); err != nil {
			log.Error().Err(err).Msg("Unable to create activities")
			return err
		}
		page++
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (s *Service) GenerateActivityStats(ctx context.Context, userID uuid.UUID) error {
	limit := 100
	page := 1

	var activitySet *models.ActivityDataSet
	var err error

	for {
		activitySet, err = s.repo.FetchUnprocessedUserActivities(ctx, userID, limit, page)
		if err != nil {
			return nil
		}

		if len(activitySet.Data) == 0 {
			break
		}

		var activityIDSet uuid.UUIDs
		currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
		yearlyData := make(map[int][]int)

		for _, record := range activitySet.Data {
			year := record.Date.Year()
			month := int(record.Date.Month())
			activityIDSet = append(activityIDSet, record.ID)

			if _, ok := yearlyData[year]; !ok {
				if year == currentYear {
					yearlyData[year] = make([]int, currentMonth)
				} else {
					yearlyData[year] = make([]int, 12)
				}
			} else if year == currentYear && len(yearlyData[year]) < currentMonth {
				additionalMonths := make([]int, currentMonth-len(yearlyData[year]))
				yearlyData[year] = append(yearlyData[year], additionalMonths...)
			}

			yearlyData[year][month-1]++
		}

		var stats []*models.ActivityStat
		for year, months := range yearlyData {
			fmt.Printf("<< Year: %d, Data: %v\n >> ", year, months)
			for month, count := range months {
				stats = append(stats, &models.ActivityStat{
					Year:   year,
					Month:  month + 1,
					Total:  count,
					UserID: userID,
				})
			}
		}

		if err = s.repo.Transaction(ctx, func(ctx context.Context, tx *repository.Postgres) error {
			err = tx.BatchUpsertActivityStat(ctx, stats)
			if err != nil {
				return nil
			}

			err = tx.SetActivityStatsToProcessedByUser(ctx, activityIDSet)
			if err != nil {
				return nil
			}

			return nil
		}); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("<< STATS DONE >> ")
	return nil
}

func (s *Service) GenerateUserYearlyData(ctx context.Context, userID uuid.UUID) (*models.YearDataStat, error) {
	activityStats, err := s.repo.GetActivityStatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	yearlyData := make(map[int]models.YearData)
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	for _, stat := range activityStats {
		if _, ok := yearlyData[stat.Year]; !ok {
			yearlyData[stat.Year] = models.YearData{
				Labels: make([]int, 0),
				Months: make([]string, 0),
			}
		}

		yearData := yearlyData[stat.Year]
		if len(yearData.Labels) < stat.Month {
			extendedLabels := make([]int, stat.Month)
			copy(extendedLabels, yearData.Labels)
			yearData.Labels = extendedLabels

			extendedMonths := make([]string, stat.Month)
			copy(extendedMonths, yearData.Months)
			for i := len(yearData.Months); i < stat.Month; i++ {
				extendedMonths[i] = months[i]
			}
			yearData.Months = extendedMonths
		}

		yearData.Labels[stat.Month-1] = stat.Total
		yearlyData[stat.Year] = yearData
	}

	currentYear := time.Now().Year()
	return &models.YearDataStat{
		YearData:    yearlyData,
		CurrentYear: strconv.Itoa(currentYear),
	}, nil
}
