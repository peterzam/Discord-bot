package controller

import (
	"Discord-bot/util"
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
)

func Server(command string, cid string, s *discordgo.Session, m *discordgo.MessageCreate) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	switch command {
	case "update":
		out, err := cli.ImagePull(ctx, "karlrees/docker_bedrockserver", types.ImagePullOptions{})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Panic: Command \"update\" error.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Image is pulling(updating).")
		}
		defer out.Close()

	case "start":
		if err := cli.ContainerStart(ctx, cid, types.ContainerStartOptions{}); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Panic: Command \"start\" error.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Container is starting.")
		}

	case "stop":
		if err := cli.ContainerStop(ctx, cid, nil); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Panic: Command \"stop\" error.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Container is stopping.")
		}

	case "restart":
		if err := cli.ContainerRestart(ctx, cid, nil); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Panic: Command \"restart\" error.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Container is restarting.")
		}

	case "status":

		out, err := cli.ContainerStats(ctx, cid, false)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Panic: Command \"status\" error.")
		}
		defer out.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(out.Body)
		statJSON := buf.String()
		value := gjson.Get(statJSON, "cpu_stats.cpu_usage.total_usage")

		if value.Int() == 0 {
			s.ChannelMessageSend(m.ChannelID, "Container is stopped.")
		} else if value.Int() > 0 {
			s.ChannelMessageSend(m.ChannelID, "Container is running.")
		}
	case "stats":
		s.ChannelMessageSend(m.ChannelID, util.GetCPUinfo())
		s.ChannelMessageSend(m.ChannelID, util.GetMeminfo())
		s.ChannelMessageSend(m.ChannelID, util.GetDiskinfo())
	default:
		s.ChannelMessageSend(m.ChannelID, "!h - Help Menu\n"+
			"!s {commands} - Server Commands\n"+
			"{commands}\n"+
			"    update   - Update Image\n"+
			"    start       - Start Container\n"+
			"    stop       - Stop Container\n"+
			"    restart   - Restart Container\n"+
			"    status     - Container Status\n"+
			"    stats      - Server Stats\n\n")
	}
}
