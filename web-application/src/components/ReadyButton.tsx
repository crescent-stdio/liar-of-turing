import { User } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";

type ReadyButtonProps = {
  userData: User;
  sendMessage: (message: WsJsonRequest) => void;
};
export default function ReadyButton({
  userData,
  sendMessage,
}: ReadyButtonProps) {
  const handleReady = () => {
    const jsonData: WsJsonRequest = {
      action: "user_is_ready",
      user: userData,
      timestamp: Date.now(),
      message: "",
    };
    sendMessage(jsonData);
  };
  return (
    <div>
      <button
        className="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-[#3b82f6] focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
        onClick={handleReady}
      >
        Join the game
      </button>
    </div>
  );
}
