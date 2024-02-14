package utils

import (
	capsolver_go "github.com/capsolver/capsolver-go"
)

func NewCapSolver(apiKey, webSite, siteKey string) CapSolverApp {
	return CapSolverApp{
		apikey:     apiKey,
		websiteURL: webSite,
		siteKey:    siteKey,
	}
}

func (cm CapSolverApp) CapSolverV3() (string, error) {
	capSolver := capsolver_go.CapSolver{ApiKey: cm.apikey}
	resp, err := capSolver.Solve(map[string]any{
		"type":       "ReCaptchaV3TaskProxyLess",
		"websiteURL": cm.websiteURL,
		"websiteKey": cm.siteKey,
	})

	if err != nil {
		return "", err
	}

	return resp.Solution.GRecaptchaResponse, nil
}
