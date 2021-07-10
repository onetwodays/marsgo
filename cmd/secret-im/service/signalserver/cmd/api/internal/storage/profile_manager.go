package storage

import (
	"database/sql"
	"encoding/json"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/model"
)

var CACHE_PREFIX = "profiles::"

type ProfileManager struct {
}

func (ProfileManager) Get(uuid string, version string) (*entities.VersionedProfile, error) {
	cmd := internal.client.HGet(CACHE_PREFIX+uuid, version)
	if cmd.Err() != nil || len(cmd.Val()) == 0 {
		dbRecode, err := internal.profileDB.FindOneByUUIDVersion(uuid, version)
		if err != nil {
			return nil, err
		}

		if dbRecode!=nil{
			return &entities.VersionedProfile{
				Version:    dbRecode.Version,
				Name:       dbRecode.Name,
				Avatar:     dbRecode.Avatar.String,
				Commitment: dbRecode.Commitment.String,
			}, nil
		}

	}

	var versionedProfile entities.VersionedProfile
	err := json.Unmarshal([]byte(cmd.Val()), &versionedProfile)
	if err != nil {
		return nil, err
	}
	return &versionedProfile, nil
}

func (ProfileManager) DeleteAll(uuid string) error {
	err := internal.client.Del(CACHE_PREFIX + uuid).Err()
	if err != nil {
		return err
	}
	return internal.profileDB.DeleteByUUID(uuid)
}

func (ProfileManager) Set(uuid string, versionedProfile *entities.VersionedProfile) error {
	jsb, _ := json.Marshal(versionedProfile)
	err := internal.client.HSet(CACHE_PREFIX+uuid, versionedProfile.Version, string(jsb)).Err()
	if err != nil {
		return err
	}

	dbRecord, err := internal.profileDB.FindOneByUUIDVersion(uuid, versionedProfile.Version)
	if err != nil && err != sqlx.ErrNotFound {
		return err
	}
	if dbRecord != nil {
		dbRecord.Name = versionedProfile.Name
		dbRecord.Version = versionedProfile.Version
		dbRecord.Avatar = sql.NullString{Valid: true, String: versionedProfile.Avatar}
		dbRecord.Commitment = sql.NullString{Valid: true, String: versionedProfile.Commitment}
		return internal.profileDB.Update(*dbRecord)

	}
	dbRecord = &model.TProfiles{
		Uuid:       uuid,
		Version:    versionedProfile.Version,
		Name:       versionedProfile.Name,
		Avatar:     sql.NullString{Valid: true, String: versionedProfile.Avatar},
		Commitment: sql.NullString{Valid: true, String: versionedProfile.Commitment},
	}
	_, err = internal.profileDB.Insert(*dbRecord)
	return err

}
