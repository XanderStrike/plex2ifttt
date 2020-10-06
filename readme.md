# plex2ifttt

This awesome project transforms Plex webhooks into webhook events for IFTTT.

I use this to automatically dim my lights when my home theater is playing media
on Plex, and raise them once it's done (or paused).

It is currently capable of generating 3 events:

  - `plex_play` - A play event (play, resume, scrobble)
  - `plex_play_day` - A play event that happens during the day (hardcoded to
    US/pacific time right now)
  - `plex_pause` - A pause or stop event (pause, stop, finish)

It also filters the webhooks by username and player UUID, which is necessary to
make sure your lights aren't dimming when someone in another room is doing
stuff.

These events are sent to your IFTTT webhook url as the event.

```
https://maker.ifttt.com/trigger/<event>/with/key/<IFTTT-KEY>
```

## setup

Get your IFTTT key by visiting
[https://ifttt.com/maker_webhooks](https://ifttt.com/maker_webhooks/) and
clicking "Documentation."

If you do not know your player's ID, the app logs all unknown player IDs.

Deploy with Docker:

```
  docker run \
    -e USER_ID=<username> \
    -e PLAYER_UUID=<uuid> \
    -e IFTTT_KEY=<key> \
    -n plex2ifttt \
    -p 8080:8080 \
    xanderstrike/plex2ifttt
```

You must make this route available to the Plex server, either by colocating on
the same machine, the same network, or making it publicly accessible on the
internet.

Set up Plex to hit your webhook URL at `/hook`.

Create an IFTTT rule to consume a webhook with one of the above events and do a
thing you'd like.
