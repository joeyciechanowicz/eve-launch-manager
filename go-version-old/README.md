# eve-launch-manager

![logo](winres/icon.png)

Tool to manage changing the set of characters that the EVE Launcher lists.

1. Download the latest version from [releases](https://github.com/joeyciechanowicz/eve-launch-manager/releases)
2. Place the `eve-launch-manager.exe` in your Documents (or wherever you want to put it really) and run it

![example of program running with list of options](image.png)

3. Before changing profiles, select the backup option. This will create a `.zip` of your entire `AppData/Roaming/EVE Online` folder and place it in `C:/%username%/eve-settings-2024-11-12_13-05.zip`. Should you need to restore from backup you will have to manually extract this folder and put the contents back in your EVE settings folder.


> [!NOTE] 
> You can press `esc` to cancel the current option and go back to the main menu at any stage. `q` or `cntrl+c` will exit the tool.

## How the tool works

The EVE Launcher stores its list of current accounts in `AppData/Roaming/EVE Online/state.json`, this tool creates multiple versions of that file. When it is run for the first ever time it will copy `state.json` to `state-main.json` and create a single entry in the tools configuration file (`C:\%username%\eve-launch-manager.json`). When you switch profiles the tool copies `state-$profilename.json` to `state.json`. The `state-*` files in your EVE settings folder **must** match what is in the `eve-launch-manager.json` configuration file. Discrepencies will cause errors.

> [!IMPORTANT] 
> Only edit the `eve-launch-manager.json` file if you understand it's structure, and ensure it is kept in sync with the `AppData/Roaming/EVE Online` folder

## Creating a profile

1. To create a profile select the `Create a profile` option
2. Give the profile a name. It must be numbers, letters, dashes and underscores
3. Select a base profile, or `None` for a new list of accounts. This is the profile that you want to extend with other accounts, and accounts that are in the `base` you select will be present in this new profile.
4. Done.

> [!NOTE] 
> The typical workflow is to create a new profile with no base and to add some common accounts to it (trading alt, corp CEO etc) and name it `common-base` or similar. You then use this as the base for your other accounts so you always start with the same set of accounts available.

## Switching profiles

1. Close the EVE Launcher
1. Select `Load a profile` and use the arrow keys to select the one you want
1. Open the EVE Launcher
