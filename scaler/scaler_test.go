package scaler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/sokar"
	"github.com/thomasobenaus/sokar/test/scaler"
)

func TestNew(t *testing.T) {

	cfg := Config{}
	scaler, err := cfg.New(nil)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg = Config{}
	scaler, err = cfg.New(scaTgt)
	assert.NoError(t, err)
	assert.NotNil(t, scaler)
}

func TestScaleBy_JobDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// dead job - error
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, fmt.Errorf("internal error"))
	result := scaler.ScaleBy(2)
	assert.Equal(t, sokar.ScaleFailed, result.State)

	// dead job
	scaTgt.EXPECT().IsJobDead(jobname).Return(true, nil)
	result = scaler.ScaleBy(2)
	assert.Equal(t, sokar.ScaleIgnored, result.State)
}

func TestScaleBy_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale up
	currentJobCount := uint(0)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	result := scaler.ScaleBy(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale up - relative
	currentJobCount = uint(1)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(3)).Return(nil)
	result = scaler.ScaleBy(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale up - max hit
	currentJobCount = uint(4)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(5)).Return(nil)
	result = scaler.ScaleBy(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}

func TestScaleBy_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale down
	currentJobCount := uint(3)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result := scaler.ScaleBy(-2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale up - min hit
	currentJobCount = uint(2)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result = scaler.ScaleBy(-5)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}

func TestScaleBy_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale down
	currentJobCount := uint(5)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	result := scaler.ScaleBy(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}
