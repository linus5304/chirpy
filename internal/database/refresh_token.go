package database

import "time"

type RefreshToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"user_id"`
}

func (db *DB) SaveRefreshToken(userID int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStructure.RefreshTokens[token] = RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// func (db *DB) GetRefreshToken(token string) (RefreshToken, error) {
// 	dbStructure, err := db.loadDB()
// 	if err != nil {
// 		return RefreshToken{}, err
// 	}
// 	refreshToken, ok := dbStructure.RefreshTokens[token]
// 	if !ok {
// 		return RefreshToken{}, err
// 	}
// 	return refreshToken, nil
// }

func (db *DB) RevokeRefreshToken(token string) error {
	dBStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dBStructure.RefreshTokens, token)
	err = db.writeDB(dBStructure)
	if err != nil {
		return err
	}
	return nil

}

func (db *DB) UserForRefreshToken(token string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return User{}, ErrNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUser(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
