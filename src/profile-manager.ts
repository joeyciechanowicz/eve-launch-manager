import React, { useState, useEffect } from "react";
import { Box, Text, useInput } from "ink";
import Spinner from "ink-spinner";
import path from "path";
import os from "os";
import fs from "fs";

const getConfigPath = (): string => path.join(os.homedir(), CONFIG_FILE);
const getEvePath = (): string =>
  path.join(os.homedir(), "AppData", "Roaming", "EVE Online");
const getStateFilePath = (profileName: string): string =>
  path.join(getEvePath(), `state-${profileName}.json`);
const getLauncherStateFilePath = () => path.join(getEvePath(), "state.json");

const CONFIG_FILE = "eve-launch-manager.json";
const DEFAULT_CONFIG = {
  version: 1,
  activeProfile: "main",
  profiles: ["main"],
};

export interface Config {
  activeProfile: string;
  profiles: string[];
}

export class ProfileManager {
  public readonly config: Config;
  private onSetState?: (newState: Config) => void;

  constructor() {
    let config: Config = DEFAULT_CONFIG;
    if (fs.existsSync(getConfigPath())) {
      config = JSON.parse(fs.readFileSync(getConfigPath()).toString());
    } else {
      fs.writeFileSync(getConfigPath(), JSON.stringify(config, null, "\t"));
      fs.copyFileSync(getLauncherStateFilePath(), getStateFilePath("main"));
    }
    this.config = config;
  }

  useConfig() {
    const [curr, setCurr] = useState(this.config);
    this.onSetState = setCurr;

    return curr;
  }

  private updateConfig() {
    fs.writeFileSync(getConfigPath(), JSON.stringify(this.config, null, "\t"));
    if (this.onSetState) {
      this.onSetState(this.config);
    }
  }

  createProfile(name: string, basedOff?: string) {
    let contents = "{}";
    if (basedOff) {
      contents = fs.readFileSync(getStateFilePath(basedOff)).toString();
    }

    fs.writeFileSync(getStateFilePath(name), contents);

    this.config.profiles.push(name);
    this.updateConfig();
  }

  switchProfile(profileName: string) {
    // backup current profile first
    fs.copyFileSync(
      getLauncherStateFilePath(),
      getStateFilePath(this.config.activeProfile)
    );

    // copy the new state over state.json
    fs.copyFileSync(getStateFilePath(profileName), getLauncherStateFilePath());

    this.config.activeProfile = profileName;

    this.updateConfig();
  }
}
