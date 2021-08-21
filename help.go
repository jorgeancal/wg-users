package main

const HeadHelp = `usage: wg-users [actions] [<users>]

[actions]:
`

const CreateHelp = " create:\tcreates the user/users. It will create the config for the WireGuard client in the home root folder.\n\n"
const DeleteHelp = " delete:\tdeletes the user/users. It will remove the users from WireGuard\n\n"
const UpdateHelp = " update:\tupdates the user/users. It will delete and create again the credentials of the user/users.\n\n"
const ListHelp = " list:\t\tlist the users of we have\n\n"
const ConfigHelp = " config:\tYou will be able to edit some values for the config. Flags:\n"
