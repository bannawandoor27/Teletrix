package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tinyzimmer/go-gst/gst"
)

// GStreamerPipeline handles the GStreamer pipeline for audio mixing
type GStreamerPipeline struct {
	ctx       context.Context
	log       *logrus.Logger
	mutex     sync.RWMutex
	pipeline  *gst.Pipeline
	tone      *gst.Element
	mic       *gst.Element
	mixer     *gst.Element
	convert   *gst.Element
	sink      *gst.Element
	micVolume *gst.Element
	level     *gst.Element
}

// NewGStreamerPipeline creates a new GStreamer pipeline instance
func NewGStreamerPipeline(ctx context.Context, log *logrus.Logger) (*GStreamerPipeline, error) {
	// Initialize GStreamer
	gst.Init(nil)

	// Create pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline: %v", err)
	}

	// Create elements
	tone, err := gst.NewElement("audiotestsrc")
	if err != nil {
		return nil, fmt.Errorf("failed to create test tone source: %v", err)
	}
	if err := tone.Set("freq", 440.0); err != nil { // 440 Hz test tone
		return nil, fmt.Errorf("failed to set tone frequency: %v", err)
	}

	mic, err := gst.NewElement("autoaudiosrc")
	if err != nil {
		return nil, fmt.Errorf("failed to create microphone source: %v", err)
	}

	micVolume, err := gst.NewElement("volume")
	if err != nil {
		return nil, fmt.Errorf("failed to create microphone volume: %v", err)
	}
	if err := micVolume.Set("volume", 1.0); err != nil {
		return nil, fmt.Errorf("failed to set microphone volume: %v", err)
	}

	level, err := gst.NewElement("level")
	if err != nil {
		return nil, fmt.Errorf("failed to create level element: %v", err)
	}
	if err := level.Set("interval", uint64(100000000)); err != nil { // 100ms interval
		return nil, fmt.Errorf("failed to set level interval: %v", err)
	}

	mixer, err := gst.NewElement("audiomixer")
	if err != nil {
		return nil, fmt.Errorf("failed to create audio mixer: %v", err)
	}

	convert, err := gst.NewElement("audioconvert")
	if err != nil {
		return nil, fmt.Errorf("failed to create audio converter: %v", err)
	}

	sink, err := gst.NewElement("autoaudiosink")
	if err != nil {
		return nil, fmt.Errorf("failed to create audio sink: %v", err)
	}

	// Add elements to pipeline
	pipeline.AddMany(tone, mic, micVolume, level, mixer, convert, sink)

	// Link elements
	if err := tone.Link(mixer); err != nil {
		return nil, fmt.Errorf("failed to link tone to mixer: %v", err)
	}
	if err := mic.Link(micVolume); err != nil {
		return nil, fmt.Errorf("failed to link mic to volume: %v", err)
	}
	if err := micVolume.Link(level); err != nil {
		return nil, fmt.Errorf("failed to link volume to level: %v", err)
	}
	if err := level.Link(mixer); err != nil {
		return nil, fmt.Errorf("failed to link level to mixer: %v", err)
	}
	if err := mixer.Link(convert); err != nil {
		return nil, fmt.Errorf("failed to link mixer to converter: %v", err)
	}
	if err := convert.Link(sink); err != nil {
		return nil, fmt.Errorf("failed to link converter to sink: %v", err)
	}

	// Set up bus watch for error messages
	bus := pipeline.GetBus()
	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageError:
			err := msg.ParseError()
			log.Errorf("Pipeline error: %v", err)
		case gst.MessageWarning:
			warn := msg.ParseWarning()
			log.Warnf("Pipeline warning: %v", warn)
		}
		return true
	})

	return &GStreamerPipeline{
		ctx:       ctx,
		log:       log,
		pipeline:  pipeline,
		tone:      tone,
		mic:       mic,
		mixer:     mixer,
		convert:   convert,
		sink:      sink,
		micVolume: micVolume,
		level:     level,
	}, nil
}

// Start begins the audio pipeline
func (p *GStreamerPipeline) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.log.Info("Starting GStreamer pipeline...")

	// Set up bus watch for messages
	bus := p.pipeline.GetBus()
	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageStateChanged:
			p.log.Debug("Pipeline state changed")
		case gst.MessageError:
			err := msg.ParseError()
			p.log.Errorf("Pipeline error: %v", err)
		case gst.MessageWarning:
			warn := msg.ParseWarning()
			p.log.Warnf("Pipeline warning: %v", warn)
		case gst.MessageEOS:
			p.log.Info("End of stream reached")
		}
		return true
	})

	// Set pipeline state to PLAYING
	if err := p.pipeline.SetState(gst.StatePlaying); err != nil {
		p.log.Errorf("Failed to set pipeline state to PLAYING: %v", err)
		return err
	}

	// Wait briefly to ensure pipeline starts
	time.Sleep(100 * time.Millisecond)

	p.log.Info("Pipeline successfully started and playing")
	return nil
}

// Close cleans up resources
func (p *GStreamerPipeline) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.log.Info("Closing GStreamer pipeline...")
	// First set to ready state to ensure clean shutdown
	if err := p.pipeline.SetState(gst.StateReady); err != nil {
		p.log.Warnf("Failed to set pipeline state to READY: %v", err)
	}

	// Then set to null state
	if err := p.pipeline.SetState(gst.StateNull); err != nil {
		p.log.Errorf("Failed to set pipeline state to NULL: %v", err)
		return err
	}

	return nil
}

// SetToneVolume sets the volume of the test tone
func (p *GStreamerPipeline) SetToneVolume(volume float64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.tone.Set("volume", volume)
}
