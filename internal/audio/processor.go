package audio

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Processor handles audio stream processing and voice activity detection
type Processor struct {
	ctx    context.Context
	log    *logrus.Logger
	mutex  sync.RWMutex
	active bool
}

// NewProcessor creates a new audio processor instance
func NewProcessor(ctx context.Context, log *logrus.Logger) (*Processor, error) {
	return &Processor{
		ctx:    ctx,
		log:    log,
		active: false,
	}, nil
}

// Start begins the audio processing
func (p *Processor) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.log.Info("Starting audio processor...")
	p.active = true
	return nil
}

// Close cleans up resources
func (p *Processor) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.log.Info("Closing audio processor...")
	p.active = false
	return nil
}

// IsActive returns whether voice activity is detected
func (p *Processor) IsActive() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.active
}
