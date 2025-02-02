package viz

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/teletrix/internal/config"
)

// Visualizer handles real-time audio visualization
type Visualizer struct {
	ctx    context.Context
	log    *logrus.Logger
	config *config.Config
	mutex  sync.RWMutex

	// Visualization state
	waveformBuffer []float64
	micVolume      float64
	toneVolume     float64
	peakVolume     float64
}

// NewVisualizer creates a new audio visualizer instance
func NewVisualizer(ctx context.Context, log *logrus.Logger, cfg *config.Config) *Visualizer {
	return &Visualizer{
		ctx:            ctx,
		log:            log,
		config:         cfg,
		waveformBuffer: make([]float64, cfg.WaveformWindowSize),
	}
}

// UpdateWaveform updates the waveform buffer with new audio samples
func (v *Visualizer) UpdateWaveform(samples []float64) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	// Copy new samples into the buffer
	copy(v.waveformBuffer, samples)
}

// UpdateVolumes updates the volume levels for visualization
func (v *Visualizer) UpdateVolumes(micVol, toneVol float64) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.micVolume = micVol
	v.toneVolume = toneVol

	// Update peak volume with decay
	if micVol > v.peakVolume {
		v.peakVolume = micVol
	} else {
		v.peakVolume *= v.config.PeakDecayRate
	}
}

// GetWaveform returns the current waveform buffer
func (v *Visualizer) GetWaveform() []float64 {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	// Return a copy to prevent external modifications
	buffer := make([]float64, len(v.waveformBuffer))
	copy(buffer, v.waveformBuffer)
	return buffer
}

// GetVolumes returns the current volume levels
func (v *Visualizer) GetVolumes() (mic, tone, peak float64) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.micVolume, v.toneVolume, v.peakVolume
}
