package context

import (
	"os"
	"sync"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/google/uuid"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type NFContext struct {
	NfId        string
	Name        string
	UriScheme   models.UriScheme
	BindingIPv4 string
	SBIPort     int

	SpyFamilyData map[string]string
	Tasks         []Task
	TaskMutex     sync.RWMutex
	NextTaskID    uint64
}

var nfContext = NFContext{}

func InitNfContext() {
	cfg := factory.NfConfig

	nfContext.NfId = uuid.New().String()
	nfContext.Name = "ANYA"

	nfContext.UriScheme = cfg.Configuration.Sbi.Scheme
	nfContext.SBIPort = cfg.Configuration.Sbi.Port
	nfContext.BindingIPv4 = os.Getenv(cfg.Configuration.Sbi.BindingIPv4)
	if nfContext.BindingIPv4 != "" {
		logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
	} else {
		nfContext.BindingIPv4 = cfg.Configuration.Sbi.BindingIPv4
		if nfContext.BindingIPv4 == "" {
			logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
			nfContext.BindingIPv4 = "0.0.0.0"
		}
	}
	nfContext.SpyFamilyData = map[string]string{
		"Loid":   "Forger",
		"Anya":   "Forger",
		"Yor":    "Forger",
		"Bond":   "Forger",
		"Becky":  "Blackbell",
		"Damian": "Desmond",
		"Franky": "Franklin",
		"Fiona":  "Frost",
		"Sylvia": "Sherwood",
		"Yuri":   "Briar",
		"Millie": "Manis",
		"Ewen":   "Egeburg",
		"Emile":  "Elman",
		"Henry":  "Henderson",
		"Martha": "Marriott",
	}
	nfContext.Tasks = make([]Task, 0)
	nfContext.NextTaskID = 0
}

func GetSelf() *NFContext {
	return &nfContext
}
