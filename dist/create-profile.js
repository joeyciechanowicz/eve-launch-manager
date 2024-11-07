import React, { useState } from 'react';
import { Box, Text, useInput } from 'ink';
import SelectInput from 'ink-select-input';
import TextInput from 'ink-text-input';
import Spinner from 'ink-spinner';
const CreateProfileSteps = {
    ENTER_NAME: 'ENTER_NAME',
    CHOOSE_BASE: 'CHOOSE_BASE',
    SELECT_BASE_PROFILE: 'SELECT_BASE_PROFILE',
    CREATING: 'CREATING'
};
const PROFILE_NAME_REGEX = /^[a-zA-Z0-9_-]+$/;
// Validation function with detailed error messages
const validateProfileName = (name, existingProfiles) => {
    const sanitizedName = name.trim();
    if (sanitizedName.length === 0) {
        return { isValid: false, error: 'Profile name cannot be empty' };
    }
    if (!PROFILE_NAME_REGEX.test(sanitizedName)) {
        return {
            isValid: false,
            error: 'Profile name can only contain letters, numbers, underscores, and dashes'
        };
    }
    if (existingProfiles.includes(sanitizedName)) {
        return { isValid: false, error: 'Profile name already exists' };
    }
    return { isValid: true };
};
export const CreateProfile = ({ config, profileManager, onComplete }) => {
    const [step, setStep] = useState(CreateProfileSteps.ENTER_NAME);
    const [profileName, setProfileName] = useState('');
    const [nameError, setNameError] = useState();
    const [error, setError] = useState();
    const [waitingForInput, setWaitingForInput] = useState(false);
    useInput((input, key) => {
        if (waitingForInput && (error)) {
            onComplete();
        }
        if (key.escape) {
            onComplete();
        }
    });
    // Profile name input handler with improved validation
    const handleNameSubmit = (name) => {
        const validation = validateProfileName(name, config.profiles);
        if (!validation.isValid) {
            setNameError(validation.error);
            // Only clear the input if it's a duplicate name
            if (validation.error === 'Profile name already exists') {
                setProfileName('');
            }
            return;
        }
        setProfileName(name.trim());
        setStep(CreateProfileSteps.CHOOSE_BASE);
    };
    // Real-time validation feedback as user types
    const handleNameChange = (name) => {
        setProfileName(name);
        // Clear error when user starts typing again
        if (nameError) {
            setNameError('');
        }
        // Optional: Provide real-time feedback
        if (name && !PROFILE_NAME_REGEX.test(name)) {
            setNameError('Only letters, numbers, underscores, and dashes allowed');
        }
    };
    // Base profile choice handler
    const handleBaseOrNoneChoice = (item) => {
        try {
            if (item.value === 'empty') {
                profileManager.createProfile(profileName);
                onComplete();
            }
            else if (item.value === 'existing') {
                setStep(CreateProfileSteps.SELECT_BASE_PROFILE);
            }
        }
        catch (err) {
            setError(`Failed to create profile: ${err.message}`);
            setWaitingForInput(true);
        }
    };
    // Existing profile selection handler
    const handleBaseProfileSelect = (item) => {
        setStep(CreateProfileSteps.CREATING);
        try {
            profileManager.createProfile(profileName, item.value);
            onComplete();
        }
        catch (err) {
            setError(`Failed to create profile: ${err.message}`);
            setWaitingForInput(true);
        }
    };
    // Render different steps
    if (error) {
        return (React.createElement(Box, { flexDirection: "column" },
            error && React.createElement(Text, { color: "red" }, error),
            React.createElement(Text, { color: "gray", italic: true }, "Press any key to return to main menu...")));
    }
    if (step === CreateProfileSteps.ENTER_NAME) {
        return (React.createElement(Box, { flexDirection: "column" },
            React.createElement(Text, null, "Enter profile name:"),
            React.createElement(Text, { color: "gray", dimColor: true }, "(Use only letters, numbers, underscores, and dashes)"),
            nameError && React.createElement(Text, { color: "red" }, nameError),
            React.createElement(TextInput, { value: profileName, onChange: handleNameChange, onSubmit: handleNameSubmit, placeholder: "profile-name" })));
    }
    if (step === CreateProfileSteps.CHOOSE_BASE) {
        const items = [
            { label: 'Create empty profile', value: 'empty' },
            { label: 'Base off existing profile', value: 'existing' }
        ];
        return (React.createElement(Box, { flexDirection: "column" },
            React.createElement(Text, null, "How would you like to create the profile?"),
            React.createElement(SelectInput, { items: items, onSelect: handleBaseOrNoneChoice })));
    }
    if (step === CreateProfileSteps.SELECT_BASE_PROFILE) {
        const items = config.profiles.map(profile => ({
            label: profile,
            value: profile
        }));
        return (React.createElement(Box, { flexDirection: "column" },
            React.createElement(Text, null, "Select base profile:"),
            React.createElement(SelectInput, { items: items, onSelect: handleBaseProfileSelect })));
    }
    if (step === CreateProfileSteps.CREATING) {
        return (React.createElement(Box, null,
            React.createElement(Text, { color: "green" },
                React.createElement(Spinner, { type: "dots" })),
            React.createElement(Text, null, " Creating profile...")));
    }
};
