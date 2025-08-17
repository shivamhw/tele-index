#!/bin/bash

cat $1 | jq ".messages[]" | jq -c '{ChatId: .raw.PeerID.ChannelID, ID: .id, Size: .raw.Media.Document.Size, File: .raw.Media.Document.Attributes[0].FileName}'