import React, { } from 'react';
import { Box, useInput } from 'ink';
import SelectInput from 'ink-select-input';
import { Config, ProfileManager } from './profile-manager.js';
import { SelectedItem } from './util.js';

export const LoadProfile = ({ config, profileManager, onComplete }: { config: Config, profileManager: ProfileManager; onComplete: () => void }) => {
    const profiles = config.profiles.map(profile => ({
        value: profile,
        label: profile === config.activeProfile ? `${profile} (active)` : profile
    }));

    const handleProfileSelected = (item: SelectedItem) => {
        if (item.value === config.activeProfile) {
            return;
        } else {
            profileManager.switchProfile(item.value);
            onComplete();
        }
    };

    useInput((input, key) => {
        if (key.escape) {
            onComplete();
        }
    });

    return (
        <Box>
            <SelectInput items={profiles} onSelect={handleProfileSelected} />
        </Box>
    );
};
