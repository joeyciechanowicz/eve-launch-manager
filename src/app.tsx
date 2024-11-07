import React, { useState } from 'react';
import { Box, Text, useApp, useInput, } from 'ink';
import SelectInput from 'ink-select-input';
import Spinner from 'ink-spinner';
import { BackupProcess } from './backup.js';
import { CreateProfile } from './create-profile.js';
import { Config, ProfileManager } from './profile-manager.js';
import { LoadProfile } from './load-profile.js';
import { useIsLauncherRunning } from './is-launcher-running.js';
import { SelectedItem } from './util.js';

const WizardSteps = {
    MAIN_MENU: 'MAIN_MENU',
    BACKUP: 'BACKUP',
    CREATE_PROFILE: 'CREATE_PROFILE',
    LOAD_PROFILE: 'LOAD_PROFILE'
};

const MainMenu = ({ onSelect, config }: { onSelect: (selected: SelectedItem) => void; config: Config }) => {
    const items = [
        { label: 'Backup existing settings', value: WizardSteps.BACKUP },
        { label: 'Create new profile', value: WizardSteps.CREATE_PROFILE },
        { label: 'Load profile', value: WizardSteps.LOAD_PROFILE }
    ];

    return (
        <Box flexDirection="column">
            <Text bold>EVE Online Settings Manager</Text>
            <Text>Active Profile: {config.activeProfile || 'None'}</Text>
            <Text>Available Profiles: {config.profiles.length}</Text>
            <Text>Select an option:</Text>
            <SelectInput items={items} onSelect={onSelect} />
        </Box>
    );
};


const profileManager = new ProfileManager();

const App = () => {
    const { isRunning, checkedOnce } = useIsLauncherRunning();
    const [currentStep, setCurrentStep] = useState(WizardSteps.MAIN_MENU);
    const config = profileManager.useConfig();

    const { exit } = useApp();
    useInput((input, key) => {
        if (input === "c" && key.ctrl) {
            exit();
        }
    });

    const handleSelect = (item: SelectedItem) => {
        setCurrentStep(item.value);
    };

    const returnToMainMenu = () => setCurrentStep(WizardSteps.MAIN_MENU);

    if (!checkedOnce) {
        return <Spinner />
    }

    if (isRunning) {
        return <Box>
            <Text color={'red'}>EVE Launcher is open, close it before making changes to profiles <Spinner /></Text>
        </Box>
    }

    return (<>
        <Text color={'red'}>Warning: do not manually edit the config file as config drift can cause your profile to get overwritten</Text>
        <Box flexDirection="column" padding={1}>
            {currentStep === WizardSteps.MAIN_MENU && (
                <MainMenu onSelect={handleSelect} config={config} />
            )}
            {currentStep === WizardSteps.BACKUP && (
                <BackupProcess onComplete={returnToMainMenu} />
            )}
            {currentStep === WizardSteps.CREATE_PROFILE && (
                <CreateProfile config={config} profileManager={profileManager} onComplete={returnToMainMenu} />
            )}
            {currentStep === WizardSteps.LOAD_PROFILE && <LoadProfile config={config} profileManager={profileManager} onComplete={returnToMainMenu} />}
        </Box>
    </>
    );
};

export default App;