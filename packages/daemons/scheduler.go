package daemons

import (
	"context"
	"fmt"
	"time"

	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/scheduler"
	"github.com/GenesisKernel/go-genesis/packages/scheduler/contract"

	log "github.com/sirupsen/logrus"
)

func loadContractTasks() error {
	stateIDs, err := model.GetAllSystemStatesIDs()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("get all system states ids")
		return err
	}

	for _, stateID := range stateIDs {
		if !model.IsTable(fmt.Sprintf("%d_vde_cron", stateID)) {
			return nil
		}

		c := model.Cron{}
		c.SetTablePrefix(fmt.Sprintf("%d_vde", stateID))
		tasks, err := c.GetAllCronTasks()
		if err != nil {
			log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("get all cron tasks")
			return err
		}

		for _, cronTask := range tasks {
			err = scheduler.UpdateTask(&scheduler.Task{
				ID:       cronTask.UID(),
				CronSpec: cronTask.Cron,
				Handler: &contract.ContractHandler{
					Contract: cronTask.Contract,
				},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// VDEScheduler starts contracts on schedule
func VDEScheduler(ctx context.Context, d *daemon) error {
	d.sleepTime = time.Hour
	return loadContractTasks()
}
