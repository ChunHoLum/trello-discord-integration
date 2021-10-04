# trello-discord-integration

## Prerequisite

- [Go 1.17.1 +](https://golang.org/)
- [Task](https://taskfile.dev/#/installation) ( Alternative to `make` )

## Installation
```sh
task build # build the program
```
```sh
cp tdi.service /etc/systemd/system/tdi.service # Copy the systemd file
```
```sh
cp tdi.sample.toml /etc/tdi.toml # Copy the config file
```
```sh
vi /etc/tdi.toml # Update the tokens and settings
```
```sh
sudo systemctl daemon-reload # Reload daemon
```
```sh
sudo systemctl enable tdi # Enable service
```
```sh
sudo systemctl start tdi # Start service
```
```sh
journalctl -u tdi # View logs
```

## Current Support Trello Events
```
EVENT_ADD_ATTACHMENT_TO_CARD = "action_add_attachment_to_card"
EVENT_ADD_CHECKLIST_TO_CARD  = "action_add_checklist_to_card"
EVENT_ADD_MEMBER_TO_CARD          = "action_member_joined_card"
EVENT_COMMENT_CARD                = "action_comment_on_card"
EVENT_DELETE_ATTACHMENT_FROM_CARD = "action_delete_attachment_from_card"
EVENT_DELETE_CARD                 = "action_delete_card"
EVENT_ARCHIVE_CARD                = "action_archived_card"
EVENT_DELETE_COMMENT              = "deleteComment"
EVENT_CREATE_CARD                 = "action_create_card"
EVENT_REMOVE_CHECKLIST_FROM_CARD  = "action_remove_checklist_from_card"
EVENT_REMOVE_MEMBER_FROM_CARD = "action_member_left_card"
EVENT_UPDATE_CHECKITEM_STATE  = "action_completed_checkitem"
EVENT_MOVE_CARD_LIST_TO_LIST  = "action_move_card_from_list_to_list"
```