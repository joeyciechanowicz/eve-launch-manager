import React, { useEffect, useState } from 'react';
import { Box, Text, useInput } from 'ink';
import Spinner from 'ink-spinner';
import path from 'path';
import os from 'os';
import fs from 'fs';
import archiver from 'archiver';

export const BackupProcess = ({ onComplete }: { onComplete: () => void }) => {
    const [status, setStatus] = useState('Starting backup...');
    const [isComplete, setIsComplete] = useState(false);
    const [error, setError] = useState(null);
    const [waitingForInput, setWaitingForInput] = useState(false);

    useInput((input, key) => {
        if (waitingForInput) {
            onComplete();
        }
    });

    useEffect(() => {
        const createBackup = async () => {
            try {
                const evePath = path.join(os.homedir(), 'AppData', 'Roaming', 'EVE Online');
                const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
                const outputFile = path.join(os.homedir(), `eve-settings-backup-${timestamp}.zip`);

                const output = fs.createWriteStream(outputFile);
                const archive = archiver('zip', { zlib: { level: 9 } });

                output.on('close', () => {
                    setStatus(`Backup completed! File saved as: ${outputFile}`);
                    setIsComplete(true);
                    setWaitingForInput(true);
                });

                archive.on('error', (err) => {
                    throw err;
                });

                archive.on('progress', (progress) => {
                    setStatus(`Backing up files... ${Math.round(progress.entries.processed / progress.entries.total * 100)}%`);
                });

                archive.pipe(output);
                archive.directory(evePath, false);
                await archive.finalize();

            } catch (err: any) {
                setError(err?.message);
                setIsComplete(true);
                setWaitingForInput(true);
            }
        };

        createBackup();
    }, []);

    if (error) {
        return (
            <Box flexDirection="column">
                <Text color="red">Error: {error}</Text>
                <Text color="gray" italic>Press any key to return to main menu...</Text>
            </Box>
        );
    }

    return (
        <Box flexDirection="column">
            {!isComplete ? (
                <Text>
                    <Text color="green">
                        <Spinner type="dots" />
                    </Text>
                    {' '}{status}
                </Text>
            ) : (
                <>
                    <Text color="green">{status}</Text>
                    <Text color="gray" italic>Press any key to return to main menu...</Text>
                </>
            )}
        </Box>
    );
};