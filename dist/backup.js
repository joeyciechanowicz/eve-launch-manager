import React, { useEffect, useState } from 'react';
import { Box, Text, useInput } from 'ink';
import Spinner from 'ink-spinner';
import path from 'path';
import os from 'os';
import fs from 'fs';
import archiver from 'archiver';
export const BackupProcess = ({ onComplete }) => {
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
            }
            catch (err) {
                setError(err?.message);
                setIsComplete(true);
                setWaitingForInput(true);
            }
        };
        createBackup();
    }, []);
    if (error) {
        return (React.createElement(Box, { flexDirection: "column" },
            React.createElement(Text, { color: "red" },
                "Error: ",
                error),
            React.createElement(Text, { color: "gray", italic: true }, "Press any key to return to main menu...")));
    }
    return (React.createElement(Box, { flexDirection: "column" }, !isComplete ? (React.createElement(Text, null,
        React.createElement(Text, { color: "green" },
            React.createElement(Spinner, { type: "dots" })),
        ' ',
        status)) : (React.createElement(React.Fragment, null,
        React.createElement(Text, { color: "green" }, status),
        React.createElement(Text, { color: "gray", italic: true }, "Press any key to return to main menu...")))));
};
