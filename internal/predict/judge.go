package predict

import (
	"context"
	"time"

	"github.com/galaxy-future/cudgx/common/types"
	"github.com/galaxy-future/cudgx/internal/clients"
	"github.com/galaxy-future/cudgx/internal/predict/config"
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	redundancy_keeper "github.com/galaxy-future/cudgx/internal/predict/redundancy-keeper"
)

var predictor *Predictor

//Predictor 指标预测器，根据指标预测冗余度
type Predictor struct {
	config *config.Param
}

//InitializeByConfig 初始化Predictor
func InitializeByConfig(theConfig *config.Config) error {
	if theConfig.Predict == nil {
		theConfig.Predict = &config.Param{}
	}
	if theConfig.Predict.MinimalSampleCount == 0 {
		theConfig.Predict.MinimalSampleCount = consts.DefaultPredictMinCount
	}
	if theConfig.Predict.RunDuration.Duration == 0 {
		theConfig.Predict.RunDuration = types.Duration{Duration: 60 * time.Second}
	}
	if theConfig.Predict.RuleConcurrency == 0 {
		theConfig.Predict.RuleConcurrency = consts.DefaultRuleConcurrency
	}
	if theConfig.Predict.LookbackDuration.Duration == 0 {
		theConfig.Predict.LookbackDuration = types.Duration{Duration: time.Minute}
	}
	if theConfig.Predict.MetricSendDuration.Duration == 0 {
		theConfig.Predict.MetricSendDuration = types.Duration{Duration: 5 * time.Second}
	}

	err := clients.InitClickhouseRdCli(theConfig.Clickhouse)
	if err != nil {
		return err
	}
	predictor = &Predictor{
		config: theConfig.Predict,
	}
	err = clients.InitDBClient(theConfig.Database)
	if err != nil {
		return err
	}
	clients.InitializeBridgxClient(theConfig.Xclient.BridgxServerAddress)
	clients.InitializeSchedulxClient(theConfig.Xclient.SchedulxServerAddress)
	redundancy_keeper.InitRedundancyKeeper(theConfig.Predict)
	return nil
}

func StartRedundancyKeeper(ctx context.Context) {
	redundancy_keeper.Start(ctx)
}
