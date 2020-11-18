package history

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	//conn, err := service.pool.Acquire(context.Background())
	//defer conn.Release()
	// TODO: panic
	_, err := service.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS historyoperationslog (
		id BIGSERIAL,
		name    TEXT NOT NULL,
		number TEXT NOT NULL,
   		recipientSender TEXT NOT NULL,
   		count INTEGER NOT NULL,
   		balanceold INTEGER NOT NULL,
   		balancenew INTEGER NOT NULL,
   		time    INTEGER NOT NULL,
		owner_id INTEGER NOT NULL
	);
`)
	if err != nil {
		log.Print(err)
	}
}

// created timestamp
// modified timestamp

// CRUD
func (service *Service) All() (models []ModelOperationsLog, err error) {
	rows, err := service.pool.Query(context.Background(), `SELECT id, name, number,recipientSender,count, balanceold, balancenew, time, owner_id FROM historyoperationslog;`)
	if err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := ModelOperationsLog{}
		err = rows.Scan(
			&model.Id,
			&model.Name,
			&model.Number,
			&model.RecipientSender,
			&model.Count,
			&model.BalanceOld,
			&model.BalanceNew,
			&model.Time,
			&model.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't get history from db: %w", err)
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	return models, nil
}

func (service *Service) UserShowTransferLogByIdCard(idCard int, idUser int) (model []ModelOperationsLog, err error) {
	modHistoryLog := ModelOperationsLog{}
	err = service.pool.QueryRow(context.Background(), `
SELECT id, name, number, recipientSender, count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog WHERE id=$1 and owner_id=$2`, idCard, idUser).Scan(
		&modHistoryLog.Id,
		&modHistoryLog.Name,
		&modHistoryLog.Number,
		&modHistoryLog.RecipientSender,
		&modHistoryLog.Count,
		&modHistoryLog.BalanceOld,
		&modHistoryLog.BalanceNew,
		&modHistoryLog.Time,
		&modHistoryLog.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	model = append(model, modHistoryLog)
	return model, nil
}

func (service *Service) AdminShowTransferLogByIdCadr(id int) (model []ModelOperationsLog, err error) {
	modHistoryLog := ModelOperationsLog{}
	err = service.pool.QueryRow(context.Background(), `
SELECT id, name, number, recipientSender, count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog WHERE id=$1`, id).Scan(
		&modHistoryLog.Id,
		&modHistoryLog.Name,
		&modHistoryLog.Number,
		&modHistoryLog.RecipientSender,
		&modHistoryLog.Count,
		&modHistoryLog.BalanceOld,
		&modHistoryLog.BalanceNew,
		&modHistoryLog.Time,
		&modHistoryLog.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	model = append(model, modHistoryLog)
	return model, nil
}

func (service *Service) ShowOperationsLogByOwnerId(id int) (models []ModelOperationsLog, err error) {
	rows, err := service.pool.Query(context.Background(), `
SELECT id, name, number,recipientSender,count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog  
WHERE owner_id= $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := ModelOperationsLog{}
		err = rows.Scan(
			&model.Id,
			&model.Name,
			&model.Number,
			&model.RecipientSender,
			&model.Count,
			&model.BalanceOld,
			&model.BalanceNew,
			&model.Time,
			&model.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't get history from db: %w", err)
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get history from db: %w", err)
	}
	return models, nil
}

func (service *Service) AddNewHistory(model ModelOperationsLog) (err error) {
	log.Print("started add new history")
	log.Print("add model to db")
	_, err = service.pool.Exec(context.Background(), `
	INSERT INTO historyoperationslog(name, number,recipientSender,count, balanceold, balancenew, time, owner_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		model.Name,
		model.Number,
		model.RecipientSender,
		model.Count,
		model.BalanceOld,
		model.BalanceNew,
		model.Time,
		model.OwnerID,
	)
	if err != nil {
		log.Printf("can't exec insert add history card: %d", err)
		return err
	}
	log.Print("saved model to db")
	log.Print("finish add model to db")
	return nil
}

type ModelOperationsLog struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	RecipientSender string `json:"recipientsender"`
	Count           int64  `json:"count"`
	BalanceOld      int64  `json:"balanceold"`
	BalanceNew      int64  `json:"balancenew"`
	Time            int64  `json:"time"`
	OwnerID         int64  `json:"ownerid"`
}

type ModelTransferMoneyCardToCard struct {
	IdCardSender        int    `json:"id_card_sender"`
	NumberCardRecipient string `json:"number_card_recipient"`
	Count               int64  `json:"count"`
}
