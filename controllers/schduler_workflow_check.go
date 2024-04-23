package controllers

import (
	"log"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"time"
)

func ScheduleWorkflowCheck() {
	// Get all projects, preload workflows if workflow.IsUpdatedLastest is false, preload WorkflowLogs
	var projects []models.Project
	for {
		log.Println("Checking for workflows...")
		DB.
			Preload("Workflows", "is_updated_lastest = ?", false).
			// Preload("Workflows.Logs").
			Find(&projects)

		for _, project := range projects {
			// log.Println("project Name : ", project.Name)
			// parse repoURL
			repo, err := github.ParseRepoURL(project.RepoURL)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, workflow := range project.Workflows {
				log.Println("workflow Name : ", workflow.Name)
				log.Println("workflow Logs : ", workflow.Logs)
				// var w models.Workflow
				// DB.Preload("Logs","workflow_id = ?",workflow.WorkflowID).Find(&w)
				var l []models.WorkflowLog
				DB.Where("workflow_id = ? and state != ? ", workflow.WorkflowID, "completed").Find(&l)
				for _, logg := range l {
					log.Println("workflow log : ", logg)
					if logg.State == "completed" {
						continue
					}
					// else

					// print all this repo.Owner, repo.Name, project.RepoApiToken, workflow.WorkflowID, workflow.AttemptNumber
					log.Println(repo.Owner, repo.Name, project.RepoApiToken, logg.WorkflowRunId, workflow.AttemptNumber)
					duration, err := github.GetLastAttemptInformationOfWorkflowRun(repo.Owner, repo.Name, project.RepoApiToken, int(logg.WorkflowRunId), workflow.AttemptNumber)
					if err != nil {
						log.Println(err.Error())
						continue
					} else {
						logg.State = "completed"
						logg.Duration = int(duration.Milliseconds())
						DB.Save(&logg)
						workflow.IsUpdatedLastest = true
						log.Println(duration)
					}
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
	// parse repoURL
	// get all workflows from repoURL
	// Get all users
	// Get all logs
	// Get all agent logs
	// Get all project logs
	// Get all user logs
	// Get all workflow logs
	// Get all workflow runs
	// Get all workflow run attempts
	// Get all workflow run timings
}
