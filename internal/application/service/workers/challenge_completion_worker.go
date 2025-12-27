package workers

import (
	"log"
	"time"

	service "challenge-app/internal/application/service"
	"challenge-app/internal/domain/repository"
)

type ChallengeCompletionWorker struct {
	completionService *service.ChallengeCompletionService
	interval          time.Duration
	stopChan          chan bool
	isRunning         bool
}

func NewChallengeCompletionWorker(
	completionRepo repository.ChallengeCompletionRepository,
	challengeRepo repository.ChallengeRepository,
	participantRepo repository.ChallengeParticipantRepository,
	userDayRepo repository.UserDayRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
	interval time.Duration,
) *ChallengeCompletionWorker {
	completionService := service.NewChallengeCompletionService(
		completionRepo,
		challengeRepo,
		participantRepo,
		userDayRepo,
		userRepo,
		categoryRepo,
	)

	return &ChallengeCompletionWorker{
		completionService: completionService,
		interval:          interval,
		stopChan:          make(chan bool),
		isRunning:         false,
	}
}

func (w *ChallengeCompletionWorker) Start() {
	if w.isRunning {
		return
	}

	w.isRunning = true
	ticker := time.NewTicker(w.interval)

	go func() {
		log.Println("Challenge completion worker started")

		for {
			select {
			case <-ticker.C:
				w.processBatch()
			case <-w.stopChan:
				ticker.Stop()
				log.Println("Challenge completion worker stopped")
				return
			}
		}
	}()
}

func (w *ChallengeCompletionWorker) Stop() {
	if w.isRunning {
		w.stopChan <- true
		w.isRunning = false
	}
}

func (w *ChallengeCompletionWorker) processBatch() {
	log.Println("Processing ended challenges...")

	startTime := time.Now()
	completedCount, err := w.completionService.ProcessEndedChallenges()

	if err != nil {
		log.Printf("Error processing ended challenges: %v", err)
		return
	}

	duration := time.Since(startTime)

	if completedCount > 0 {
		log.Printf("Marked %d challenge(s) as completed in %v", completedCount, duration)
	} else {
		log.Printf("No challenges to process (took %v)", duration)
	}
}

func (w *ChallengeCompletionWorker) ProcessNow() (int, error) {
	return w.completionService.ProcessEndedChallenges()
}
