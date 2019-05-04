package database

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/jackc/pgx"
	"github.com/xlab/closer"
	"os"
	"time"
)

const (
	maxConnections = 3
	acquireTimeout = 3 * time.Second
)

const (
	SQLNoRows   = "sql: no rows in result set"
	NoUserFound = "БД: Не был найден юзер"
)

var db *dbManager

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
}

type dbManager struct {
	dataBase *pgx.ConnPool
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func loadConfiguration(file string) (pgxConfig pgx.ConnConfig) {
	configFile, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var config dbConfig
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = configFile.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	pgxConfig.Host = config.Host
	pgxConfig.User = config.User
	pgxConfig.Password = config.Password
	pgxConfig.Database = config.DBName
	pgxConfig.Port = config.Port
	return
}

func init() {

	pgxConfig := loadConfiguration(os.Getenv("BASEPATH") + "/user_microservice/database/config.json")
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig, MaxConnections: maxConnections, AcquireTimeout: acquireTimeout}

	dataBase, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}

	fmt.Println("DB connection opened")

	db = &dbManager{
		dataBase: dataBase,
	}

	closer.Bind(closeConnection)
}

func closeConnection() {
	db.dataBase.Close()
	fmt.Println("DB connection closed")
}

func GetInstance() *dbManager {
	return db
}

func (db *dbManager) GetUserByID(userID uint64) (user services.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT id, email, nickname, avatar_path FROM public.users WHERE id = $1`, userID)
	err = row.Scan(&user.Id, &user.Email, &user.Nickname, &user.AvatarPath)
	if err != nil && err.Error() == SQLNoRows {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetUserByEmail(email string) (user services.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT id, email, nickname, avatar_path FROM public.users WHERE email = $1`, email)
	err = row.Scan(&user.Id, &user.Email, &user.Nickname, &user.AvatarPath)
	if err != nil {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true

	if !txOK {
		err = tx.Rollback()
		return
	}
	return
}

func (db *dbManager) AddUser(email, nickname, avatarPath string, password, salt []byte) (user services.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`INSERT INTO public.users (email, password, salt, nickname, avatar_path)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, email, nickname, avatar_path`,
		email, password, salt, nickname, avatarPath)
	err = row.Scan(&user.Id, &user.Email, &user.Nickname, &user.AvatarPath)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) UpdateUser(user *services.User) (err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`UPDATE public.users 
			SET nickname = CASE
				WHEN $1 = '' THEN nickname ELSE $1 END,
			    avatar_path = CASE
				WHEN $2 = '' THEN avatar_path ELSE $2 END
			WHERE id = $3`, user.Nickname, user.AvatarPath, user.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetUsers(page *services.PageData) (leaderBoardPage *services.LeaderBoardPage, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	rows, err := db.dataBase.Query(`SELECT id, email, nickname, avatar_path, rating, win_rate, matches FROM public.users ORDER BY rating DESC`)
	if err != nil {
		return
	}

	var profiles []*services.Profile
	var profile services.Profile
	for rows.Next() {
		err = rows.Scan(
			&profile.User.Id,
			&profile.User.Email,
			&profile.User.Nickname,
			&profile.User.AvatarPath,
			&profile.Rating,
			&profile.WinRate,
			&profile.Matches)
		if err != nil {
			return
		}
		profiles = append(profiles, &profile)
	}
	leaderBoardPage = &services.LeaderBoardPage{
		Leaders: profiles,
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetProfile(userID uint64) (profile services.Profile, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT id, email, nickname, avatar_path, rating, win_rate, matches FROM public.users WHERE id = $1`, userID)
	err = row.Scan(
		&profile.User.Id,
		&profile.User.Email,
		&profile.User.Nickname,
		&profile.User.AvatarPath,
		&profile.Rating,
		&profile.WinRate,
		&profile.Matches)
	if err != nil {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true

	if !txOK {
		err = tx.Rollback()
		return
	}
	return
}

func (db *dbManager) GetUserSignature(email string) (id uint64, password, salt []byte, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT id, password, salt FROM public.users WHERE email = $1`, email)
	err = row.Scan(&id, &password, &salt)
	if err != nil {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true

	if !txOK {
		err = tx.Rollback()
		return
	}
	return
}

func (db *dbManager) CleanerDBForTests() (err error) {
	//TODO remove?
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(`TRUNCATE TABLE public.users RESTART IDENTITY`)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}
