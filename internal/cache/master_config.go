package config_cache

import (
	"context"
	"fmt"

	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/cache"
	"github.com/loop-kar/pixie/errs"
)

var masterConfig *cache.Cache

const scopedKeyFormat = "%s.%s-chan-%d"
const contextScopedKeyFormat = "%s-chan-%d"

func InitMasterConfig(repo repository.MasterConfigRepository) *errs.XError {

	if masterConfig != nil {
		return nil
	} else {
		masterConfig = &cache.Cache{}
	}

	ctx := context.Background()
	configs, err := repo.LoadAll(&ctx)
	if err != nil {
		return err
	}

	var value string
	for _, config := range configs {
		if config.UseDefault {
			value = config.DefaultValue
		} else {
			value = config.CurrentValue
		}

		masterConfig.Set(scopedKey(config.Type, config.Name, config.ChannelId), value)
	}

	return nil
}

// GetValue checks if the value is present in the cache, if it does not
// it returns nil,false
func GetValue(ctx *context.Context, key string) (interface{}, bool) {
	return masterConfig.Get(getScopedKeyFromContext(ctx, key))
}

// SetValue just sets the value to the key, if key already exists
// it updates the value
func SetValue(ctx *context.Context, key string, value interface{}) {
	masterConfig.Set(getScopedKeyFromContext(ctx, key), value)
}

func getScopedKeyFromContext(ctx *context.Context, key string) string {
	return fmt.Sprintf(contextScopedKeyFormat, key, utils.GetChannelId(ctx))
}

func scopedKey(keyType, name string, channelId uint) string {
	return fmt.Sprintf(scopedKeyFormat, keyType, name, channelId)
}
