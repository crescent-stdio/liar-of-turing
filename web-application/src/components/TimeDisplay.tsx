import { getTimeFromTimer } from "@/utils/timeHelper";

type TimerDisplayProps = {
  timerTime: number;
  startTimer: () => void;
};

const TimerDisplay: React.FC<TimerDisplayProps> = ({
  timerTime,
  startTimer,
}) => {
  return (
    <div className="flex flex-row">
      <div className="text-center text-xl font-medium mx-2 text-gray-900">
        {getTimeFromTimer(timerTime)}
      </div>
      <button
        className="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md"
        onClick={startTimer}
      >
        Start
      </button>
    </div>
  );
};

export default TimerDisplay;
