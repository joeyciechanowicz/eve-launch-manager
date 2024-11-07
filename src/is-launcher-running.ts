import { useEffect, useState } from "react";
import { exec } from 'node:child_process';

export function useIsLauncherRunning(): { isRunning: boolean, checkedOnce: boolean } {
    const [isRunning, setIsRunning] = useState(false);
    const [checkedOnce, setCheckedOnce] = useState(false);

    useEffect(() => {
        exec('tasklist', (err, stdout, stderr) => {
            if (stdout.indexOf('eve-online.exe') !== -1) {
                setIsRunning(true);
            } else {
                setIsRunning(false);
            }
            setCheckedOnce(true);
        });

        const id = setInterval(() => {
            exec('tasklist', (err, stdout, stderr) => {
                if (stdout.indexOf('eve-online.exe') !== -1) {
                    setIsRunning(true);
                } else {
                    setIsRunning(false);
                }
            });
        }, 250);

        return () => clearInterval(id);
    }, []);

    return { isRunning, checkedOnce };
}