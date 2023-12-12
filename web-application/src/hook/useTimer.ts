import { useState, useEffect } from "react";
import { TIMER_TIME } from "@/store/gameStore";

const useTimer = () => {
  const [timerTime, setTimerTime] = useState<number>(TIMER_TIME);
  const [isRunning, setIsRunning] = useState<boolean>(false);

  useEffect(() => {
    let interval: NodeJS.Timeout | null = null;

    if (isRunning) {
      interval = setInterval(() => {
        setTimerTime((prevTime) => (prevTime > 0 ? prevTime - 1 : 0));
      }, 1000);
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [isRunning]);

  const startTimer = () => setIsRunning(true);
  const pauseTimer = () => setIsRunning(false);
  const resetTimer = () => {
    setIsRunning(false);
    setTimerTime(TIMER_TIME);
  };

  return { timerTime, startTimer, pauseTimer, resetTimer };
};

export default useTimer;
