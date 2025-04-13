package auth

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitializeCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	modelPath := filepath.Join("internal", "config", "model.conf")
	policyPath := filepath.Join("internal", "config", "policy.csv")

	// Initialize Gorm adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM adapter: %v", err)
	}

	// Create enforcer from model and adapter
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %v", err)
	}

	enforcer.EnableAutoSave(true)

	// Clear all current policies from both memory and database
	enforcer.ClearPolicy()
	if err := enforcer.SavePolicy(); err != nil {
		return nil, fmt.Errorf("failed to save cleared policies: %v", err)
	}

	// Open and load from CSV
	file, err := os.Open(policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open policy file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read policy file: %v", err)
		}

		// Inside the CSV reading loop:
		if len(record) >= 4 && record[0] == "p" {
			sub := strings.TrimSpace(record[1]) // Trim spaces
			obj := strings.TrimSpace(record[2]) // Trim spaces
			act := strings.TrimSpace(record[3]) // Trim spaces

			if _, err := enforcer.AddPolicy(sub, obj, act); err != nil {
				return nil, fmt.Errorf("failed to add policy: %v", err)
			}
			logrus.Infof("✅ Added policy: %s, %s, %s", sub, obj, act)
		}
	}

	if err := enforcer.SavePolicy(); err != nil {
		return nil, fmt.Errorf("failed to persist loaded policies: %v", err)
	}

	policies, err := enforcer.GetPolicy()
	if err != nil {
		logrus.Errorf("❌ Failed to get policies: %v", err)
	} else {
		logrus.Infof("✅ Final loaded policies: %+v", policies)
	}

	return enforcer, nil
}
