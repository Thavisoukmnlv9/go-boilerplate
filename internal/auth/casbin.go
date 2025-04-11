package auth

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// InitializeCasbin initializes the Casbin enforcer with the model and policies.
func InitializeCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	// Define paths to configuration files relative to the project root
	modelPath := filepath.Join("internal", "config", "model.conf")

	// Create GORM adapter for database integration
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM adapter: %v", err)
	}

	// Initialize the Casbin enforcer with the model and adapter
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %v", err)
	}
	policy, err := enforcer.GetPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get existing policies: %v", err)
	}
	logrus.Infof("Existing policies: %v", policy)
	// Check if there are any existing policies in the database
	if len(policy) == 0 {
		policyPath := filepath.Join("internal", "config", "policy.csv")

		// Open the policy CSV file
		file, err := os.Open(policyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open policy file: %v", err)
		}
		defer file.Close()

		// Create a CSV reader
		reader := csv.NewReader(file)
		for {
			// Read each record (line) from the CSV
			record, err := reader.Read()
			if err == io.EOF {
				break // End of file reached
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read policy file: %v", err)
			}

			// Process policy lines starting with "p" (permission policies)
			if len(record) >= 4 && record[0] == "p" {
				sub, obj, act := record[1], record[2], record[3]
				if _, err := enforcer.AddPolicy(sub, obj, act); err != nil {
					return nil, fmt.Errorf("failed to add policy: %v", err)
				}
				logrus.Infof("Added policy: %s, %s, %s", sub, obj, act)
			}
		}

		// Save the loaded policies to the database
		if err := enforcer.SavePolicy(); err != nil {
			return nil, fmt.Errorf("failed to save policy to database: %v", err)
		}
	}

	// Enable auto-save for future policy changes
	enforcer.EnableAutoSave(true)
	return enforcer, nil
}
