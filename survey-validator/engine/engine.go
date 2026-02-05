package engine

// engine.go - runs all validation checks concurrently

import (
	"sync"
	"time"

	"github.com/survey-validator/domain"
	"github.com/survey-validator/models"
)

type ValidationCheck func(data *models.SurveyData) []models.ValidationIssue

type Engine struct {
	checks map[string]ValidationCheck
}

func NewEngine() *Engine {
	e := &Engine{
		checks: make(map[string]ValidationCheck),
	}

	// add all the checks we want to run
	e.RegisterCheck("input_validation", domain.ValidateInput)
	e.RegisterCheck("duplicate_detection", domain.DetectDuplicates)
	e.RegisterCheck("distance_bearing_check", domain.CheckDistanceAndBearing)
	e.RegisterCheck("outlier_detection", domain.DetectOutliers)
	e.RegisterCheck("traverse_closure", domain.CheckTraverseClosure)

	return e
}

func (e *Engine) RegisterCheck(name string, check ValidationCheck) {
	e.checks[name] = check
}

type checkResult struct {
	checkName string
	issues    []models.ValidationIssue
}

// Validate - runs all checks in parallel, collects results
func (e *Engine) Validate(data *models.SurveyData) *models.ValidationReport {
	return e.ValidateWithOptions(data, nil)
}

// ValidateWithOptions - validation with optional traverse adjustment settings
func (e *Engine) ValidateWithOptions(data *models.SurveyData, traverseInput *models.TraverseInput) *models.ValidationReport {
	startTime := time.Now()

	report := models.NewValidationReport(data.ProjectID)
	resultChan := make(chan checkResult, len(e.checks))
	var wg sync.WaitGroup

	for name, check := range e.checks {
		wg.Add(1)
		go func(checkName string, checkFunc ValidationCheck) {
			defer wg.Done()

			issues := checkFunc(data)
			resultChan <- checkResult{
				checkName: checkName,
				issues:    issues,
			}
		}(name, check)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		report.ChecksPerformed = append(report.ChecksPerformed, result.checkName)
		for _, issue := range result.issues {
			report.AddIssue(issue)
		}
	}

	report.Summary = domain.CalculateSummaryStatistics(data)

	// run traverse adjustment if we have traverse points
	if report.Summary.TraversePoints >= 3 {
		report.TraverseResult = domain.ComputeTraverseAdjustment(data, traverseInput)
		report.ChecksPerformed = append(report.ChecksPerformed, "bowditch_adjustment")
	}

	report.CalculateConfidenceScore()
	report.ProcessingTime = time.Since(startTime).String()

	return report
}
