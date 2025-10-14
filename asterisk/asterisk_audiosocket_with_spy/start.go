package main

import (
	"context"
	"log"
	"strings"

	"github.com/CyCoreSystems/ari"
	"github.com/google/uuid"
)

func handleARIStasisStart(ctx context.Context, client ari.Client, ev *ari.StasisStart) error {

	log.Printf("StasisStart. channel_name: %s", ev.Channel.Name)
	switch {
	case strings.HasPrefix(strings.ToLower(ev.Channel.Name), "pjsip/"):
		handleARIStasisStartPJSIP(ctx, client, ev)
		return nil

	case strings.HasPrefix(strings.ToLower(ev.Channel.Name), "snoop/"):
		handleARIStasisStartSnoop(ctx, client, ev)
		return nil

	case strings.HasPrefix(strings.ToLower(ev.Channel.Name), "audiosocket/"):
		handleARIStasisStartAudiosocket(ctx, client, ev)
		return nil

	case strings.HasPrefix(strings.ToLower(ev.Channel.Name), "local/"):
		handleARIStasisStartLocal(ctx, client, ev)
		return nil

	default:
	}

	return nil
}

func handleARIStasisStartPJSIP(ctx context.Context, client ari.Client, ev *ari.StasisStart) {

	log.Printf("handleARIStasisStartPJSIP. channel_name: %s", ev.Channel.Name)

	// note: Generate snoop channel from the pjsip channel and bridging the snoop channel and external media channel
	// does not work, because the snoop channel's whisper is working only when the pjsip channel has media streaming.
	// basically, the snoop channel injects the media streaming, but if the pjsip channel has no media, it does not work.
	// so, we need to play some sound on the pjsip channel to make the snoop channel work.

	// start snoop channel
	id := uuid.New().String()
	tmp, err := client.Channel().Snoop(ev.Keys().Channels().First(), id, &ari.SnoopOptions{
		App: appName,
		// Spy:     ari.DirectionBoth,
		// Whisper: ari.DirectionBoth,
		Whisper: ari.DirectionBoth,
		Spy:     ari.DirectionOut,
	})
	if err != nil {
		log.Printf("Error starting snoop: %v", err)
		return
	}
	log.Printf("Started snoop with ID: %s", tmp.ID())

	// local, err := client.Channel().Create(nil, ari.ChannelCreateRequest{
	// 	Endpoint: "Local/dummy@default",
	// 	App:      appName,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Started local channel with ID: %s", local.ID())

	play, err := client.Channel().Play(ev.Keys().Channels().First(), "test_silence", "tmp/silence_slin16_8000.wav")
	if err != nil {
		log.Printf("Error playing silence: %v", err)
		return
	}
	log.Printf("Started playing silence on channel: %s", play.ID())

}

func handleARIStasisStartSnoop(ctx context.Context, client ari.Client, ev *ari.StasisStart) {

	log.Printf("handleARIStasisStartSnoop. channel_name: %s", ev.Channel.Name)

	if errAdd := gBridgeHandleExternal.AddChannel(ev.Channel.ID); errAdd != nil {
		log.Printf("❌ Failed to add channel %s to bridge: %v", ev.Channel.ID, errAdd)
		return
	}
	log.Printf("✅ Channel added to bridge. channel_id: %s", ev.Channel.ID)
}

func handleARIStasisStartAudiosocket(ctx context.Context, client ari.Client, ev *ari.StasisStart) {

	log.Printf("handleARIStasisStartAudiosocket. channel_name: %s", ev.Channel.Name)

	if errAdd := gBridgeHandleExternal.AddChannel(ev.Channel.ID); errAdd != nil {
		log.Printf("❌ Failed to add channel %s to bridge: %v", ev.Channel.ID, errAdd)
		return
	}
	log.Printf("✅ Channel added to bridge. channel_id: %s", ev.Channel.ID)
}

func handleARIStasisStartLocal(ctx context.Context, client ari.Client, ev *ari.StasisStart) {

	log.Printf("handleARIStasisStartLocal. channel_name: %s", ev.Channel.Name)

	if errAdd := gBridgeHandleCall.AddChannel(ev.Channel.ID); errAdd != nil {
		log.Printf("❌ Failed to add channel %s to bridge: %v", ev.Channel.ID, errAdd)
		return
	}
	log.Printf("✅ Channel added to bridge. channel_id: %s", ev.Channel.ID)
}

func handleARIChannelStateChange(ctx context.Context, client ari.Client, ev *ari.ChannelStateChange) {
	log.Printf("Channel state change. channel_id: %s, new_state: %s", ev.Channel.ID, ev.Channel.State)

	if !strings.HasPrefix(strings.ToLower(ev.Channel.Name), "pjsip/") {
		return
	}

	if strings.ToLower(ev.Channel.State) != "up" {
		return
	}

	play, err := client.Channel().Play(ev.Keys().Channels().First(), "test_silence", "sound:silence_slin16_8000")
	if err != nil {
		log.Printf("Error playing silence: %v", err)
		return
	}
	log.Printf("Started playing silence on channel: %s", play.ID())

}

func handleARIPlaybackFinished(ctx context.Context, client ari.Client, ev *ari.PlaybackFinished) {
	log.Printf("handleARIPlaybackFinished. data: %v, id: %s, lang: %s, media_uri: %s, state: %s, target_uri: %s",
		ev.Playback, ev.Playback.ID, ev.Playback.Language, ev.Playback.MediaURI, ev.Playback.State, ev.Playback.TargetURI)
}
