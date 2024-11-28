package workers

import (
	"go-todo-api/internal/config"
	"log"

	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
)

type MailWorker struct {
	Log    *logrus.Logger
	Mailer *config.MailerConfig
}

func NewMailWorker(logger *logrus.Logger, mailer *config.MailerConfig) *MailWorker {
	return &MailWorker{
		Log:    logger,
		Mailer: mailer,
	}
}

func (w *MailWorker) SendEmail(job *work.Job) error {
	to := job.ArgString("to")
	subject := job.ArgString("subject")
	body := job.ArgString("body")

	if err := job.ArgError(); err != nil {
		return err
	}

	if err := w.Mailer.SendMail(to, subject, body); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func RegisterMailJobs(workerPool *work.WorkerPool) {
	workerPool.Job("send_email", (*MailWorker).SendEmail)
}
