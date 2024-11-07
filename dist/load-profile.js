import React from 'react';
import { Box, useInput } from 'ink';
import SelectInput from 'ink-select-input';
export const LoadProfile = ({ config, profileManager, onComplete }) => {
    const profiles = config.profiles.map(profile => ({
        value: profile,
        label: profile === config.activeProfile ? `${profile} (active)` : profile
    }));
    const handleProfileSelected = (item) => {
        if (item.value === config.activeProfile) {
            return;
        }
        else {
            profileManager.switchProfile(item.value);
            onComplete();
        }
    };
    useInput((input, key) => {
        if (key.escape) {
            onComplete();
        }
    });
    return (React.createElement(Box, null,
        React.createElement(SelectInput, { items: profiles, onSelect: handleProfileSelected })));
};
