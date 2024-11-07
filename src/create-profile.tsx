import React, { useState } from 'react';
import { Box, Text, useInput } from 'ink';
import SelectInput from 'ink-select-input';
import TextInput from 'ink-text-input';
import Spinner from 'ink-spinner';
import { SelectedItem } from './util.js';
import { Config, ProfileManager } from './profile-manager.js';

const CreateProfileSteps = {
  ENTER_NAME: 'ENTER_NAME',
  CHOOSE_BASE: 'CHOOSE_BASE',
  SELECT_BASE_PROFILE: 'SELECT_BASE_PROFILE',
  CREATING: 'CREATING'
};

const PROFILE_NAME_REGEX = /^[a-zA-Z0-9_-]+$/;

// Validation function with detailed error messages
const validateProfileName = (name: string, existingProfiles: string[]) => {
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

export const CreateProfile = ({ config, profileManager, onComplete }: { config: Config, profileManager: ProfileManager; onComplete: () => void }) => {
  const [step, setStep] = useState(CreateProfileSteps.ENTER_NAME);
  const [profileName, setProfileName] = useState('');
  const [nameError, setNameError] = useState<string>();
  const [error, setError] = useState<string>();
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
  const handleNameSubmit = (name: string) => {
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
  const handleNameChange = (name: string) => {
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
  const handleBaseOrNoneChoice = (item: SelectedItem) => {
    try {
      if (item.value === 'empty') {
        profileManager.createProfile(profileName);
        onComplete();
      } else if (item.value === 'existing') {
        setStep(CreateProfileSteps.SELECT_BASE_PROFILE);
      }
    } catch (err: any) {
      setError(`Failed to create profile: ${err.message}`);
      setWaitingForInput(true);
    }
  };

  // Existing profile selection handler
  const handleBaseProfileSelect = (item: SelectedItem) => {
    setStep(CreateProfileSteps.CREATING);
    try {
      profileManager.createProfile(profileName, item.value);
      onComplete();
    } catch (err: any) {
      setError(`Failed to create profile: ${err.message}`);
      setWaitingForInput(true);
    }
  };

  // Render different steps
  if (error) {
    return (
      <Box flexDirection="column">
        {error && <Text color="red">{error}</Text>}
        <Text color="gray" italic>Press any key to return to main menu...</Text>
      </Box>
    );
  }

  if (step === CreateProfileSteps.ENTER_NAME) {
    return (
      <Box flexDirection="column">
        <Text>Enter profile name:</Text>
        <Text color="gray" dimColor>
          (Use only letters, numbers, underscores, and dashes)
        </Text>
        {nameError && <Text color="red">{nameError}</Text>}
        <TextInput
          value={profileName}
          onChange={handleNameChange}
          onSubmit={handleNameSubmit}
          placeholder="profile-name"
        />
      </Box>
    );
  }

  if (step === CreateProfileSteps.CHOOSE_BASE) {
    const items = [
      { label: 'Create empty profile', value: 'empty' },
      { label: 'Base off existing profile', value: 'existing' }
    ];
    return (
      <Box flexDirection="column">
        <Text>How would you like to create the profile?</Text>
        <SelectInput items={items} onSelect={handleBaseOrNoneChoice} />
      </Box>
    );
  }

  if (step === CreateProfileSteps.SELECT_BASE_PROFILE) {
    const items = config.profiles.map(profile => ({
      label: profile,
      value: profile
    }));
    return (
      <Box flexDirection="column">
        <Text>Select base profile:</Text>
        <SelectInput items={items} onSelect={handleBaseProfileSelect} />
      </Box>
    );
  }

  if (step === CreateProfileSteps.CREATING) {
    return (
      <Box>
        <Text color="green">
          <Spinner type="dots" />
        </Text>
        <Text> Creating profile...</Text>
      </Box>
    );
  }
};
