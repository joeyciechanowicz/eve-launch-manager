import { useState } from 'react';
import path from 'path';
import os from 'os';
import fs from 'fs';
const getConfigPath = () => path.join(os.homedir(), CONFIG_FILE);
const getEvePath = () => path.join(os.homedir(), 'AppData', 'Roaming', 'EVE Online');
const getStateFilePath = (profileName) => path.join(getEvePath(), `state-${profileName}.json`);
const getLauncherStateFilePath = () => path.join(getEvePath(), 'state.json');
const CONFIG_FILE = 'eve-launch-manager.json';
const DEFAULT_CONFIG = {
    version: 1,
    activeProfile: 'main',
    profiles: ['main']
};
export class ProfileManager {
    constructor() {
        let config = DEFAULT_CONFIG;
        if (fs.existsSync(getConfigPath())) {
            config = JSON.parse(fs.readFileSync(getConfigPath()).toString());
        }
        else {
            fs.writeFileSync(getConfigPath(), JSON.stringify(config, null, '\t'));
            fs.copyFileSync(getLauncherStateFilePath(), getStateFilePath('main'));
        }
        this.config = config;
    }
    useConfig() {
        const [curr, setCurr] = useState(this.config);
        this.onSetState = setCurr;
        return curr;
    }
    updateConfig() {
        fs.writeFileSync(getConfigPath(), JSON.stringify(this.config, null, '\t'));
        if (this.onSetState) {
            this.onSetState(this.config);
        }
    }
    createProfile(name, basedOff) {
        let contents = '{}';
        if (basedOff) {
            contents = fs.readFileSync(getStateFilePath(basedOff)).toString();
        }
        fs.writeFileSync(getStateFilePath(name), contents);
        this.config.profiles.push(name);
        this.updateConfig();
    }
    switchProfile(profileName) {
        // backup current profile first
        fs.copyFileSync(getLauncherStateFilePath(), getStateFilePath(this.config.activeProfile));
        // copy the new state over state.json
        fs.copyFileSync(getStateFilePath(profileName), getLauncherStateFilePath());
        this.config.activeProfile = profileName;
        this.updateConfig();
    }
}
